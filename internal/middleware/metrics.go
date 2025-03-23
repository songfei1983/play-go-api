package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/songfei1983/play-go-api/internal/metrics"
)

// MetricsMiddleware returns a middleware that collects HTTP metrics
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			// Record request duration
			duration := time.Since(start).Seconds()
			metrics.HTTPRequestDuration.WithLabelValues(
				c.Request().Method,
				c.Path(),
			).Observe(duration)

			// Record request count
			status := c.Response().Status
			metrics.HTTPRequestsTotal.WithLabelValues(
				c.Request().Method,
				c.Path(),
				strconv.Itoa(status),
			).Inc()

			return err
		}
	}
}
