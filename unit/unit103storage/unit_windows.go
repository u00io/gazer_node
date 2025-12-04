package unit103storage

import (
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

func (c *Unit103Storage) bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}

func (c *Unit103Storage) drives() []string {
	drives := make([]string, 0)
	drivesBits, err := windows.GetLogicalDrives()
	if err == nil {
		drives = c.bitsToDrives(drivesBits)
	}
	return drives
}

func (c *Unit103Storage) processTick() []item {
	result := make([]item, 0)
	drives := c.drives()
	var err error

	var TotalSpace uint64
	var UsedSpace uint64

	for _, disk := range drives {
		var free, total, avail uint64
		var diskName *uint16
		diskName, err = syscall.UTF16PtrFromString(disk + ":\\")

		err = windows.GetDiskFreeSpaceEx(
			diskName,
			&free,
			&total,
			&avail,
		)

		if err != nil {
			result = append(result, item{"/error", "Status", err.Error(), "error"})
		} else {
			result = append(result, item{"/", "Status", "", ""})
			result = append(result, item{"/" + disk + "/Total", "Total " + disk, strconv.FormatUint(total/1024/1024, 10), "MB"})
			result = append(result, item{"/" + disk + "/Free", "Free " + disk, strconv.FormatUint(free/1024/1024, 10), "MB"})
			result = append(result, item{"/" + disk + "/Used", "Used " + disk, strconv.FormatUint((total-free)/1024/1024, 10), "MB"})
			result = append(result, item{"/" + disk + "/Utilization", "Utilization " + disk, strconv.FormatFloat(100*float64(total-free)/float64(total), 'f', 2, 64), "%"})

			TotalSpace += total
			UsedSpace += total - free
		}
	}

	summaryUtilization := strconv.FormatFloat(100*float64(UsedSpace)/float64(TotalSpace), 'f', 1, 64)
	result = append(result, item{"/", "Used Percents", summaryUtilization, "%"})

	return result
}
