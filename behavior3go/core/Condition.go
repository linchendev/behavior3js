package core

type ConditionOptions struct {
	Name       string
	Title      string
	Properties map[string]any
}

type Condition struct {
	BaseNode
}

func NewCondition(options ...ConditionOptions) *Condition {
	config := ConditionOptions{
		Name: "Condition",
	}
	if len(options) > 0 {
		config = options[0]
		if config.Name == "" {
			config.Name = "Condition"
		}
	}

	node := &Condition{}
	node.BaseNode = NewBaseNode(BaseNodeOptions{
		Category:   CONDITION,
		Name:       config.Name,
		Title:      config.Title,
		Properties: config.Properties,
	})
	return node
}

func (node *Condition) Execute(tick *Tick) Status {
	return node.BaseNode.ExecuteNode(tick, node)
}
