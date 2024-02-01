package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func (p *PrometheusMetrics) IncHits(status int, method, path string) {
	p.HitsTotal.Inc()
	p.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (p *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	p.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

func CreateMetric(address, name string) (Metrics, error) {
	var prometheusMetric PrometheusMetrics
	prometheusMetric.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(prometheusMetric.HitsTotal); err != nil {
		return nil, err
	}

	prometheusMetric.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_hits",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(prometheusMetric.Hits); err != nil {
		return nil, err
	}

	prometheusMetric.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name + "_times",
	}, []string{"status", "method", "path"})
	if err := prometheus.Register(prometheusMetric.Times); err != nil {
		return nil, err
	}

	return &prometheusMetric, nil
}
