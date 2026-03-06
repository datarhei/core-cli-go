package cmd

import (
	"fmt"
	"sort"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var clusterProcessRebalanceCmd = &cobra.Command{
	Use:   "rebalance",
	Short: "Rebalance the processes in the cluster",
	Long:  "Rebalance the processes in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		execute, _ := cmd.Flags().GetBool("execute")
		prio, _ := cmd.Flags().GetString("prio")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		aboutv1, aboutv2, err := client.Cluster()
		if err != nil {
			return err
		}

		var about api.ClusterAbout

		if aboutv1 != nil {
			about = aboutv1.ClusterAbout
		} else {
			about = aboutv2.ClusterAbout
		}

		eligibleNodes := []api.ClusterNode{}

		for _, n := range about.Nodes {
			if n.Status != "online" {
				continue
			}

			if n.Resources.IsThrottling {
				continue
			}

			eligibleNodes = append(eligibleNodes, n)
		}

		nEligibleNodes := len(eligibleNodes)

		list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
			Filter: []string{"state"},
		})
		if err != nil {
			return err
		}

		type relocateProcess struct {
			id       coreclient.ProcessID
			fromNode string
			toNode   string
		}

		relocateList := []relocateProcess{}

		// Rebalance processes WITHOUT reference

		processListWithoutReference := []api.Process{}
		nodeProcessCount := map[string]int{}
		nodeCPU := map[string]float64{}
		totalCPU := float64(0)

		for _, p := range list {
			if len(p.Reference) != 0 {
				continue
			}

			processListWithoutReference = append(processListWithoutReference, p)
			nodeProcessCount[p.CoreID]++
			nodeCPU[p.CoreID] += p.State.Resources.CPU.Current
			totalCPU += p.State.Resources.CPU.Current
		}

		// Sort processes by runtime, shortest to longest
		sort.SliceStable(processListWithoutReference, func(a, b int) bool {
			return processListWithoutReference[a].State.Runtime < processListWithoutReference[b].State.Runtime
		})

		// Group processes by node
		processListNode := map[string][]api.Process{}

		for _, p := range processListWithoutReference {
			list := processListNode[p.CoreID]
			list = append(list, p)
			processListNode[p.CoreID] = list
		}

		nProcesses := len(processListWithoutReference)         // Number of processes
		nProcessesPerNode := (nProcesses / nEligibleNodes) + 1 // Desired number of processes per node
		nCPUPerNode := (totalCPU / float64(nEligibleNodes))    // Desired CPU per node

		if prio == "cpu" {
			for nodeid, cpu := range nodeCPU {
				diff := cpu - nCPUPerNode
				if diff <= 0 {
					continue
				}

				// This node has too many processes, move some away
				for _, p := range processListNode[nodeid] {
					relocateList = append(relocateList, relocateProcess{
						id:       coreclient.NewProcessID(p.ID, p.Domain),
						fromNode: p.CoreID,
					})

					diff -= p.State.Resources.CPU.Current
					if diff <= 0 {
						break
					}
				}
			}
		} else {
			// Redistribute processes
			for nodeid, count := range nodeProcessCount {
				diff := count - nProcessesPerNode
				if diff <= 0 {
					continue
				}

				// This node has too many processes, move some away
				for _, p := range processListNode[nodeid] {
					relocateList = append(relocateList, relocateProcess{
						id:       coreclient.NewProcessID(p.ID, p.Domain),
						fromNode: p.CoreID,
					})

					diff--
					if diff <= 0 {
						break
					}
				}
			}
		}

		// Rebalance processes WITH reference

		processReferenceMap := map[string][]api.Process{} // List of processes grouped by reference
		nodeProcessCount = map[string]int{}               // Number of processes per node
		nodeCPU = map[string]float64{}
		totalCPU = float64(0)

		// Group processes by their reference
		for _, p := range list {
			if len(p.Reference) == 0 {
				continue
			}

			ref := processReferenceMap[p.Reference]
			ref = append(ref, p)
			processReferenceMap[p.Reference] = ref

			nodeProcessCount[p.CoreID]++
			nodeCPU[p.CoreID] += p.State.Resources.CPU.Current

			totalCPU += p.State.Resources.CPU.Current
		}

		processListWithReference := []api.Process{}
		nProcesses = 0

		// Sort list of processes grouped by reference by their runtime, longest first, because this is the main process.
		// The first member is the representant of that reference group and we put in the processReferenceList
		for key, list := range processReferenceMap {
			sort.SliceStable(list, func(a, b int) bool {
				return list[a].State.Runtime > list[b].State.Runtime
			})
			processReferenceMap[key] = list

			processListWithReference = append(processListWithReference, list[0])
			nProcesses += len(list)
		}

		// Sort processes by runtime, shortest to longest
		sort.SliceStable(processListWithReference, func(a, b int) bool {
			if prio == "cpu" {
				return processListWithReference[a].State.Resources.CPU.Current > processListWithReference[b].State.Resources.CPU.Current
			}
			return processListWithReference[a].State.Runtime < processListWithReference[b].State.Runtime
		})

		// Group processed by node
		processListNode = map[string][]api.Process{}

		for _, p := range processListWithReference {
			list := processListNode[p.CoreID]
			list = append(list, p)
			processListNode[p.CoreID] = list
		}

		// Calculate desired number of processes per node
		nProcessesPerNode = (nProcesses / nEligibleNodes) + 1
		nCPUPerNode = (totalCPU / float64(nEligibleNodes)) // Desired CPU per node

		if prio == "cpu" {
			for nodeid, cpu := range nodeCPU {
				diff := cpu - nCPUPerNode
				if diff <= 0 {
					continue
				}

				// This node has too many processes, move some away. Here we have to
				// move all processes with the same reference
				for _, p := range processListNode[nodeid] {
					reference := p.Reference

					for _, p := range processReferenceMap[reference] {
						relocateList = append(relocateList, relocateProcess{
							id:       coreclient.NewProcessID(p.ID, p.Domain),
							fromNode: p.CoreID,
						})

						diff -= p.State.Resources.CPU.Current
					}

					if diff <= 0 {
						break
					}
				}
			}
		} else {
			for nodeid, count := range nodeProcessCount {
				diff := count - nProcessesPerNode
				if diff <= 0 {
					continue
				}

				// This node has too many processes, move some away. Here we have to
				// move all processes with the same reference
				for _, p := range processListNode[nodeid] {
					reference := p.Reference

					for _, p := range processReferenceMap[reference] {
						relocateList = append(relocateList, relocateProcess{
							id:       coreclient.NewProcessID(p.ID, p.Domain),
							fromNode: p.CoreID,
						})

						diff--
					}

					if diff <= 0 {
						break
					}
				}
			}
		}

		total := len(relocateList)

		if total == 0 {
			fmt.Printf("nothing to rebalance\n")
			return nil
		}

		for i, p := range relocateList {
			fmt.Printf("%4d / %4d: relocating %s away from %s ... ", i+1, total, p.id, p.fromNode)

			if execute {
				if err := client.ClusterRelocateProcess(p.id, p.toNode); err != nil {
					fmt.Printf("failed: %s\n", err.Error())
				} else {
					fmt.Printf("OK\n")
				}
			} else {
				fmt.Printf("demo\n")
			}
		}

		if !execute {
			fmt.Printf("Use -x to actually rebalance the processes.\n")
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessRebalanceCmd)

	clusterProcessRebalanceCmd.Flags().BoolP("execute", "x", false, "Actually execute the cleanup")
	clusterProcessRebalanceCmd.Flags().StringP("prio", "p", "", "Rebalance prio count|cpu|memory")
}
