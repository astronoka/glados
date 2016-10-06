package glados

// Logger is glados logger interface
type Logger interface {
	Debugf(string, ...interface{})
	Debugln(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}
