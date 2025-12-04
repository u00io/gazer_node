package unit102process

import (
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/procfs"
)

func (c *Unit102Process) processTick() []item {
	var err error
	result := make([]item, 0)

	if c.actualProcessId == -1 {

		var err error

		allProcesses, err := procfs.AllProcs()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			return result
		}

		for _, p := range allProcesses {
			matchId := false
			matchName := false

			if c.configProcessIdActive {
				if int(c.configProcessId) == p.PID {
					matchId = true
				}
			} else {
				matchId = true
			}

			if c.configProcessNameActive {
				if comm, err := p.Comm(); err == nil && strings.Contains(comm, c.configProcessName) {
					matchName = true
				}
			} else {
				matchName = true
			}

			if matchId && matchName {
				c.actualProcessId = p.PID
				c.proc = p

				comm, err := p.Comm()
				if err == nil {
					result = append(result, item{"/Command", "Command", comm, ""})
				}
				exe, err := p.Executable()
				if err == nil {
					result = append(result, item{"/Executable", "Executable", exe, ""})
					//c.SetString("Executable", exe, "")
				}
				//c.SetFloat64("PID", float64(p.PID), "", 0)
				result = append(result, item{"/PID", "PID", strconv.Itoa(p.PID), ""})
			}
		}
	}

	if c.actualProcessId == -1 {
		time.Sleep(100 * time.Millisecond)
		{
			result = append(result, item{"/Status", "Status", "no process found", "error"})
		}
		return result
	}

	pStat, err := c.proc.Stat()
	if err == nil {
		tNow := time.Now()
		dur := tNow.Sub(c.lastCpuTime).Seconds()
		cpuTime := pStat.CPUTime()

		if c.lastCpuValid {
			value := cpuTime - c.lastCpuValue
			if dur > 0.0000001 {
				usage := value / dur
				result = append(result, item{"/CPU", "CPU", strconv.FormatFloat(usage*100, 'f', 2, 64), "%"})
			}
		}

		c.lastCpuTime = tNow
		c.lastCpuValue = cpuTime
		c.lastCpuValid = true

		result = append(result, item{"/ResidentMemory", "Resident Memory", strconv.FormatFloat(float64(pStat.ResidentMemory()/1024), 'f', 0, 64), "KB"})
		result = append(result, item{"/VirtualMemory", "Virtual Memory", strconv.FormatFloat(float64(pStat.VirtualMemory()/1024), 'f', 0, 64), "KB"})
		result = append(result, item{"/Status", "Status", pStat.State, ""})
	} else {
		result = append(result, item{"/Status", "Status", err.Error(), "error"})
		c.lastCpuValid = false
		c.actualProcessId = -1
	}

	fdInfo, err := c.proc.FileDescriptorsInfo()
	if err == nil {
		result = append(result, item{"/FileDescriptors", "File Descriptors", strconv.Itoa(fdInfo.Len()), ""})
	} else {
		result = append(result, item{"/FileDescriptors", "File Descriptors", "", "error"})
		c.actualProcessId = -1
		c.lastCpuValid = false
	}

	c.dtOperationTime = time.Now().UTC()

	return result
}
