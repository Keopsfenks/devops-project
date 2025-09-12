package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CompressionConfig struct {
	Level          int
	MinLength      int
	ExcludedTypes  []string
	ExcludedPaths  []string
	EnableForHTTPS bool
}

var DefaultCompressionConfig = CompressionConfig{
	Level:          gzip.DefaultCompression,
	MinLength:      1024,
	EnableForHTTPS: true,
	ExcludedTypes: []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"video/mp4",
		"video/webm",
		"application/pdf",
		"application/zip",
		"application/gzip",
	},
	ExcludedPaths: []string{},
}

func Compression(next echo.HandlerFunc) echo.HandlerFunc {
	return CompressionWithConfig(DefaultCompressionConfig)(next)
}

func CompressionWithConfig(config CompressionConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
				return next(c)
			}

			if !config.EnableForHTTPS && req.TLS != nil {
				return next(c)
			}

			for _, path := range config.ExcludedPaths {
				if strings.HasPrefix(req.URL.Path, path) {
					return next(c)
				}
			}

			res.Header().Set("Content-Encoding", "gzip")
			res.Header().Set("Vary", "Accept-Encoding")

			gw, err := gzip.NewWriterLevel(res.Writer, config.Level)
			if err != nil {
				return err
			}
			defer func(gw *gzip.Writer) {
				err := gw.Close()
				if err != nil {
					c.Logger().Errorf("gzip writer close error: %v", err)
				}
			}(gw)

			grw := &gzipResponseWriter{
				Writer:         gw,
				ResponseWriter: res.Writer,
				config:         config,
			}
			res.Writer = grw

			return next(c)
		}
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	config        CompressionConfig
	headerWritten bool
	contentLength int
	contentType   string
}

func (g *gzipResponseWriter) WriteHeader(code int) {
	if g.headerWritten {
		return
	}

	g.contentType = g.Header().Get("Content-Type")

	if g.shouldCompress() {
		g.Header().Del("Content-Length")
		g.ResponseWriter.WriteHeader(code)
		g.headerWritten = true
	} else {
		g.Header().Del("Content-Encoding")
		g.Header().Del("Vary")
		g.ResponseWriter.WriteHeader(code)
		g.headerWritten = true
	}
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.headerWritten {
		g.contentLength += len(b)

		if g.contentLength < g.config.MinLength {
			return g.ResponseWriter.Write(b)
		}

		g.WriteHeader(http.StatusOK)
	}

	if g.shouldCompress() {
		return g.Writer.Write(b)
	}
	return g.ResponseWriter.Write(b)
}

func (g *gzipResponseWriter) shouldCompress() bool {
	if g.contentLength < g.config.MinLength && g.contentLength > 0 {
		return false
	}

	for _, excludedType := range g.config.ExcludedTypes {
		if strings.Contains(g.contentType, excludedType) {
			return false
		}
	}

	return true
}

func (g *gzipResponseWriter) Flush() {
	if gw, ok := g.Writer.(*gzip.Writer); ok {
		err := gw.Flush()
		if err != nil {
			return
		}
	}
	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
