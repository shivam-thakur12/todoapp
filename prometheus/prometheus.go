package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var CREATE_REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "todo_app_create_requests_count",
	Help: "Total App create HTTP Requests count.",
})
var GET_REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "todo_app_get_requests_count",
	Help: "Total App get HTTP Requests count.",
})
var PATCH_REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "todo_app_patch_requests_count",
	Help: "Total App patch HTTP Requests count.",
})
var DELETE_REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "todo_app_delete_requests_count",
	Help: "Total App delete HTTP Requests count.",
})
