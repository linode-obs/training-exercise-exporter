# training-exercise-exporter

This exporter provides two metrics that can be scraped by Prometheus:
* The current temperature in Boston, MA (in Celsius)
* The number of times it has scraped the NWS API since it was started.

This was written as reference code to accompany a presentation on how to write an exporter for Prometheus.

Additional resources used in training:
* https://prometheus.io/docs/instrumenting/writing_exporters/
* https://prometheus.io/docs/instrumenting/clientlibs/
* https://pkg.go.dev/github.com/prometheus/client_golang/prometheus
* https://github.com/prometheus/prometheus/wiki/Default-port-allocations