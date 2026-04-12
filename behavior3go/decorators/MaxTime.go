package decorators

import (
	"fmt"

	"github.com/behavior3/behavior3go/core"
)

type MaxTime struct {
	core.Decorator
	MaxTime int
}

func NewMaxTime(maxTime int, child core.Node) *MaxTime {
	if maxTime == 0 {
		panic("maxTime parameter in MaxTime decorator is an obligatory parameter")
	}

	node := &MaxTime{MaxTime: maxTime}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:       "MaxTime",
		Title:      "Max <maxTime>ms",
		Child:      child,
		Properties: map[string]any{"maxTime": 0},
	})
	return node
}

func (node *MaxTime) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *MaxTime) Open(tick *core.Tick) {
	tick.Blackboard.Set("startTime", core.Now(), tick.Tree.Id, node.Id)
}

func (node *MaxTime) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	currTime := core.Now()
	startTime, _ := core.ToInt64(tick.Blackboard.Get("startTime", tick.Tree.Id, node.Id))
	status := node.Child.Execute(tick)
	if currTime-startTime > int64(node.MaxTime) {
		return core.FAILURE
	}
	return status
}

func newMaxTime(properties map[string]any) (core.Node, error) {
	maxTime, ok := core.GetIntProperty(properties, "maxTime", 0)
	if !ok || maxTime == 0 {
		return nil, fmt.Errorf("maxTime parameter in MaxTime decorator is an obligatory parameter")
	}
	return NewMaxTime(maxTime, nil), nil
}

func init() {
	core.Register("MaxTime", newMaxTime)
}
