package actions

import "github.com/behavior3/behavior3go/core"

type Wait struct {
	core.Action
	EndTime int
}

func NewWait(milliseconds int) *Wait {
	node := &Wait{EndTime: milliseconds}
	node.Action = *core.NewAction(core.ActionOptions{
		Name:       "Wait",
		Title:      "Wait <milliseconds>ms",
		Properties: map[string]any{"milliseconds": 0},
	})
	return node
}

func (node *Wait) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Wait) Open(tick *core.Tick) {
	startTime := core.Now()
	tick.Blackboard.Set("startTime", startTime, tick.Tree.Id, node.Id)
}

func (node *Wait) Tick(tick *core.Tick) core.Status {
	currTime := core.Now()
	startTime, _ := core.ToInt64(tick.Blackboard.Get("startTime", tick.Tree.Id, node.Id))
	if currTime-startTime > int64(node.EndTime) {
		return core.SUCCESS
	}
	return core.RUNNING
}

func newWait(properties map[string]any) (core.Node, error) {
	milliseconds, _ := core.GetIntProperty(properties, "milliseconds", 0)
	return NewWait(milliseconds), nil
}

func init() {
	core.Register("Wait", newWait)
}
