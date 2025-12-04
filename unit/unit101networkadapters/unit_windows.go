package unit101networkadapters

import (
	"net"
	"strconv"
	"time"

	"github.com/kbinani/win"
)

func (c *Unit101NetworkAdapters) processTick() []item {
	var err error
	result := make([]item, 0)
	totalSpeed := 0.0
	totalInSpeed := 0.0
	totalOutSpeed := 0.0
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()

	if err != nil {
		result = append(result, item{"/", "Status", err.Error(), "error"})
		return result
	}

	for _, ni := range interfaces {
		c.writeAddresses(ni)

		var table win.MIB_IF_ROW2
		table.InterfaceIndex = win.NET_IFINDEX(ni.Index)
		win.GetIfEntry2(&table)

		if table.Type == 24 {
			continue
		}

		totalIn := uint64(table.InUcastPkts) + uint64(table.InNUcastPkts) + uint64(table.InDiscards)
		totalInBytes := uint64(table.InOctets)
		totalOut := uint64(table.OutUcastPkts) + uint64(table.OutNUcastPkts) + uint64(table.OutDiscards)
		totalOutBytes := uint64(table.OutOctets)

		nowTime := time.Now()

		if table.OperStatus == 1 {
			if cs, ok := c.lastCounters[ni.Index]; ok {
				seconds := nowTime.Sub(cs.DT).Seconds()
				if seconds > 0.001 {
					result = append(result, item{"/" + ni.Name + "/InSpeed", "In Speed " + ni.Name, strconv.FormatFloat(float64(totalInBytes-cs.TotalInBytes)/seconds/1024.0, 'f', 2, 64), "KB/sec"})
					result = append(result, item{"/" + ni.Name + "/OutSpeed", "Out Speed " + ni.Name, strconv.FormatFloat(float64(totalOutBytes-cs.TotalOutBytes)/seconds/1024.0, 'f', 2, 64), "KB/sec"})
					totalInSpeed += float64(totalInBytes-cs.TotalInBytes) / seconds / 1024.0
					totalOutSpeed += float64(totalOutBytes-cs.TotalOutBytes) / seconds / 1024.0
				}
			}

			c.lastCounters[ni.Index] = LastCounters{
				DT:            nowTime,
				TotalIn:       totalIn,
				TotalOut:      totalOut,
				TotalInBytes:  totalInBytes,
				TotalOutBytes: totalOutBytes,
			}
		} else {
			delete(c.lastCounters, ni.Index)
			result = append(result, item{"/" + ni.Name + "/InSpeed", "In Speed " + ni.Name, "0", "KB/sec"})
			result = append(result, item{"/" + ni.Name + "/OutSpeed", "Out Speed " + ni.Name, "0", "KB/sec"})
		}

	}

	totalSpeed = totalInSpeed + totalOutSpeed

	result = append(result, item{"/TotalInSpeed", "Total In Speed", strconv.FormatFloat(totalInSpeed, 'f', 2, 64), "KB/sec"})
	result = append(result, item{"/TotalOutSpeed", "Total Out Speed", strconv.FormatFloat(totalOutSpeed, 'f', 2, 64), "KB/sec"})
	result = append(result, item{"/TotalSpeed", "Total Speed", strconv.FormatFloat(totalSpeed, 'f', 2, 64), "KB/sec"})
	result = append(result, item{"/", "Total Speed", strconv.FormatFloat(totalSpeed, 'f', 2, 64), "KB/sec"})

	return result
}
