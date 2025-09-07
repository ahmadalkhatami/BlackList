package monitor

import "runtime"

// ThreadCollector menghitung jumlah goroutine aktif.
type ThreadCollector struct{}

func NewThreadCollector() *ThreadCollector {
	return &ThreadCollector{}
}

func (t *ThreadCollector) Collect() (string, interface{}) {
	return "threads", runtime.NumGoroutine()
}
