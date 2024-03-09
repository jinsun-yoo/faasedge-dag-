package dag

type Dag struct {
	Name      string      `yaml:"name"`
	Functions []Function  `yaml:"functions"`
}

type Function struct {
	Name                 string       `yaml:"name"`
	Image                string       `yaml:"image"`
	DownstreamConnections []Connection `yaml:"downstream_connections"`
}

type Connection struct {
	DownstreamFunctions string         `yaml:"downstream_functions"`
	ConnectionType      ConnectionType `yaml:"connection_type"`
}

type ConnectionType int

const (
	Unknown ConnectionType = iota // 0
	Immediate                    // 1
	FanIn                        // 2
	FanOut                       // 3
)
