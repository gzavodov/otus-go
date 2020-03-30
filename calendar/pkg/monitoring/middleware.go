package monitoring

import (
	"net/http"
)

type Middleware interface {
	PrepareMetricExportHandler() http.Handler
}
