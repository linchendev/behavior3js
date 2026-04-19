package actions

import "github.com/behavior3/behavior3go/core"

type Runner struct {
	core.Action
}

func NewRunner() *Runner {
	node := &Runner{}
	node.Action = *core.NewAction(core.ActionOptions{Name: "Runner"})
	return node
}

func (node *Runner) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Runner) Tick(tick *core.Tick) core.Status {
	return core.RUNNING
}

func newRunner(properties map[string]any) (core.Node, error) {
	return NewRunner(), nil
}

func loadRunner(spec core.NodeData) (core.Node, error) {
	node := &Runner{}
	node.Action = *core.NewActionForLoad(spec.Id, core.ActionOptions{Name: "Runner", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Runner", newRunner)
	core.RegisterLoadConstructor("Runner", loadRunner)
}
