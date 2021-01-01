package configuration

// Container wraps all subconfiguration
type Container struct {
	// ApplicationSettings apply to the whole app.
	ApplicationSettings AppSettings `yaml:"applicationSettings"`
	// SourceConfigs defines the behaviours of sinks/bridges.
	SourceConfigs []NodeConfig `yaml:"sourceConfigurations"`
	// RelayConfigs defines the flow of data after being received.
	RelayConfigs []NodeConfig `yaml:"relayConfigs"`
	// DistributorConfigs defines the behaviour of distribution mechanisms.
	DistributorConfigs []NodeConfig `yaml:"distributorConfigurations"`
}

// -------------------------------------------------------------------------------------

// AppSettings is application level configuration.
type AppSettings struct{}

// -------------------------------------------------------------------------------------

// NodeConfig defines the configuration for any given pluggable structure.
type NodeConfig struct {
	Name       string                 `yaml:"name"`
	Config     map[string]interface{} `yaml:"config"`
	SubConfigs []NodeConfig           `yaml:"subConfigs"`
}
