package decorators

import "github.com/behavior3/behavior3go/core"

type RepeatUntilSuccess struct {
	core.Decorator
	MaxLoop int
}

func NewRepeatUntilSuccess(maxLoop int, child core.Node) *RepeatUntilSuccess {
	node := &RepeatUntilSuccess{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:       "RepeatUntilSuccess",
		Title:      "Repeat Until Success",
		Child:      child,
		Properties: map[string]any{"maxLoop": -1},
	})
	return node
}

func (node *RepeatUntilSuccess) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *RepeatUntilSuccess) Open(tick *core.Tick) {
	tick.Blackboard.Set("i", 0, tick.Tree.Id, node.Id)
}

func (node *RepeatUntilSuccess) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	i, _ := core.GetIntProperty(map[string]any{
		"i": tick.Blackboard.Get("i", tick.Tree.Id, node.Id),
	}, "i", 0)
	status := core.ERROR

	for node.MaxLoop < 0 || i < node.MaxLoop {
		status = node.Child.Execute(tick)
		if status == core.FAILURE {
			i++
			continue
		}
		break
	}

	tick.Blackboard.Set("i", i, tick.Tree.Id, node.Id)
	return status
}

func newRepeatUntilSuccess(properties map[string]any) (core.Node, error) {
	maxLoop, _ := core.GetIntProperty(properties, "maxLoop", -1)
	return NewRepeatUntilSuccess(maxLoop, nil), nil
}

func init() {
	core.Register("RepeatUntilSuccess", newRepeatUntilSuccess)
}
