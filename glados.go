package glados

// New is create glados instance
func New(c Context) *Glados {
	return &Glados{
		context: c,
	}
}

// Glados is Github Lifeform and Disk Operating System
type Glados struct {
	context Context
}

// Boot is boot up glados instance
func (g Glados) Boot() {
	logger := g.context.Logger()
	logger.Infoln("Glados: Boot..")
	defer logger.Infoln("Glados: Shutdown..")
	port := g.context.ListenPort()
	g.context.Router().RunWithPort(port)
}

// Install is register glados logic
func (g *Glados) Install(p Program) {
	p.Initialize(g.context)
}
