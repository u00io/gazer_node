package unit102process

import (
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func (c *Unit102Process) processTick() []item {
	//var err error
	result := make([]item, 0)

	processJustFound := false

	if c.actualProcessId < 0 {
		if !c.configProcessIdActive && !c.configProcessNameActive {
			return result
		}

		handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
		if err == nil {
			var entry windows.ProcessEntry32
			entry.Size = uint32(unsafe.Sizeof(entry))
			err = windows.Process32First(handle, &entry)
			for err == nil {
				nameSize := 0
				for i := 0; i < 260; i++ {
					if entry.ExeFile[nameSize] == 0 {
						break
					}
					nameSize++
				}

				id := int(entry.ProcessID)
				name := syscall.UTF16ToString(entry.ExeFile[:nameSize])

				// Filtering
				matchId := false
				matchName := false
				if c.configProcessIdActive {
					if id == c.configProcessId {
						matchId = true
					}
				} else {
					matchId = true
				}
				if c.configProcessNameActive {
					if strings.ToLower(name) == strings.ToLower(c.configProcessName) {
						matchName = true
					}
				} else {
					matchName = true
				}
				if matchId && matchName {
					c.actualProcessId = int(entry.ProcessID)
					processJustFound = true
					c.actualProcessName = name
					break
				}
				// /////////////////

				err = windows.Process32Next(handle, &entry)
			}

			_ = windows.CloseHandle(handle)
		}
	}

	if c.actualProcessId >= 0 {
		hProcess, err := windows.OpenProcess(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|windows.SPECIFIC_RIGHTS_ALL, false, uint32(c.actualProcessId))
		if err == nil {
			// Common

			result = append(result, item{"/Common/Name", "Name", c.actualProcessName, ""})
			result = append(result, item{"/Common/ProcessID", "Process ID", strconv.Itoa(c.actualProcessId), ""})

			{
				res, _ := GetProcessMemoryInfo(hProcess)
				result = append(result, item{"/Memory/WorkingSetSize", "Working Set Size", strconv.FormatUint(res.WorkingSetSize/1024, 10), "KB"})
				result = append(result, item{"/Memory/PageFaults", "Page Faults", strconv.FormatUint(uint64(res.PageFaultCount), 10), ""})
				result = append(result, item{"/Memory/PeakWorkingSetSize", "Peak Working Set Size", strconv.FormatUint(res.PeakWorkingSetSize/1024, 10), "KB"})
				result = append(result, item{"/Memory/PrivateUsage", "Private Usage", strconv.FormatUint(res.PrivateUsage/1024, 10), "KB"})
			}

			result = append(result, item{"/Main/ThreadCount", "Thread Count", strconv.Itoa(ProcessThreadsCount(uint32(c.actualProcessId))), ""})
			result = append(result, item{"/Main/HandleCount", "Handle Count", strconv.Itoa(GetProcessHandleCount(hProcess)), ""})
			{
				cntGDI, cntUser, cntGDIPeak, cntUserPeak, _ := GetGuiResources(hProcess)

				result = append(result, item{"/Main/GDIObjects", "GDI Objects", strconv.FormatInt(cntGDI, 10), ""})
				result = append(result, item{"/Main/GDIObjectsPeak", "GDI Objects Peak", strconv.FormatInt(cntGDIPeak, 10), ""})
				result = append(result, item{"/Main/UserObjects", "User Objects", strconv.FormatInt(cntUser, 10), ""})
				result = append(result, item{"/Main/UserObjectsPeak", "User Objects Peak", strconv.FormatInt(cntUserPeak, 10), ""})
			}

			{
				var ftStart windows.Filetime
				var ftEnd windows.Filetime
				var ftKernel windows.Filetime
				var ftUser windows.Filetime
				err = windows.GetProcessTimes(hProcess, &ftStart, &ftEnd, &ftKernel, &ftUser)
				if err == nil {
					kernelTimeMs := (int64(ftKernel.HighDateTime)<<32 + int64(ftKernel.LowDateTime)) / 100000
					userTimeMs := (int64(ftUser.HighDateTime)<<32 + int64(ftUser.LowDateTime)) / 100000
					deltaKernel := kernelTimeMs - c.lastKernelTimeMs
					deltaUser := userTimeMs - c.lastUserTimeMs
					duringMs := time.Now().UTC().Sub(c.lastReadProcessTimes).Milliseconds()

					usageCpuKernel := float64(0)
					usageCpuUser := float64(0)
					usageCpu := float64(0)

					if duringMs > 0 {
						usageCpuKernel = float64(deltaKernel) / float64(duringMs)
						usageCpuUser = float64(deltaUser) / float64(duringMs)
						usageCpu = float64(deltaKernel+deltaUser) / float64(duringMs)
					}

					c.lastReadProcessTimes = time.Now().UTC()
					c.lastKernelTimeMs = kernelTimeMs
					c.lastUserTimeMs = userTimeMs

					if processJustFound {
						usageCpuKernel = 0
						usageCpuUser = 0
						usageCpu = 0
					}

					/*c.SetInt64("CPU/Kernel Mode Time", kernelTimeMs, uom.MS)
					c.SetInt64("CPU/User Mode Time", userTimeMs, uom.MS)
					c.SetFloat64("CPU/Usage", usageCpu*100, uom.PERCENTS, 1)
					c.SetFloat64("CPU/Usage Kernel", usageCpuKernel*100, uom.PERCENTS, 1)
					c.SetFloat64("CPU/Usage User", usageCpuUser*100, uom.PERCENTS, 1)*/

					result = append(result, item{"/CPU/KernelModeTime", "Kernel Mode Time", strconv.FormatInt(kernelTimeMs, 10), "ms"})
					result = append(result, item{"/CPU/UserModeTime", "User Mode Time", strconv.FormatInt(userTimeMs, 10), "ms"})
					result = append(result, item{"/CPU/Usage", "CPU Usage", strconv.FormatFloat(usageCpu*100, 'f', 2, 64), "%"})
					result = append(result, item{"/", "CPU Usage", strconv.FormatFloat(usageCpu*100, 'f', 2, 64), "%"})
					result = append(result, item{"/CPU/UsageKernel", "CPU Usage Kernel", strconv.FormatFloat(usageCpuKernel*100, 'f', 2, 64), "%"})
					result = append(result, item{"/CPU/UsageUser", "CPU Usage User", strconv.FormatFloat(usageCpuUser*100, 'f', 2, 64), "%"})

					{
						res, _ := GetProcessIoCounters(hProcess)
						result = append(result, item{"/IO/ReadOperationCount", "Read Operation Count", strconv.FormatUint(res.ReadOperationCount, 10), ""})
						result = append(result, item{"/IO/ReadTransferCount", "Read Transfer Count", strconv.FormatUint(res.ReadTransferCount, 10), ""})
						result = append(result, item{"/IO/WriteOperationCount", "Write Operation Count", strconv.FormatUint(res.WriteOperationCount, 10), ""})
						result = append(result, item{"/IO/WriteTransferCount", "Write Transfer Count", strconv.FormatUint(res.WriteTransferCount, 10), ""})
						result = append(result, item{"/IO/OtherOperationCount", "Other Operation Count", strconv.FormatUint(res.OtherOperationCount, 10), ""})
						result = append(result, item{"/IO/OtherTransferCount", "Other Transfer Count", strconv.FormatUint(res.OtherTransferCount, 10), ""})
					}

				}
			}

			//GetGuiResources function (winuser.h)

			_ = windows.CloseHandle(hProcess)
		} else {
			c.actualProcessId = -1
		}
	}

	if c.actualProcessId < 0 {
		result = append(result, item{"/Status", "Error", "No process found", ""})
	} else {
		result = append(result, item{"/Status", "Error", "", ""})
	}

	c.dtOperationTime = time.Now().UTC()

	return result
}

var (
	modpsapi    = windows.NewLazySystemDLL("psapi.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	moduser32   = windows.NewLazySystemDLL("user32.dll")
)

var getProcessMemoryInfo = modpsapi.NewProc("GetProcessMemoryInfo")
var getProcessHandleCount = modkernel32.NewProc("GetProcessHandleCount")
var getProcessIoCounters = modkernel32.NewProc("GetProcessIoCounters")
var getGuiResources = moduser32.NewProc("GetGuiResources")

type PROCESS_MEMORY_COUNTERS_EX struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uint64
	WorkingSetSize             uint64
	QuotaPeakPagedPoolUsage    uint64
	QuotaPagedPoolUsage        uint64
	QuotaPeakNonPagedPoolUsage uint64
	QuotaNonPagedPoolUsage     uint64
	PagefileUsage              uint64
	PeakPagefileUsage          uint64
	PrivateUsage               uint64
}

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	}
	return e
}

func GetProcessMemoryInfo(handle windows.Handle) (pc PROCESS_MEMORY_COUNTERS_EX, err error) {
	var res PROCESS_MEMORY_COUNTERS_EX
	r1, _, e1 := syscall.Syscall(getProcessMemoryInfo.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(&res)), uintptr(80))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return res, err
}

func ProcessThreadsCount(processId uint32) int {
	count := 0
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err == nil {
		var entry windows.ThreadEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))
		err := windows.Thread32First(handle, &entry)
		for err == nil {
			if entry.OwnerProcessID == processId {
				count++
			}
			entry.Size = uint32(unsafe.Sizeof(entry))
			err = windows.Thread32Next(handle, &entry)
		}

		_ = windows.CloseHandle(handle)
	}

	return count
}

func GetProcessHandleCount(handle windows.Handle) int {
	var res uint32
	syscall.Syscall(getProcessHandleCount.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(&res)), 0)
	return int(res)
}

func GetProcessIoCounters(handle windows.Handle) (cnt windows.IO_COUNTERS, err error) {
	var res windows.IO_COUNTERS
	r1, _, e1 := syscall.Syscall(getProcessIoCounters.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(&res)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return res, err
}

func GetGuiResources(handle windows.Handle) (cntGDI int64, cntUser int64, cntGDIPeak int64, cntUserPeak int64, err error) {
	var flags uint32
	flags = 0
	r1, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntGDI = int64(r1)
	flags = 1
	r2, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntUser = int64(r2)
	flags = 2
	r3, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntGDIPeak = int64(r3)
	flags = 4
	r4, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntUserPeak = int64(r4)
	return
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err == nil {
		var entry windows.ProcessEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))
		err = windows.Process32First(handle, &entry)
		for err == nil {
			nameSize := 0
			for i := 0; i < 260; i++ {
				if entry.ExeFile[nameSize] == 0 {
					break
				}
				nameSize++
			}
			name := syscall.UTF16ToString(entry.ExeFile[:nameSize])

			var pi ProcessInfo
			pi.Id = int(entry.ProcessID)
			pi.Name = name
			result = append(result, pi)

			err = windows.Process32Next(handle, &entry)
		}

		_ = windows.CloseHandle(handle)
	}

	return result
}

type ProcessInfo struct {
	Name string
	Id   int
	Info string
}
