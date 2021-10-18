package handler

import "net/http"

type responseWriter struct {
	w http.ResponseWriter

	statusCode int
}

// Header 实现http.ResponseWriter接口
func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

// Write 实现http.ResponseWriter接口
func (r *responseWriter) Write(bs []byte) (int, error) {
	return r.w.Write(bs)
}

// WriteHeader 实现http.ResponseWriter接口
func (r *responseWriter) WriteHeader(statusCode int) {
	// 此处拦截http的返回码
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

func wrapperResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		w: w,
	}
}

