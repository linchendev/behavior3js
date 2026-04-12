package decorators

import "github.com/behavior3/behavior3go/core"

type Repeater struct {
	core.Decorator
	MaxLoop int
}

func NewRepeater(maxLoop int, child core.Node) *Repeater {
	node := &Repeater{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:       "Repeater",
		Title:      "Repeat <maxLoop>x",
		Child:      child,
		Properties: map[string]any{"maxLoop": -1},
	})
	return node
}

func (node *Repeater) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Repeater) Open(tick *core.Tick) {
	tick.Blackboard.Set("i", 0, tick.Tree.Id, node.Id)
}

func (node *Repeater) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	i, _ := core.GetIntProperty(map[string]any{
		"i": tick.Blackboard.Get("i", tick.Tree.Id, node.Id),
	}, "i", 0)
	status := core.SUCCESS

	for node.MaxLoop < 0 || i < node.MaxLoop {
		status = node.Child.Execute(tick)
		if status == core.SUCCESS || status == core.FAILURE {
			i++
			continue
		}
		break
	}

	tick.Blackboard.Set("i", i, tick.Tree.Id, node.Id)
	return status
}

func newRepeater(properties map[string]any) (core.Node, error) {
	maxLoop, _ := core.GetIntProperty(properties, "maxLoop", -1)
	return NewRepeater(maxLoop, nil), nil
}

func init() {
	core.Register("Repeater", newRepeater)
}
