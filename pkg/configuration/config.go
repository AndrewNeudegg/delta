package configuration

// ConfigLoader is a generic configuration loader.
type ConfigLoader interface {
	Load() (Container, error)
	// Write(Container) error
}
