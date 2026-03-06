package cmd

import (
	"bufio"
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterProcessStressLoadCmd = &cobra.Command{
	Use:   "load [threads] [pause in milliseconds] [core[,core[,..]]]",
	Short: "Process API stress",
	Long:  "Process API stress",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		restart, _ := cmd.Flags().GetBool("restart")
		nThreads, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		if nThreads <= 0 {
			return fmt.Errorf("number or threads must be greater than 0")
		}

		pause, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		cores := strings.Split(args[2], ",")

		if pause < 0 {
			return fmt.Errorf("the pasue must be a positive number")
		}

		maxjitter := int(float64(pause) * 0.9)
		if maxjitter < 10 {
			maxjitter = 10
		}

		clients := []coreclient.RestClient{}

		for _, core := range cores {
			client, err := connectCore(core)
			if err != nil {
				fmt.Printf("%s: %s\n", core, err.Error())
				continue
			}

			clients = append(clients, client)
		}

		if len(clients) == 0 {
			return fmt.Errorf("no clients available")
		}

		requestsChan := make(chan int, nThreads)
		requestsTotal := uint64(0)

		go func(requests <-chan int) {
			for i := range requests {
				requestsTotal += uint64(i)
			}
		}(requestsChan)

		start := time.Now()

		ctx, cancel := context.WithCancel(context.Background())

		wg := sync.WaitGroup{}

		for i := 0; i < nThreads; i++ {
			wg.Add(1)

			go func(ctx context.Context, requests chan<- int) {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					jitter := rand.IntN(maxjitter) - maxjitter/2
					time.Sleep(time.Duration(pause+jitter) * time.Millisecond)

					i := rand.IntN(len(clients))
					client := clients[i]

					list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
						Filter:    []string{"state", "config"},
						IDPattern: "test_*",
					})
					if err != nil {
						continue
					}

					requests <- 1

					if len(list) == 0 {
						continue
					}

					i = rand.IntN(len(list))

					process := list[i]
					config := process.Config

					client.ClusterProcessAdd(*config)
					config.Reference = StringAlphanumeric(16)
					client.ClusterProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), *config, false)
					client.ClusterProcess(coreclient.NewProcessID(config.ID, config.Domain), []string{})
					if restart {
						client.ClusterProcessCommand(coreclient.NewProcessID(config.ID, config.Domain), "restart")
					}

					requests <- 4
				}
			}(ctx, requestsChan)
		}

		fmt.Println("press 'q' to quit")

		bufio.NewReader(os.Stdin).ReadBytes('q')

		fmt.Println("canelling all threads")

		cancel()

		wg.Wait()

		close(requestsChan)

		duration := time.Since(start)

		fmt.Printf("%d requests in %s => %.3f requests/s\n", requestsTotal, duration, float64(requestsTotal)/duration.Seconds())

		return nil
	},
}

func init() {
	clusterProcessStressCmd.AddCommand(clusterProcessStressLoadCmd)

	clusterProcessStressLoadCmd.Flags().Bool("restart", false, "Whether to issue an additional restart")
}
