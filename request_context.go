package glados

import "io"

// H is hash map
type H map[string]interface{}

// RequestContext is request system context
type RequestContext interface {
	Header(string) string
	Param(string) string
	RequestBody() io.ReadCloser
	JSON(int, interface{})
}
