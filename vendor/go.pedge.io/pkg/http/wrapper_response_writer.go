package pkghttp

import "net/http"

type wrapperResponseWriter interface {
	http.ResponseWriter
	StatusCode() int
	WriteError() error
}

func newWrapperResponseWriter(responseWriter http.ResponseWriter) wrapperResponseWriter {
	flusher, flush := responseWriter.(http.Flusher)
	hijacker, hijack := responseWriter.(http.Hijacker)
	closeNotifier, closeNotify := responseWriter.(http.CloseNotifier)
	if flush {
		if hijack {
			if closeNotify {
				return newWrapperResponseWriteFlushHijackCloseNotify(responseWriter, flusher, hijacker, closeNotifier)
			}
			return newWrapperResponseWriteFlushHijack(responseWriter, flusher, hijacker)
		}
		if closeNotify {
			return newWrapperResponseWriteFlushCloseNotify(responseWriter, flusher, closeNotifier)
		}
		return newWrapperResponseWriteFlush(responseWriter, flusher)
	}
	if hijack {
		if closeNotify {
			return newWrapperResponseWriteHijackCloseNotify(responseWriter, hijacker, closeNotifier)
		}
		return newWrapperResponseWriteHijack(responseWriter, hijacker)
	}
	if closeNotify {
		return newWrapperResponseWriteCloseNotify(responseWriter, closeNotifier)
	}
	return newWrapperResponseWrite(responseWriter)
}

func newWrapperResponseWrite(responseWriter http.ResponseWriter) wrapperResponseWriter {
	return &wrapperResponseWrite{responseWriter, 0, nil}
}

type wrapperResponseWrite struct {
	http.ResponseWriter
	statusCode int
	writeError error
}

func (w *wrapperResponseWrite) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.writeError = err
	return n, err
}

func (w *wrapperResponseWrite) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *wrapperResponseWrite) StatusCode() int {
	return w.statusCode
}

func (w *wrapperResponseWrite) WriteError() error {
	return w.writeError
}

type wrapperResponseWriteFlush struct {
	wrapperResponseWriter
	http.Flusher
}

func newWrapperResponseWriteFlush(
	responseWriter http.ResponseWriter,
	flusher http.Flusher,
) wrapperResponseWriter {
	return &wrapperResponseWriteFlush{
		newWrapperResponseWrite(responseWriter),
		flusher,
	}
}

type wrapperResponseWriteHijack struct {
	wrapperResponseWriter
	http.Hijacker
}

func newWrapperResponseWriteHijack(
	responseWriter http.ResponseWriter,
	hijacker http.Hijacker,
) wrapperResponseWriter {
	return &wrapperResponseWriteHijack{
		newWrapperResponseWrite(responseWriter),
		hijacker,
	}
}

type wrapperResponseWriteCloseNotify struct {
	wrapperResponseWriter
	http.CloseNotifier
}

func newWrapperResponseWriteCloseNotify(
	responseWriter http.ResponseWriter,
	closeNotifier http.CloseNotifier,
) wrapperResponseWriter {
	return &wrapperResponseWriteCloseNotify{
		newWrapperResponseWrite(responseWriter),
		closeNotifier,
	}
}

type wrapperResponseWriteFlushHijack struct {
	wrapperResponseWriter
	http.Flusher
	http.Hijacker
}

func newWrapperResponseWriteFlushHijack(
	responseWriter http.ResponseWriter,
	flusher http.Flusher,
	hijacker http.Hijacker,
) wrapperResponseWriter {
	return &wrapperResponseWriteFlushHijack{
		newWrapperResponseWrite(responseWriter),
		flusher,
		hijacker,
	}
}

type wrapperResponseWriteFlushCloseNotify struct {
	wrapperResponseWriter
	http.Flusher
	http.CloseNotifier
}

func newWrapperResponseWriteFlushCloseNotify(
	responseWriter http.ResponseWriter,
	flusher http.Flusher,
	closeNotifier http.CloseNotifier,
) wrapperResponseWriter {
	return &wrapperResponseWriteFlushCloseNotify{
		newWrapperResponseWrite(responseWriter),
		flusher,
		closeNotifier,
	}
}

type wrapperResponseWriteHijackCloseNotify struct {
	wrapperResponseWriter
	http.Hijacker
	http.CloseNotifier
}

func newWrapperResponseWriteHijackCloseNotify(
	responseWriter http.ResponseWriter,
	hijacker http.Hijacker,
	closeNotifier http.CloseNotifier,
) wrapperResponseWriter {
	return &wrapperResponseWriteHijackCloseNotify{
		newWrapperResponseWrite(responseWriter),
		hijacker,
		closeNotifier,
	}
}

type wrapperResponseWriteFlushHijackCloseNotify struct {
	wrapperResponseWriter
	http.Flusher
	http.Hijacker
	http.CloseNotifier
}

func newWrapperResponseWriteFlushHijackCloseNotify(
	responseWriter http.ResponseWriter,
	flusher http.Flusher,
	hijacker http.Hijacker,
	closeNotifier http.CloseNotifier,
) wrapperResponseWriter {
	return &wrapperResponseWriteFlushHijackCloseNotify{
		newWrapperResponseWrite(responseWriter),
		flusher,
		hijacker,
		closeNotifier,
	}
}
