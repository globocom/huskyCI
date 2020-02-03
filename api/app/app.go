package app

// Application holds all information regarding the application
type Application struct {
}

// New returns a new Application struct
func New() *Application {
	return &Application{}
}

// Start starts the application
func (a *Application) Start() error {
	return nil
}

// Stop stops the applicaiton
func (a *Application) Stop() error {
	return nil
}
