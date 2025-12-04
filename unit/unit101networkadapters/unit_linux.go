package unit101networkadapters

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
		rxPackets := int64(0)
		rxBytes := int64(0)
		txPackets := int64(0)
		txBytes := int64(0)

		// Addresses
		addrs, err := ni.Addrs()
		if err != nil {
			addrsString := ""
			for _, a := range addrs {
				if len(addrsString) > 0 {
					addrsString += " "
				}
				addrsString += a.String()
			}
			//c.SetString(ni.Name+"/Addresses", addrsString, "-")
			result = append(result, item{"/" + ni.Name + "/Addresses", "Addresses " + ni.Name, addrsString, "-"})
		}

		rxPacketsStr, errParamRxPackets := os.ReadFile("/sys/class/net/" + ni.Name + "/statistics/rx_packets")
		if errParamRxPackets == nil {
			rxPackets, errParamRxPackets = strconv.ParseInt(strings.ReplaceAll(string(rxPacketsStr), "\n", ""), 10, 64)
		} else {
			// logger.Println(errParamRxPackets)
		}

		//c.SetString("errParamRxPackets", errParamRxPackets.Error(), "q")

		rxBytesStr, errParamRxBytes := os.ReadFile("/sys/class/net/" + ni.Name + "/statistics/rx_bytes")
		if errParamRxBytes == nil {
			rxBytes, errParamRxBytes = strconv.ParseInt(strings.ReplaceAll(string(rxBytesStr), "\n", ""), 10, 64)
		} else {
			// logger.Println(errParamRxBytes)
		}

		//c.SetString("errParamRxBytes", errParamRxBytes.Error(), "q")

		txPacketsStr, errParamTxPackets := os.ReadFile("/sys/class/net/" + ni.Name + "/statistics/tx_packets")
		if errParamTxPackets == nil {
			txPackets, errParamTxPackets = strconv.ParseInt(strings.ReplaceAll(string(txPacketsStr), "\n", ""), 10, 64)
		} else {
			// logger.Println(errParamTxPackets)
		}

		//c.SetString("errParamTxPackets", errParamTxPackets.Error(), "q")

		txBytesStr, errParamTxBytes := os.ReadFile("/sys/class/net/" + ni.Name + "/statistics/tx_bytes")
		if errParamTxBytes == nil {
			txBytes, errParamTxBytes = strconv.ParseInt(strings.ReplaceAll(string(txBytesStr), "\n", ""), 10, 64)
		} else {
			// logger.Println(errParamTxBytes)
		}

		//c.SetString("errParamTxBytes", errParamTxBytes.Error(), "q")

		totalIn := uint64(rxPackets)
		totalInBytes := uint64(rxBytes)
		totalOut := uint64(txPackets)
		totalOutBytes := uint64(txBytes)

		nowTime := time.Now()

		if true {
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
