package collector

import (
	"github.com/ghostbaby/sentinel/pkg/work"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector collects HARBOR metrics. It implements prometheus.Collector interface.
type Collector struct {
	metrics map[string]*prometheus.Desc
}

const (
	LatencyMaxType = "max"
	LatencyMinType = "min"
	LatencyAvgType = "avg"
)

// NewHarborCollector creates an HarborCollector.
func NewCollector(namespace string) *Collector {
	return &Collector{
		metrics: map[string]*prometheus.Desc{
			"sentinel_network_latency_max": newGlobalMetric(namespace, "sentinel_network_latency_max",
				"host network max latency", []string{"host", "url", "name"}),
			"sentinel_network_latency_min": newGlobalMetric(namespace, "sentinel_network_latency_min",
				"host network min latency", []string{"host", "url", "name"}),
			"sentinel_network_latency_avg": newGlobalMetric(namespace, "sentinel_network_latency_avg",
				"host network avg latency", []string{"host", "url", "name"}),
			"sentinel_network_latency_loss": newGlobalMetric(namespace, "sentinel_network_latency_loss",
				"host network loss", []string{"host", "url", "name"}),

			"sentinel_network_latency_details": newGlobalMetric(namespace, "sentinel_network_latency_details",
				"host network loss details", []string{"host", "url", "name", "node", "area", "isp"}),
			"sentinel_network_latency_details_min": newGlobalMetric(namespace, "sentinel_network_latency_details_min",
				"host network min details", []string{"host", "url", "name", "node", "area", "isp"}),
			"sentinel_network_latency_details_max": newGlobalMetric(namespace, "sentinel_network_latency_details_max",
				"host network max details", []string{"host", "url", "name", "node", "area", "isp"}),
			"sentinel_network_latency_details_avg": newGlobalMetric(namespace, "sentinel_network_latency_details_avg",
				"host network avg details", []string{"host", "url", "name", "node", "area", "isp"}),

			"sentinel_network_scp_result": newGlobalMetric(namespace, "sentinel_network_scp_result",
				"host network avg details", []string{"host", "name"}),
			//"sentinel_network_latency": newGlobalMetric(namespace, "sentinel_network_latency",
			//	"host network all latency", []string{"host", "url", "type", "name"}),
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

	pingCnResult := work.PingCn()

	for _, v := range pingCnResult {

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_max"],
			prometheus.GaugeValue, v.Max, v.Host, v.Provider, v.Name)

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_min"],
			prometheus.GaugeValue, v.Min, v.Host, v.Provider, v.Name)

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_avg"],
			prometheus.GaugeValue, v.Avg, v.Host, v.Provider, v.Name)

		ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_loss"],
			prometheus.GaugeValue, v.Loss, v.Host, v.Provider, v.Name)

		for _, detail := range v.Details {
			if detail.Loss == 0 {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_details"],
				prometheus.GaugeValue, detail.Loss, v.Host, v.Provider, v.Name, detail.Name, detail.Area, detail.IspName)

			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_details_min"],
				prometheus.GaugeValue, detail.Min, v.Host, v.Provider, v.Name, detail.Name, detail.Area, detail.IspName)
			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_details_max"],
				prometheus.GaugeValue, detail.Max, v.Host, v.Provider, v.Name, detail.Name, detail.Area, detail.IspName)
			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_latency_details_avg"],
				prometheus.GaugeValue, detail.Avg, v.Host, v.Provider, v.Name, detail.Name, detail.Area, detail.IspName)
		}
	}

	scpResult := work.Scp()
	for _, result := range scpResult {
		if result.IsReady {
			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_scp_result"],
				prometheus.GaugeValue, 1, result.IP, result.Name)
		} else {
			ch <- prometheus.MustNewConstMetric(c.metrics["sentinel_network_scp_result"],
				prometheus.GaugeValue, 0, result.IP, result.Name)
		}
	}
}

func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}
