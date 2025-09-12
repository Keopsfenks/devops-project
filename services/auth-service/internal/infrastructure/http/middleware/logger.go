package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			err := next(c)

			end := time.Now()
			latency := end.Sub(start)

			method := req.Method
			uri := req.RequestURI
			status := res.Status
			size := res.Size
			userAgent := req.UserAgent()
			clientIP := c.RealIP()

			logMsg := fmt.Sprintf("[%s] %s %s %d %d %v %s %s",
				end.Format("2006-01-02 15:04:05"),
				method,
				uri,
				status,
				size,
				latency,
				clientIP,
				userAgent,
			)

			if status >= 500 {
				log.Printf("ERROR: %s", logMsg)
			} else if status >= 400 {
				log.Printf("WARN: %s", logMsg)
			} else {
				log.Printf("INFO: %s", logMsg)
			}

			return err
		}
	})
}

type LoggerConfig struct {
	Format     string
	TimeFormat string
	SkipPaths  []string
	Output     *log.Logger
}

func CustomLogger(config LoggerConfig) echo.MiddlewareFunc {
	if config.TimeFormat == "" {
		config.TimeFormat = "2006-01-02 15:04:05"
	}

	if config.Format == "" {
		config.Format = "[${time}] ${method} ${uri} ${status} ${latency} ${remote_ip}"
	}

	if config.Output == nil {
		config.Output = log.Default()
	}

	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, path := range config.SkipPaths {
				if c.Request().URL.Path == path {
					return next(c)
				}
			}

			start := time.Now()

			err := next(c)

			end := time.Now()
			latency := end.Sub(start)

			logLine := config.Format
			logLine = replaceLogPlaceholders(logLine, c, start, end, latency)

			config.Output.Println(logLine)

			return err
		}
	})
}

func replaceLogPlaceholders(format string, c echo.Context, start, end time.Time, latency time.Duration) string {
	req := c.Request()
	res := c.Response()

	replacements := map[string]string{
		"${time}":       end.Format("2006-01-02 15:04:05"),
		"${method}":     req.Method,
		"${uri}":        req.RequestURI,
		"${status}":     fmt.Sprintf("%d", res.Status),
		"${latency}":    latency.String(),
		"${remote_ip}":  c.RealIP(),
		"${user_agent}": req.UserAgent(),
		"${size}":       fmt.Sprintf("%d", res.Size),
	}

	result := format
	for _, value := range replacements {
		result = fmt.Sprintf("%s", result)
		result = fmt.Sprintf(result, value)
	}

	return result
}
