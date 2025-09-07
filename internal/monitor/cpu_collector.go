package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct {
	interval time.Duration
}

func NewCPUCollector(interval time.Duration) *CPUCollector {
	return &CPUCollector{interval: interval}
}

func (c *CPUCollector) Collect() (string, interface{}) {
	percent, err := cpu.Percent(c.interval, false)
	if err != nil || len(percent) == 0 {
		return "cpu", 0.0
	}
	return "cpu", percent[0]
}
