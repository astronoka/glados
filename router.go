package glados

// RequestHandler is golados request interface
type RequestHandler func(RequestContext)

// Router is request handler mapper
type Router interface {
	GET(string, RequestHandler)
	POST(string, RequestHandler)
	RunWithPort(string)
}
