package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

type clusterHLSSessionCollector struct {
	client   coreclient.RestClient
	node     string
	tx_bytes uint64
	last     map[string]uint64

	hlsSessionsDesc      *prometheus.Desc
	hlsSessionsBytesDesc *prometheus.Desc
}

func newClusterHLSSessionCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterHLSSessionCollector{
		client:   client,
		node:     node,
		tx_bytes: 0,
		last:     map[string]uint64{},
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

	check := map[string]struct{}{}

	for id := range c.last {
		check[id] = struct{}{}
	}

	for _, sess := range sessions.Active.SessionList {
		c.tx_bytes += sess.TxBytes - c.last[sess.ID]
		c.last[sess.ID] = sess.TxBytes
		delete(check, sess.ID)
	}

	for id := range check {
		delete(c.last, id)
	}

	ch <- prometheus.MustNewConstMetric(c.hlsSessionsDesc, prometheus.GaugeValue, float64(len(sessions.Active.SessionList)), c.node)
	ch <- prometheus.MustNewConstMetric(c.hlsSessionsBytesDesc, prometheus.CounterValue, float64(c.tx_bytes), c.node)
}

type clusterNodeCollector struct {
	client coreclient.RestClient

	cpuLimitDesc      *prometheus.Desc
	cpuCurrentDesc    *prometheus.Desc
	cpuNCoresDesc     *prometheus.Desc
	cpuCoreDesc       *prometheus.Desc
	memLimitDesc      *prometheus.Desc
	memCurrentDesc    *prometheus.Desc
	memCoreDesc       *prometheus.Desc
	throttlingDesc    *prometheus.Desc
	degradedDesc      *prometheus.Desc
	gpuLimitDesc      *prometheus.Desc
	gpuGeneralDesc    *prometheus.Desc
	gpuDecoderDesc    *prometheus.Desc
	gpuEncoderDesc    *prometheus.Desc
	gpuMemLimitDesc   *prometheus.Desc
	gpuMemCurrentDesc *prometheus.Desc
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
		cpuNCoresDesc: prometheus.NewDesc(
			"cluster_node_cpu_cores",
			"Cluster node CPU cores",
			[]string{"node"}, nil),
		cpuCoreDesc: prometheus.NewDesc(
			"cluster_core_cpu_current_percent",
			"Cluster core CPU current in percent",
			[]string{"node"}, nil),
		memLimitDesc: prometheus.NewDesc(
			"cluster_node_mem_limit_bytes",
			"Cluster node memory limit in bytes",
			[]string{"node"}, nil),
		memCurrentDesc: prometheus.NewDesc(
			"cluster_node_mem_current_bytes",
			"Cluster node memory current in bytes",
			[]string{"node"}, nil),
		memCoreDesc: prometheus.NewDesc(
			"cluster_core_mem_current_bytes",
			"Cluster core memory current in bytes",
			[]string{"node"}, nil),
		throttlingDesc: prometheus.NewDesc(
			"cluster_node_throttling",
			"Cluster node throttling",
			[]string{"node"}, nil),
		degradedDesc: prometheus.NewDesc(
			"cluster_node_degraded",
			"Cluster node degraded",
			[]string{"node"}, nil),
		gpuLimitDesc: prometheus.NewDesc(
			"cluster_node_gpu_usage_limit_percent",
			"Cluster node GPU usage limit in percent per GPU",
			[]string{"node", "gpu"}, nil),
		gpuGeneralDesc: prometheus.NewDesc(
			"cluster_node_gpu_usage_general_percent",
			"Cluster node GPU general usage in percent per GPU",
			[]string{"node", "gpu"}, nil),
		gpuDecoderDesc: prometheus.NewDesc(
			"cluster_node_gpu_usage_decoder_percent",
			"Cluster node GPU decoder usage in percent per GPU",
			[]string{"node", "gpu"}, nil),
		gpuEncoderDesc: prometheus.NewDesc(
			"cluster_node_gpu_usage_encoder_percent",
			"Cluster node GPU encoder usage in percent per GPU",
			[]string{"node", "gpu"}, nil),
		gpuMemLimitDesc: prometheus.NewDesc(
			"cluster_node_gpu_mem_limit_bytes",
			"Cluster node GPU memory limit in bytes",
			[]string{"node", "gpu"}, nil),
		gpuMemCurrentDesc: prometheus.NewDesc(
			"cluster_node_gpu_mem_current_bytes",
			"Cluster node GPU current memory in bytes",
			[]string{"node", "gpu"}, nil),
	}
}

func (c *clusterNodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuLimitDesc
	ch <- c.cpuCurrentDesc
	ch <- c.cpuNCoresDesc
	ch <- c.cpuCoreDesc
	ch <- c.memLimitDesc
	ch <- c.memCurrentDesc
	ch <- c.memCoreDesc
	ch <- c.throttlingDesc
	ch <- c.degradedDesc
	ch <- c.gpuLimitDesc
	ch <- c.gpuGeneralDesc
	ch <- c.gpuDecoderDesc
	ch <- c.gpuEncoderDesc
	ch <- c.gpuMemLimitDesc
	ch <- c.gpuMemCurrentDesc
}

func (c *clusterNodeCollector) Collect(ch chan<- prometheus.Metric) {
	coreabout, _ := c.client.About(true)

	aboutv1, aboutv2, err := c.client.Cluster()
	if err != nil {
		return
	}

	var about api.ClusterAbout

	if aboutv1 != nil {
		about = aboutv1.ClusterAbout
	} else {
		about = aboutv2.ClusterAbout
	}

	for _, node := range about.Nodes {
		if node.ID != coreabout.ID {
			continue
		}

		throttling := .0
		if node.Resources.IsThrottling {
			throttling = 1.0
		}

		ch <- prometheus.MustNewConstMetric(c.cpuLimitDesc, prometheus.GaugeValue, node.Resources.CPULimit, node.ID)
		ch <- prometheus.MustNewConstMetric(c.cpuCurrentDesc, prometheus.GaugeValue, node.Resources.CPU, node.ID)
		ch <- prometheus.MustNewConstMetric(c.cpuNCoresDesc, prometheus.GaugeValue, node.Resources.NCPU, node.ID)
		ch <- prometheus.MustNewConstMetric(c.cpuCoreDesc, prometheus.GaugeValue, node.Resources.CPUCore, node.ID)
		ch <- prometheus.MustNewConstMetric(c.memLimitDesc, prometheus.GaugeValue, float64(node.Resources.MemLimit), node.ID)
		ch <- prometheus.MustNewConstMetric(c.memCurrentDesc, prometheus.GaugeValue, float64(node.Resources.Mem), node.ID)
		ch <- prometheus.MustNewConstMetric(c.memCoreDesc, prometheus.GaugeValue, float64(node.Resources.MemCore), node.ID)
		ch <- prometheus.MustNewConstMetric(c.throttlingDesc, prometheus.GaugeValue, throttling, node.ID)

		for i, gpu := range node.Resources.GPU {
			gpuid := strconv.Itoa(i)
			ch <- prometheus.MustNewConstMetric(c.gpuLimitDesc, prometheus.GaugeValue, gpu.UsageLimit, node.ID, gpuid)
			ch <- prometheus.MustNewConstMetric(c.gpuGeneralDesc, prometheus.GaugeValue, gpu.Usage, node.ID, gpuid)
			ch <- prometheus.MustNewConstMetric(c.gpuDecoderDesc, prometheus.GaugeValue, gpu.Decoder, node.ID, gpuid)
			ch <- prometheus.MustNewConstMetric(c.gpuEncoderDesc, prometheus.GaugeValue, gpu.Encoder, node.ID, gpuid)
			ch <- prometheus.MustNewConstMetric(c.gpuMemLimitDesc, prometheus.GaugeValue, float64(gpu.MemLimit), node.ID, gpuid)
			ch <- prometheus.MustNewConstMetric(c.gpuMemCurrentDesc, prometheus.GaugeValue, float64(gpu.Mem), node.ID, gpuid)
		}

		break
	}

	degraded := .0
	if about.Degraded {
		degraded = 1.0
	}

	ch <- prometheus.MustNewConstMetric(c.degradedDesc, prometheus.GaugeValue, degraded, coreabout.ID)
}

type clusterProcessCollector struct {
	client coreclient.RestClient
	node   string

	processDesc             *prometheus.Desc
	processResourcesCPUDesc *prometheus.Desc
	processResourcesMemDesc *prometheus.Desc
	processInputFPSDesc     *prometheus.Desc
	processInputSpeedDesc   *prometheus.Desc
}

func newClusterProcessCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterProcessCollector{
		client: client,
		node:   node,
		processDesc: prometheus.NewDesc(
			"cluster_process",
			"Cluster processes by state",
			[]string{"node", "state"}, nil),
		processResourcesCPUDesc: prometheus.NewDesc(
			"cluster_process_resources_cpu",
			"Cluster process CPU consumption by id",
			[]string{"node", "id"}, nil),
		processResourcesMemDesc: prometheus.NewDesc(
			"cluster_process_resources_mem",
			"Cluster process memory consumption by id",
			[]string{"node", "id"}, nil),
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
	ch <- c.processResourcesCPUDesc
	ch <- c.processResourcesMemDesc
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

		ch <- prometheus.MustNewConstMetric(c.processResourcesCPUDesc, prometheus.GaugeValue, p.State.Resources.CPU.Current, c.node, p.ID)
		ch <- prometheus.MustNewConstMetric(c.processResourcesMemDesc, prometheus.GaugeValue, float64(p.State.Resources.Memory.Current), c.node, p.ID)

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

type clusterHTTPStatusCollector struct {
	client coreclient.RestClient
	node   string

	statusDesc *prometheus.Desc
}

func newClusterHTTPStatusCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterHTTPStatusCollector{
		client: client,
		node:   node,
		statusDesc: prometheus.NewDesc(
			"cluster_http_request_count",
			"Cluster HTTP requests by status, method, and path",
			[]string{"node", "status", "method", "path"}, nil),
	}
}

func (c *clusterHTTPStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.statusDesc
}

func (c *clusterHTTPStatusCollector) Collect(ch chan<- prometheus.Metric) {
	query := api.MetricsQuery{
		Metrics: []api.MetricsQueryMetric{
			{
				Name:   "http_status",
				Labels: map[string]string{},
			},
		},
	}

	resp, err := c.client.Metrics(query)
	if err != nil {
		return
	}

	for _, m := range resp.Metrics {
		ch <- prometheus.MustNewConstMetric(c.statusDesc, prometheus.CounterValue, m.Values[0].Value, c.node, m.Labels["code"], m.Labels["method"], m.Labels["path"])
	}
}

type clusterBufferpoolCollector struct {
	client coreclient.RestClient
	node   string

	allocDesc       *prometheus.Desc
	reuseDesc       *prometheus.Desc
	recycleDesc     *prometheus.Desc
	dumpDesc        *prometheus.Desc
	defaultSizeDesc *prometheus.Desc
	maxSizeDesc     *prometheus.Desc
}

func newClusterBufferpoolCollector(client coreclient.RestClient, node string) prometheus.Collector {
	return &clusterBufferpoolCollector{
		client: client,
		node:   node,
		allocDesc: prometheus.NewDesc(
			"cluster_bufferpool_alloc",
			"Cluster bufferpool allocations",
			[]string{"node"}, nil),
		reuseDesc: prometheus.NewDesc(
			"cluster_bufferpool_reuse",
			"Cluster bufferpool reuses of an already allocated buffer",
			[]string{"node"}, nil),
		recycleDesc: prometheus.NewDesc(
			"cluster_bufferpool_recycle",
			"Cluster bufferpool recycling a buffer",
			[]string{"node"}, nil),
		dumpDesc: prometheus.NewDesc(
			"cluster_bufferpool_dump",
			"Cluster bufferpool throwing away a buffer",
			[]string{"node"}, nil),
		defaultSizeDesc: prometheus.NewDesc(
			"cluster_bufferpool_default_size_bytes",
			"Cluster bufferpool min. size of a buffer on allocation",
			[]string{"node"}, nil),
		maxSizeDesc: prometheus.NewDesc(
			"cluster_bufferpool_max_size_bytes",
			"Cluster bufferpool max. size of a buffer in order the get recycled",
			[]string{"node"}, nil),
	}
}

func (c *clusterBufferpoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.allocDesc
	ch <- c.reuseDesc
	ch <- c.recycleDesc
	ch <- c.dumpDesc
	ch <- c.defaultSizeDesc
	ch <- c.maxSizeDesc
}

func (c *clusterBufferpoolCollector) Collect(ch chan<- prometheus.Metric) {
	query := api.MetricsQuery{
		Metrics: []api.MetricsQueryMetric{
			{
				Name:   "self_bufferpool_alloc",
				Labels: map[string]string{},
			},
			{
				Name:   "self_bufferpool_reuse",
				Labels: map[string]string{},
			},
			{
				Name:   "self_bufferpool_recycle",
				Labels: map[string]string{},
			},
			{
				Name:   "self_bufferpool_dump",
				Labels: map[string]string{},
			},
			{
				Name:   "self_bufferpool_default_size",
				Labels: map[string]string{},
			},
			{
				Name:   "self_bufferpool_max_size",
				Labels: map[string]string{},
			},
		},
	}

	resp, err := c.client.Metrics(query)
	if err != nil {
		return
	}

	var desc *prometheus.Desc
	var vtype prometheus.ValueType

	for _, m := range resp.Metrics {
		switch m.Name {
		case "self_bufferpool_alloc":
			desc = c.allocDesc
			vtype = prometheus.CounterValue
		case "self_bufferpool_reuse":
			desc = c.reuseDesc
			vtype = prometheus.CounterValue
		case "self_bufferpool_recycle":
			desc = c.recycleDesc
			vtype = prometheus.CounterValue
		case "self_bufferpool_dump":
			desc = c.dumpDesc
			vtype = prometheus.CounterValue
		case "self_bufferpool_default_size":
			desc = c.defaultSizeDesc
			vtype = prometheus.GaugeValue
		case "self_bufferpool_max_size":
			desc = c.maxSizeDesc
			vtype = prometheus.GaugeValue
		default:
			continue
		}

		ch <- prometheus.MustNewConstMetric(desc, vtype, m.Values[0].Value, c.node)
	}
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

		coreabout, _ := client.About(true)

		nodeCollector := newClusterNodeCollector(client)
		sessionCollector := newClusterHLSSessionCollector(client, coreabout.ID)
		processCollector := newClusterProcessCollector(client, coreabout.ID)
		filesCollector := newClusterFilesCollector(client, coreabout.ID, "mem")
		statusCollector := newClusterHTTPStatusCollector(client, coreabout.ID)
		bufferpoolCollector := newClusterBufferpoolCollector(client, coreabout.ID)

		registry := prometheus.NewRegistry()

		registry.Register(nodeCollector)
		registry.Register(sessionCollector)
		registry.Register(processCollector)
		registry.Register(filesCollector)
		registry.Register(statusCollector)
		registry.Register(bufferpoolCollector)

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
