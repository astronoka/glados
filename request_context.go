package glados

import "net/http"

// H is hash map
type H map[string]interface{}

// RequestContext is request system context
type RequestContext interface {
	Header(string) string
	Param(string) string
	Request() *http.Request
	JSON(int, interface{})
}
