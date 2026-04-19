package composites

import "github.com/behavior3/behavior3go/core"

type Priority struct {
	core.Composite
}

func NewPriority(children ...core.Node) *Priority {
	node := &Priority{}
	node.Composite = *core.NewComposite(core.CompositeOptions{
		Name:     "Priority",
		Children: children,
	})
	return node
}

func (node *Priority) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Priority) Tick(tick *core.Tick) core.Status {
	for _, child := range node.Children {
		status := child.Execute(tick)
		if status != core.FAILURE {
			return status
		}
	}
	return core.FAILURE
}

func newPriority(properties map[string]any) (core.Node, error) {
	return NewPriority(), nil
}

func loadPriority(spec core.NodeData) (core.Node, error) {
	node := &Priority{}
	node.Composite = *core.NewCompositeForLoad(spec.Id, core.CompositeOptions{Name: "Priority", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Priority", newPriority)
	core.RegisterLoadConstructor("Priority", loadPriority)
}
