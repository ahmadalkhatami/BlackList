package monitor

import "github.com/shirou/gopsutil/v3/mem"

// MemoryCollector menghitung persentase memory yang digunakan.
type MemoryCollector struct{}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (m *MemoryCollector) Collect() (string, interface{}) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "memory", 0.0
	}
	return "memory", vmStat.UsedPercent
}
