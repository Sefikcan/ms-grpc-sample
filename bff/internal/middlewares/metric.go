package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sefikcan/ms-grpc-sample/bff/pkg/metric"
	"time"
)

func (mw *MiddlewareManager) MetricMiddleware(metrics metric.Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			var status int
			if err != nil {
				var HTTPError *echo.HTTPError
				errors.As(err, &HTTPError)
				status = HTTPError.Code
			} else {
				status = c.Response().Status
			}

			metrics.ObserveResponseTime(status, c.Request().Method, c.Path(), time.Since(start).Seconds())
			metrics.IncHits(status, c.Request().Method, c.Path())

			return err
		}
	}
}
