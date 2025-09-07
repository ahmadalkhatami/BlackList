package monitor

import (
	"fmt"
	"time"
)

type MonitorService struct {
	interval   time.Duration
	collectors []Collector
}

type MonitorOption func(*MonitorService)

func WithCollector(c Collector) MonitorOption {
	return func(m *MonitorService) {
		m.collectors = append(m.collectors, c)
	}
}

func NewMonitorService(interval time.Duration, opts ...MonitorOption) *MonitorService {
	svc := &MonitorService{
		interval:   interval,
		collectors: []Collector{},
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (m *MonitorService) Start() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, collector := range m.collectors {
			name, value := collector.Collect()
			fmt.Printf("[%s] %v\n", name, value)
		}
	}
}
