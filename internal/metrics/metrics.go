package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(HttpResponseTime, OrderReceptionsCreatedTotal,
		PvzCreatedTotal, RequestsTotal, ProductsAddedTotal)
}

var (
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Кол-во запросов",
		},
	)

	HttpResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_time_seconds",
			Help:    "Длительность запросов",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 1.5, 2, 3, 5},
		},
		[]string{"method", "path"},
	)

	PvzCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_pvz_created_total",
			Help: "Кол-во созданных ПВЗ",
		},
	)

	OrderReceptionsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_order_receptions_created_total",
			Help: "Кол-во созданных приемок",
		},
	)

	ProductsAddedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_products_added_total",
			Help: "Кол-во добавленных продуктов",
		},
	)
)