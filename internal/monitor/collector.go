package monitor

type Collector interface {
	Collect() (string, interface{})
}
