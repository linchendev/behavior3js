package core

type DecoratorOptions struct {
	Child      Node
	Name       string
	Title      string
	Properties map[string]any
}

type Decorator struct {
	BaseNode
	Child Node
}

func NewDecorator(options ...DecoratorOptions) *Decorator {
	config := DecoratorOptions{
		Name: "Decorator",
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Decorator"
		}
	}

	node := &Decorator{}
	node.BaseNode = NewBaseNode(BaseNodeOptions{
		Category:   DECORATOR,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	node.Child = config.Child
	return node
}

func NewDecoratorForLoad(id string, options ...DecoratorOptions) *Decorator {
	config := DecoratorOptions{
		Name: "Decorator",
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Decorator"
		}
	}

	node := &Decorator{}
	node.BaseNode = NewBaseNodeForLoad(BaseNodeOptions{
		ID:         id,
		Category:   DECORATOR,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	node.Child = config.Child
	return node
}

func (node *Decorator) Execute(tick *Tick) Status {
	return node.BaseNode.ExecuteNode(tick, node, node)
}

func (node *Decorator) GetChild() Node {
	return node.Child
}

func (node *Decorator) SetChild(child Node) {
	node.Child = child
}
