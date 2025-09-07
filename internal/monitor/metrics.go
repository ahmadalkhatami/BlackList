package monitor

type Metric struct {
	Name  string
	Value float64
}

type MetricCollector interface {
	Collect() Metric
}

// type MetricHandler interface {
// 	Handle(Metric)
// }

type MonitorServiceInterface interface {
	Start()
}
