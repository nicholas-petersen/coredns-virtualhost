package virtualhost

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	hostnameCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "virtualhost",
		Name:      "hostname_count",
		Help:      "The count of hostname.",
	}, []string{"hostname"})
)
