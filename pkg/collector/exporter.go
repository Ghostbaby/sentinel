package collector

import (
	"github.com/ghostbaby/sentinel/pkg/work"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector collects HARBOR metrics. It implements prometheus.Collector interface.
type Collector struct {
	metrics map[string]*prometheus.Desc
}

// NewHarborCollector creates an HarborCollector.
func NewCollector(namespace string) *Collector {
	return &Collector{
		metrics: map[string]*prometheus.Desc{
			"sentinel_network_latency_max": newGlobalMetric(namespace, "sentinel_network_latency_max",
				"host network max latency", []string{"host", "url"}),
			"sentinel_network_latency_min": newGlobalMetric(namespace, "sentinel_network_latency_min",
				"host network min latency", []string{"host", "url"}),
			"sentinel_network_latency_avg": newGlobalMetric(namespace, "sentinel_network_latency_avg",
				"host network avg latency", []string{"host", "url"}),
		},
	}
}

// Describe sends the super-set of all possible descriptors of HARBOR metrics
// to the provided channel.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from HARBOR and sends them to the provided channel.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	hosts := work.PingCn()

	for _, v := range hosts {

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_max"],
			prometheus.GaugeValue, v.Max, v.Host, v.Name)

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_min"],
			prometheus.GaugeValue, v.Min, v.Host, v.Name)

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_avg"],
			prometheus.GaugeValue, v.Avg, v.Host, v.Name)
	}
}

func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}
