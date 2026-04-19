package core

type ActionOptions struct {
	Name       string
	Title      string
	Properties map[string]any
}

type Action struct {
	BaseNode
}

func NewAction(options ...ActionOptions) *Action {
	config := ActionOptions{
		Name: "Action",
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Action"
		}
	}

	node := &Action{}
	node.BaseNode = NewBaseNode(BaseNodeOptions{
		Category:   ACTION,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	return node
}

func NewActionForLoad(id string, options ...ActionOptions) *Action {
	config := ActionOptions{
		Name: "Action",
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Action"
		}
	}

	node := &Action{}
	node.BaseNode = NewBaseNodeForLoad(BaseNodeOptions{
		ID:         id,
		Category:   ACTION,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	return node
}

func (node *Action) Execute(tick *Tick) Status {
	return node.BaseNode.ExecuteNode(tick, node, node)
}
