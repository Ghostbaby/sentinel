package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/ghostbaby/sentinel/pkg/collector"
	"github.com/ghostbaby/sentinel/pkg/g"

	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	cfg := flag.String("config", "", "configuration file")
	listen := flag.String("listen", ":80", "Address to listen on for web endpoints.")
	path := flag.String("path", "/metrics",
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	flag.Parse()
	g.ParseConfig(*cfg, false)

	log := ctrl.Log.WithName("controllers").WithName("sentinel")

	go g.ConfigReload(log)

	log.Info("Starting Prometheus Exporter", "Version", 0.1)

	registry := prometheus.NewRegistry()

	registry.MustRegister(collector.NewCollector("sentinel"))

	http.Handle(*path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Sentinel Exporter</title></head>
			<body>
			<h1>Sentinel Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	go http.ListenAndServe(*listen, nil)

	select {}
}
