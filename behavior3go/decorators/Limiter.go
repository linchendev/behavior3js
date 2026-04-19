package decorators

import (
	"fmt"

	"github.com/behavior3/behavior3go/core"
)

type Limiter struct {
	core.Decorator
	MaxLoop int
}

func NewLimiter(maxLoop int, child core.Node) *Limiter {
	if maxLoop == 0 {
		panic("maxLoop parameter in Limiter decorator is an obligatory parameter")
	}

	node := &Limiter{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:       "Limiter",
		Title:      "Limit <maxLoop> Activations",
		Child:      child,
		Properties: map[string]any{"maxLoop": 1},
	})
	return node
}

func (node *Limiter) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Limiter) Open(tick *core.Tick) {
	tick.Blackboard.Set("i", 0, tick.Tree.Id, node.Id)
}

func (node *Limiter) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	i, _ := core.GetIntProperty(map[string]any{
		"i": tick.Blackboard.Get("i", tick.Tree.Id, node.Id),
	}, "i", 0)

	if i < node.MaxLoop {
		status := node.Child.Execute(tick)
		if status == core.SUCCESS || status == core.FAILURE {
			tick.Blackboard.Set("i", i+1, tick.Tree.Id, node.Id)
		}
		return status
	}

	return core.FAILURE
}

func newLimiter(properties map[string]any) (core.Node, error) {
	maxLoop, ok := core.GetIntProperty(properties, "maxLoop", 0)
	if !ok || maxLoop == 0 {
		return nil, fmt.Errorf("maxLoop parameter in Limiter decorator is an obligatory parameter")
	}
	return NewLimiter(maxLoop, nil), nil
}

func loadLimiter(spec core.NodeData) (core.Node, error) {
	maxLoop, ok := core.GetIntProperty(spec.Properties, "maxLoop", 0)
	if !ok || maxLoop == 0 {
		return nil, fmt.Errorf("maxLoop parameter in Limiter decorator is an obligatory parameter")
	}
	node := &Limiter{MaxLoop: maxLoop}
	node.Decorator = *core.NewDecoratorForLoad(spec.Id, core.DecoratorOptions{
		Name:       "Limiter",
		Title:      "Limit <maxLoop> Activations",
		Properties: spec.Properties,
	})
	return node, nil
}

func init() {
	core.Register("Limiter", newLimiter)
	core.RegisterLoadConstructor("Limiter", loadLimiter)
}
