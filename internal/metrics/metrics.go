package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(HttpResponseTime, OrderReceptionsCreatedTotal,
		PvzCreatedTotal, RequestsTotal, ProductsAddedTotal)
}

var (
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
	)

	HttpResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 1.5, 2, 3, 5},
		},
		[]string{"method", "path"},
	)

	PvzCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_pvz_created_total",
			Help: "Total number of PVZ (Пункты Выдачи Заказов) created",
		},
	)

	OrderReceptionsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_order_receptions_created_total",
			Help: "Total number of order receptions (Приёмки заказов) created",
		},
	)

	ProductsAddedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "business_products_added_total",
			Help: "Total number of products (Товары) added",
		},
	)
)