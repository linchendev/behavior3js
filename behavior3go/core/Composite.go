package core

type CompositeOptions struct {
	Children   []Node
	Name       string
	Title      string
	Properties map[string]any
}

type Composite struct {
	BaseNode
	Children []Node
}

func NewComposite(options ...CompositeOptions) *Composite {
	config := CompositeOptions{
		Name:     "Composite",
		Children: []Node{},
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Composite"
		}
		if config.Children == nil {
			config.Children = []Node{}
		}
	}

	node := &Composite{}
	node.BaseNode = NewBaseNode(BaseNodeOptions{
		Category:   COMPOSITE,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	node.Children = append([]Node{}, config.Children...)
	return node
}

func (node *Composite) Execute(tick *Tick) Status {
	return node.BaseNode.ExecuteNode(tick, node, node)
}

func (node *Composite) GetChildren() []Node {
	return node.Children
}

func (node *Composite) AppendChild(child Node) {
	node.Children = append(node.Children, child)
}
