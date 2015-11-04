package plugins

import (
	"github.com/boxtown/verto"
	"io"
	"net/http"
	"strings"
)

type compressionWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w compressionWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w compressionWriter) Write(b []byte) (int, error) {
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w compressionWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// CompressionPlugin returns a VertoPluginFunc that handles
// gzip/deflate encoding.
func CompressionPlugin() verto.PluginFunc {
	return verto.PluginFunc(compressionFunc)
}

func compressionFunc(c *verto.Context, next http.HandlerFunc) {
	r := c.Request
	w := c.Response

	w.Header().Add("Vary", "Accept-Encoding")

	enc := strings.Split(r.Header.Get("Accept-Encoding"), ",")
	for _, v := range enc {
		v = strings.TrimSpace(v)
		if v == "gzip" {
			w.Header().Add("Content-Encoding", "gzip")

			ref := pool.get(w, ctGzip)
			defer ref.dispose()

			w = &compressionWriter{
				Writer:         ref.w,
				ResponseWriter: w,
			}
			next(w, r)
			return
		}
		if v == "deflate" {
			w.Header().Add("Content-Encoding", "deflate")

			ref := pool.get(w, ctFlate)
			defer ref.dispose()

			w = &compressionWriter{
				Writer:         ref.w,
				ResponseWriter: w,
			}
			next(w, r)
			return
		}
	}

	next(w, r)
}
