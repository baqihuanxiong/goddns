package metrics

type Metrics interface {
	CollectHttp(handler, method, statusCode string, duration float64)
	IncreaseDDNSCounter()
	Serve() error
}
