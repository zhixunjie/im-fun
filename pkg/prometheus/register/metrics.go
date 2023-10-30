package register

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var M *Metrics

func init() {
	M = NewMetrics()
}

type Metrics struct {
	MetricsCollectors

	// registerer
	R prometheus.Registerer
	G prometheus.Gatherer
}

type MetricsCollectors struct {
	TestCollectors

	// 默认的收集器
	ProcessCollector prometheus.Collector
	GoCollector      prometheus.Collector
}

type TestCollectors struct {
	// 网络请求的收集器
	// refer：https://go-kratos.dev/docs/component/middleware/metrics/
	ReqCounterVec   *prometheus.CounterVec
	ReqHistogramVec *prometheus.HistogramVec

	// 时间戳瞬时值
	TsGaugeVec *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	m := new(Metrics)
	m.newRegistry()

	return m
}

// 创建注册器
func (m *Metrics) newRegistry() {
	// init
	registry := prometheus.NewRegistry()
	m.R = registry
	m.G = registry

	// create collector: 默认的收集器
	m.ProcessCollector = collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	m.GoCollector = collectors.NewGoCollector()

	// 注册默认的收集器
	m.R.MustRegister(m.ProcessCollector, m.GoCollector)
}
