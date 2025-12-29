package metrics

type Metric struct {
	Value float64
}

type UploadMetricsResponse struct {
	Sum     float64
	Average float64
	Count   int64
}
