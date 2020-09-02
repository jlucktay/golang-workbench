package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
)

const (
	intervals   = 10
	processName = "gopls"

	msToS = 1000
)

func main() {
	usages := map[int][]float64{}

	for {
		procs, err := ps.Processes()
		if err != nil {
			log.Printf("error getting processes: %v", err)

			return
		}

		for _, proc := range procs {
			if !strings.EqualFold(proc.Executable(), processName) {
				continue
			}

			pid := proc.Pid()

			psu, errNew := process.NewProcess(int32(pid))
			if errNew != nil {
				log.Printf("error with new process: %v", errNew)

				return
			}

			procCreate, errCreateTime := psu.CreateTime()
			if errCreateTime != nil {
				log.Printf("error getting process create time: %v", errCreateTime)

				return
			}

			procStart := time.Unix(procCreate/msToS, 0)

			if time.Since(procStart) < 3*intervals*time.Second { // give new processes a chance to warm up
				fmt.Printf("%s [%d] is new (started %v ago at %v)\n", processName, pid,
					time.Since(procStart).Truncate(time.Second), procStart)

				continue
			}

			perc, errPerc := psu.CPUPercent()
			if errPerc != nil {
				log.Printf("error with process CPU percent: %v", errPerc)

				return
			}

			if usages[pid] == nil {
				usages[pid] = []float64{}
			}

			usages[pid] = append(usages[pid], perc)

			if len(usages[pid]) > intervals {
				usages[pid] = usages[pid][1:] // pop one off the front
			}

			fmt.Printf("%s [%d]: %v\n", processName, pid, usages[pid])
		}

		for pid := range usages {
			if len(usages[pid]) < intervals {
				continue
			}

			sum := float64(0)

			for _, usage := range usages[pid] {
				sum += usage
			}

			if sum >= intervals*50 {
				osp, errFind := os.FindProcess(pid)
				if errFind != nil {
					log.Printf("error finding process %d: %v", pid, errFind)

					continue
				}

				fmt.Printf("%s [%d] will be killed... ", processName, pid)

				if errKill := osp.Kill(); errKill != nil {
					log.Printf("error killing process %d: %v", pid, errFind)

					continue
				}

				fmt.Println("done.")

				delete(usages, pid)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
