package unit301serialportkeyvalue

import (
	"strings"
	"time"

	"github.com/tarm/serial"
	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit301SerialPortKeyValue struct {
	unit000base.Unit

	serialConfig *serial.Config
	serialPort   *serial.Port
	inputBuffer  []byte
}

func New() unit000base.IUnit {
	var c Unit301SerialPortKeyValue
	c.SetType("unit301serialportkeyvalue")
	c.Init(&c)

	c.Config().SetParameterString("0000_00_name_str", "Serial Port Key=Value")

	c.serialConfig = &serial.Config{
		Name:        "COM3",
		Baud:        int(9600),
		ReadTimeout: 100 * time.Millisecond,
		Size:        byte(8),
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	return &c
}

func (c *Unit301SerialPortKeyValue) Tick() {
	var err error

	if c.serialPort == nil {
		c.serialPort, err = serial.OpenPort(c.serialConfig)
		if err != nil {
			c.serialPort = nil
			c.SetValue("/status", "Status", err.Error(), "error")
		} else {
			c.SetValue("/status", "Status", "connected", "")
		}
	}

	if c.serialPort != nil {
		buffer := make([]byte, 32)
		n, err := c.serialPort.Read(buffer)
		if err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "eof") {
				c.serialPort.Close()
				c.serialPort = nil
				c.SetValue("/status", "Status", err.Error(), "error")
			}
		} else {
			if n > 0 {
				c.inputBuffer = append(c.inputBuffer, buffer[:n]...)

				found := true
				for found {
					found = false
					currentLine := make([]byte, 0)
					for index, b := range c.inputBuffer {
						if b == 10 || b == 13 {
							// parse currentLine
							if len(currentLine) > 0 {
								parts := strings.Split(string(currentLine), "=")
								if len(parts) > 1 {
									if len(parts[0]) > 0 {
										key := parts[0]
										value := parts[1]

										finalValue := value
										c.SetValue("/"+key, key, finalValue, "")
										c.SetValue("/", key, finalValue, "")

										time.Sleep(100 * time.Microsecond)
									}
								}

							}
							c.inputBuffer = c.inputBuffer[index+1:]
							found = true
							break
						} else {
							if b >= 32 && b < 128 {
								currentLine = append(currentLine, b)
							}
						}
					}
				}
			}
		}
	}
}
