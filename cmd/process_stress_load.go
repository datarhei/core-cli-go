package cmd

import (
	"bufio"
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var processStressLoadCmd = &cobra.Command{
	Use:   "load [threads] [pause in milliseconds]",
	Short: "Process API stress",
	Long:  "Process API stress",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if pause < 0 {
			return fmt.Errorf("the pasue must be a positive number")
		}

		maxjitter := int(float64(pause) * 0.9)
		if maxjitter < 10 {
			maxjitter = 10
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
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

					list, err := client.ProcessList(coreclient.ProcessListOptions{
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

					client.ProcessAdd(*config)
					config.Reference = StringAlphanumeric(16)
					client.ProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), *config, false)
					client.Process(coreclient.NewProcessID(config.ID, config.Domain), []string{})

					requests <- 3
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
	processStressCmd.AddCommand(processStressLoadCmd)
}
