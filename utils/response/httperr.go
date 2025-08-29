package httperr

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// HTTPError mirrors your TS t_http_error shape.
type HTTPError struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Request    RequestInfo `json:"request"`
	Message    string      `json:"message"`
	Data       any         `json:"data"`            // keep as null in JSON (no omitempty)
	Trace      *Trace      `json:"trace,omitempty"` // omitted in PROD
}

type RequestInfo struct {
	IP     *string `json:"ip,omitempty"` // omitted in PROD
	Method string  `json:"method"`
	URL    string  `json:"url"`
}

type Trace struct {
	Error string `json:"error"`
}

// Logger is a tiny interface so you can plug in zap/logrus/slog.
type Logger interface {
	Error(msg string, args ...any)
}

func messageFromError(err error) string {
	if err == nil {
		return "Something went wrong"
	}
	if msg := err.Error(); msg != "" {
		return msg
	}
	return "Something went wrong"
}

// Build builds the error object. Pass prod=true to redact ip & trace.
func Build(c *gin.Context, err error, status int, prod bool) HTTPError {
	reqInfo := RequestInfo{
		Method: c.Request.Method,
		URL:    c.Request.URL.RequestURI(), // like Express originalUrl
	}

	if !prod {
		ip := c.ClientIP()
		reqInfo.IP = &ip
	}

	var tr *Trace
	if !prod && err != nil {
		tr = &Trace{Error: string(debug.Stack())}
	}

	return HTTPError{
		Success:    false,
		StatusCode: status,
		Request:    reqInfo,
		Message:    messageFromError(err),
		Data:       nil,
		Trace:      tr,
	}
}

// Write logs and writes the JSON error, aborting the context.
func Write(c *gin.Context, log Logger, err error, status int, prod bool) {
	obj := Build(c, err, status, prod)

	if log != nil {
		// Log the full meta object like your TS logger
		log.Error("CONTROLLER_ERROR", "meta", obj)
	}

	c.AbortWithStatusJSON(status, obj)
}

// Fail is a helper to be used inside handlers for early returns.
func Fail(c *gin.Context, log Logger, err error, status int, prod bool) {
	Write(c, log, err, status, prod)
}
