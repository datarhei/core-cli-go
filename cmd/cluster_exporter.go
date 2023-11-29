package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

type clusterHLSSessionCollector struct {
	client coreclient.RestClient
	node   string

	hlsSessionsDesc      *prometheus.Desc
	hlsSessionsBytesDesc *prometheus.Desc
}

func newClusterHLSSessionCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterHLSSessionCollector{
		client: client,
		node:   node,
		hlsSessionsDesc: prometheus.NewDesc(
			"cluster_node_hls_sessions",
			"Cluster node HLS sessions",
			[]string{"node"}, nil),
		hlsSessionsBytesDesc: prometheus.NewDesc(
			"cluster_node_hls_tx_bytes",
			"Cluster node HLS sent bytes",
			[]string{"node"}, nil),
	}
}

func (c *clusterHLSSessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hlsSessionsDesc
	ch <- c.hlsSessionsBytesDesc
}

func (c *clusterHLSSessionCollector) Collect(ch chan<- prometheus.Metric) {
	list, err := c.client.Sessions([]string{"hls"})
	if err != nil {
		return
	}

	sessions := list["hls"]

	bytes := sessions.Summary.TotalTxBytes * 1024 * 1024

	for _, sess := range sessions.Active.SessionList {
		bytes += sess.TxBytes
	}

	ch <- prometheus.MustNewConstMetric(c.hlsSessionsDesc, prometheus.GaugeValue, float64(len(sessions.Active.SessionList)), c.node)
	ch <- prometheus.MustNewConstMetric(c.hlsSessionsBytesDesc, prometheus.CounterValue, float64(bytes), c.node)
}

type clusterNodeCollector struct {
	client coreclient.RestClient

	cpuLimitDesc   *prometheus.Desc
	cpuCurrentDesc *prometheus.Desc
	cpuCoresDesc   *prometheus.Desc
	memLimitDesc   *prometheus.Desc
	memCurrentDesc *prometheus.Desc
	throttlingDesc *prometheus.Desc
	degradedDesc   *prometheus.Desc
}

func newClusterNodeCollector(client coreclient.RestClient) prometheus.Collector {
	return &clusterNodeCollector{
		client: client,
		cpuLimitDesc: prometheus.NewDesc(
			"cluster_node_cpu_limit_percent",
			"Cluster node CPU limit in percent",
			[]string{"node"}, nil),
		cpuCurrentDesc: prometheus.NewDesc(
			"cluster_node_cpu_current_percent",
			"Cluster node CPU current in percent",
			[]string{"node"}, nil),
		cpuCoresDesc: prometheus.NewDesc(
			"cluster_node_cpu_cores",
			"Cluster node CPU cores",
			[]string{"node"}, nil),
		memLimitDesc: prometheus.NewDesc(
			"cluster_node_mem_limit_bytes",
			"Cluster node memory limit in bytes",
			[]string{"node"}, nil),
		memCurrentDesc: prometheus.NewDesc(
			"cluster_node_mem_current_bytes",
			"Cluster node memory current in bytes",
			[]string{"node"}, nil),
		throttlingDesc: prometheus.NewDesc(
			"cluster_node_throttling",
			"Cluster node throttling",
			[]string{"node"}, nil),
		degradedDesc: prometheus.NewDesc(
			"cluster_node_degraded",
			"Cluster node degraded",
			[]string{"node"}, nil),
	}
}

func (c *clusterNodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuLimitDesc
	ch <- c.cpuCurrentDesc
	ch <- c.cpuCoresDesc
	ch <- c.memLimitDesc
	ch <- c.memCurrentDesc
	ch <- c.throttlingDesc
	ch <- c.degradedDesc
}

func (c *clusterNodeCollector) Collect(ch chan<- prometheus.Metric) {
	about, err := c.client.Cluster()
	if err != nil {
		return
	}

	for _, node := range about.Nodes {
		if node.ID != about.ID {
			continue
		}

		throttling := .0
		if node.Resources.IsThrottling {
			throttling = 1.0
		}

		ch <- prometheus.MustNewConstMetric(c.cpuLimitDesc, prometheus.GaugeValue, node.Resources.CPULimit, about.ID)
		ch <- prometheus.MustNewConstMetric(c.cpuCurrentDesc, prometheus.GaugeValue, node.Resources.CPU, about.ID)
		ch <- prometheus.MustNewConstMetric(c.cpuCoresDesc, prometheus.GaugeValue, node.Resources.NCPU, about.ID)
		ch <- prometheus.MustNewConstMetric(c.memLimitDesc, prometheus.GaugeValue, float64(node.Resources.MemLimit), about.ID)
		ch <- prometheus.MustNewConstMetric(c.memCurrentDesc, prometheus.GaugeValue, float64(node.Resources.Mem), about.ID)
		ch <- prometheus.MustNewConstMetric(c.throttlingDesc, prometheus.GaugeValue, throttling, about.ID)

		break
	}

	degraded := .0
	if about.Degraded {
		degraded = 1.0
	}

	ch <- prometheus.MustNewConstMetric(c.degradedDesc, prometheus.GaugeValue, degraded, about.ID)
}

type clusterProcessCollector struct {
	client coreclient.RestClient
	node   string

	processDesc           *prometheus.Desc
	processInputFPSDesc   *prometheus.Desc
	processInputSpeedDesc *prometheus.Desc
}

func newClusterProcessCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterProcessCollector{
		client: client,
		node:   node,
		processDesc: prometheus.NewDesc(
			"cluster_process",
			"Cluster processes by state",
			[]string{"node", "state"}, nil),
		processInputFPSDesc: prometheus.NewDesc(
			"cluster_process_input_fps",
			"Cluster process input FPS by id and input",
			[]string{"node", "id", "input"}, nil),
		processInputSpeedDesc: prometheus.NewDesc(
			"cluster_process_input_speed",
			"Cluster process input speed by id and input",
			[]string{"node", "id", "input"}, nil),
	}
}

func (c *clusterProcessCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.processDesc
	ch <- c.processInputFPSDesc
	ch <- c.processInputSpeedDesc
}

func (c *clusterProcessCollector) Collect(ch chan<- prometheus.Metric) {
	processes, err := c.client.ProcessList(coreclient.ProcessListOptions{
		Filter: []string{"state"},
	})
	if err != nil {
		return
	}

	states := map[string]uint64{}

	for _, p := range processes {
		if p.State == nil {
			continue
		}

		states[p.State.State]++

		if p.State.State != "running" {
			continue
		}

		if !strings.HasSuffix(p.ID, ":main") {
			continue
		}

		for _, input := range p.State.Progress.Input {
			if input.Type != "video" {
				continue
			}

			ch <- prometheus.MustNewConstMetric(c.processInputFPSDesc, prometheus.GaugeValue, float64(input.FPS), c.node, p.ID, input.ID+":"+strconv.FormatUint(input.Stream, 10))
		}
	}

	for state, value := range states {
		ch <- prometheus.MustNewConstMetric(c.processDesc, prometheus.GaugeValue, float64(value), c.node, state)
	}
}

type clusterFilesCollector struct {
	client  coreclient.RestClient
	node    string
	storage string

	filesDesc           *prometheus.Desc
	filesCollisionsDesc *prometheus.Desc
}

func newClusterFilesCollector(client coreclient.RestClient, node, storage string) prometheus.Collector {
	return &clusterFilesCollector{
		client:  client,
		node:    node,
		storage: storage,
		filesDesc: prometheus.NewDesc(
			"cluster_files",
			"Cluster number of files by storage",
			[]string{"node", "storage"}, nil),
		filesCollisionsDesc: prometheus.NewDesc(
			"cluster_files_collisions",
			"Cluster number of file collisions by storage",
			[]string{"node", "storage"}, nil),
	}
}

func (c *clusterFilesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.filesDesc
	ch <- c.filesCollisionsDesc
}

func (c *clusterFilesCollector) Collect(ch chan<- prometheus.Metric) {
	files, err := c.client.ClusterFilesystemList(c.storage, "/**", "asc", "name")
	if err != nil {
		return
	}

	ch <- prometheus.MustNewConstMetric(c.filesDesc, prometheus.GaugeValue, float64(len(files)), c.node, c.storage)

	collisions := map[string]uint64{}

	for _, file := range files {
		collisions[file.Name]++
	}

	ncollisions := float64(0)

	for _, counter := range collisions {
		if counter == 1 {
			continue
		}

		ncollisions += 1
	}

	ch <- prometheus.MustNewConstMetric(c.filesCollisionsDesc, prometheus.GaugeValue, ncollisions, c.node, c.storage)
}

var clusterExporterCmd = &cobra.Command{
	Use:   "exporter [clustername] [listenaddress]",
	Short: "Cluster exporter related commands",
	Long:  "Cluster exporter related commands",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		address, _ := cmd.Flags().GetString("address")

		var client coreclient.RestClient
		var err error

		if len(address) == 0 {
			client, err = connectSelectedCore()
			if err != nil {
				return err
			}
		} else {
			client, err = coreclient.New(coreclient.Config{
				Address: address,
			})
			if err != nil {
				return fmt.Errorf("can't connect to core at %s: %w", address, err)
			}
		}

		about, err := client.Cluster()
		if err != nil {
			return err
		}

		nodeCollector := newClusterNodeCollector(client)
		sessionCollector := newClusterHLSSessionCollector(client, about.ID)
		processCollector := newClusterProcessCollector(client, about.ID)
		filesCollector := newClusterFilesCollector(client, about.ID, "mem")

		registry := prometheus.NewRegistry()

		registry.Register(nodeCollector)
		registry.Register(sessionCollector)
		registry.Register(processCollector)
		registry.Register(filesCollector)

		http.Handle("/metrics", promhttp.InstrumentMetricHandler(registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

		quit := make(chan os.Signal, 1)

		go func() {
			if err := http.ListenAndServe(args[1], nil); err != nil && err != http.ErrServerClosed {
				if proc, err := os.FindProcess(os.Getpid()); err != nil {
					proc.Signal(os.Interrupt)
				}
			}
		}()

		signal.Notify(quit, os.Interrupt)
		<-quit

		return err
	},
}

func init() {
	clusterCmd.AddCommand(clusterExporterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	clusterExporterCmd.Flags().StringP("address", "a", "", "Alternative address for Core")
}
