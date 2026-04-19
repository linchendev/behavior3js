package decorators

import "github.com/behavior3/behavior3go/core"

type RepeatUntilFailure struct {
	core.Decorator
	MaxLoop int
}

func NewRepeatUntilFailure(maxLoop int, child core.Node) *RepeatUntilFailure {
	node := &RepeatUntilFailure{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:       "RepeatUntilFailure",
		Title:      "Repeat Until Failure",
		Child:      child,
		Properties: map[string]any{"maxLoop": -1},
	})
	return node
}

func (node *RepeatUntilFailure) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *RepeatUntilFailure) Open(tick *core.Tick) {
	tick.Blackboard.Set("i", 0, tick.Tree.Id, node.Id)
}

func (node *RepeatUntilFailure) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	i, _ := core.GetIntProperty(map[string]any{
		"i": tick.Blackboard.Get("i", tick.Tree.Id, node.Id),
	}, "i", 0)
	status := core.ERROR

	for node.MaxLoop < 0 || i < node.MaxLoop {
		status = node.Child.Execute(tick)
		if status == core.SUCCESS {
			i++
			continue
		}
		break
	}

	tick.Blackboard.Set("i", i, tick.Tree.Id, node.Id)
	return status
}

func newRepeatUntilFailure(properties map[string]any) (core.Node, error) {
	maxLoop, _ := core.GetIntProperty(properties, "maxLoop", -1)
	return NewRepeatUntilFailure(maxLoop, nil), nil
}

func loadRepeatUntilFailure(spec core.NodeData) (core.Node, error) {
	maxLoop, _ := core.GetIntProperty(spec.Properties, "maxLoop", -1)
	node := &RepeatUntilFailure{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecoratorForLoad(spec.Id, core.DecoratorOptions{
		Name:       "RepeatUntilFailure",
		Title:      "Repeat Until Failure",
		Properties: spec.Properties,
	})
	return node, nil
}

func init() {
	core.Register("RepeatUntilFailure", newRepeatUntilFailure)
	core.RegisterLoadConstructor("RepeatUntilFailure", loadRepeatUntilFailure)
}
