package composites

import "github.com/behavior3/behavior3go/core"

type Sequence struct {
	core.Composite
}

func NewSequence(children ...core.Node) *Sequence {
	node := &Sequence{}
	node.Composite = *core.NewComposite(core.CompositeOptions{
		Name:     "Sequence",
		Children: children,
	})
	return node
}

func (node *Sequence) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Sequence) Tick(tick *core.Tick) core.Status {
	for _, child := range node.Children {
		status := child.Execute(tick)
		if status != core.SUCCESS {
			return status
		}
	}
	return core.SUCCESS
}

func newSequence(properties map[string]any) (core.Node, error) {
	return NewSequence(), nil
}

func loadSequence(spec core.NodeData) (core.Node, error) {
	node := &Sequence{}
	node.Composite = *core.NewCompositeForLoad(spec.Id, core.CompositeOptions{Name: "Sequence", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Sequence", newSequence)
	core.RegisterLoadConstructor("Sequence", loadSequence)
}
