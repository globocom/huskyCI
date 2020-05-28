package runner

// RSession is the runner interface
type RSession interface {
	Close() error
	Ping() error
	Run(containerImage string, command string) (string, error)
}
