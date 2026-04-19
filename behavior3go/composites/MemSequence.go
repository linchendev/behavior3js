package composites

import "github.com/behavior3/behavior3go/core"

type MemSequence struct {
	core.Composite
}

func NewMemSequence(children ...core.Node) *MemSequence {
	node := &MemSequence{}
	node.Composite = *core.NewComposite(core.CompositeOptions{
		Name:     "MemSequence",
		Children: children,
	})
	return node
}

func (node *MemSequence) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *MemSequence) Open(tick *core.Tick) {
	tick.Blackboard.Set("runningChild", 0, tick.Tree.Id, node.Id)
}

func (node *MemSequence) Tick(tick *core.Tick) core.Status {
	childIndex, _ := core.GetIntProperty(map[string]any{
		"runningChild": tick.Blackboard.Get("runningChild", tick.Tree.Id, node.Id),
	}, "runningChild", 0)

	for i := childIndex; i < len(node.Children); i++ {
		status := node.Children[i].Execute(tick)
		if status != core.SUCCESS {
			if status == core.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.Tree.Id, node.Id)
			}
			return status
		}
	}

	return core.SUCCESS
}

func newMemSequence(properties map[string]any) (core.Node, error) {
	return NewMemSequence(), nil
}

func loadMemSequence(spec core.NodeData) (core.Node, error) {
	node := &MemSequence{}
	node.Composite = *core.NewCompositeForLoad(spec.Id, core.CompositeOptions{Name: "MemSequence", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("MemSequence", newMemSequence)
	core.RegisterLoadConstructor("MemSequence", loadMemSequence)
}
