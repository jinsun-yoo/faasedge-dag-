package dag

import "log"

type Dag struct {
	Name      string     `yaml:"name"`
	Functions []Function `yaml:"functions"`
}

type FlowCallback func(d *Dag, id string, parentResults []FlowResult)

type FlowResult struct {
	ID string
	Result map[string]interface{}
}

func (dag *Dag) DescendantsFlow(startID string, callback FlowCallback) {
	var function Function;
	var found bool = false;
	for _, dag := range dag.Functions {
		if dag.Name == startID {
			function = dag;
			found = true;
		}
	}

	if (!found) {
		log.Println("Descendants flow error: couldn't find provided start ID")
	}

	log.Println(function)

}

func (dag *Dag) descendantsFlowHelper(startID string, callback FlowCallback, parentResults []FlowResult) {

}

type Function struct {
	Name                  string       `yaml:"name"`
	Image                 string       `yaml:"image"`
	DownstreamConnections []Connection `yaml:"downstream_connections"`
}

type Connection struct {
	DownstreamFunctions string         `yaml:"downstream_functions"`
	ConnectionType      ConnectionType `yaml:"connection_type"`
}

type ConnectionType int

const (
	Unknown   ConnectionType = iota // 0
	Immediate                       // 1
	FanIn                           // 2
	FanOut                          // 3
)
