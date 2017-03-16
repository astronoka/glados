package ginbind

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/astronoka/glados"
	"github.com/gin-gonic/gin"
)

// NewRouter is create ginRouter instance
func NewRouter(c glados.Context) glados.Router {
	engine := gin.New()
	engine.Use(accessLogger(c.Logger()), gin.Recovery())
	return ginRouter{
		context: c,
		engine:  engine,
	}
}

func accessLogger(logger glados.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		start := time.Now()
		log := gin.H{}
		log["path"] = ginContext.Request.URL.Path
		ginContext.Next()
		end := time.Now()
		log["end"] = end
		log["latency"] = end.Sub(start)
		log["client_ip"] = ginContext.ClientIP()
		log["method"] = ginContext.Request.Method
		log["status_code"] = ginContext.Writer.Status()
		log["comment"] = ginContext.Errors.ByType(gin.ErrorTypePrivate).String()
		bytes, err := json.Marshal(log)
		if err != nil {
			logger.Errorln("ginbind: output access log failed. " + err.Error())
		} else {
			logger.Infoln(string(bytes))
		}
	}
}

type ginRouter struct {
	context glados.Context
	engine  *gin.Engine
}

func (router ginRouter) GET(path string, handler glados.RequestHandler) {
	router.engine.GET(path, func(ginContext *gin.Context) {
		handler(ginContextWrapper{ginContext})
	})
}

func (router ginRouter) POST(path string, handler glados.RequestHandler) {
	router.engine.POST(path, func(ginContext *gin.Context) {
		handler(ginContextWrapper{ginContext})
	})
}

func (router ginRouter) RunWithPort(port string) {
	router.engine.Run(":" + port)
}

type ginContextWrapper struct {
	*gin.Context
}

func (w ginContextWrapper) Header(key string) string {
	return w.Request.Header.Get(key)
}

func (w ginContextWrapper) Request() *http.Request {
	return w.Request
}
