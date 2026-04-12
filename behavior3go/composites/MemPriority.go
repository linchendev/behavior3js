package composites

import "github.com/behavior3/behavior3go/core"

type MemPriority struct {
	core.Composite
}

func NewMemPriority(children ...core.Node) *MemPriority {
	node := &MemPriority{}
	node.Composite = *core.NewComposite(core.CompositeOptions{
		Name:     "MemPriority",
		Children: children,
	})
	return node
}

func (node *MemPriority) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node)
}

func (node *MemPriority) Open(tick *core.Tick) {
	tick.Blackboard.Set("runningChild", 0, tick.Tree.Id, node.Id)
}

func (node *MemPriority) Tick(tick *core.Tick) core.Status {
	childIndex, _ := core.GetIntProperty(map[string]any{
		"runningChild": tick.Blackboard.Get("runningChild", tick.Tree.Id, node.Id),
	}, "runningChild", 0)

	for i := childIndex; i < len(node.Children); i++ {
		status := node.Children[i].Execute(tick)
		if status != core.FAILURE {
			if status == core.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.Tree.Id, node.Id)
			}
			return status
		}
	}

	return core.FAILURE
}

func newMemPriority(properties map[string]any) (core.Node, error) {
	return NewMemPriority(), nil
}

func init() {
	core.Register("MemPriority", newMemPriority)
}
