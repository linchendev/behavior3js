package actions

import "github.com/behavior3/behavior3go/core"

type Succeeder struct {
	core.Action
}

func NewSucceeder() *Succeeder {
	node := &Succeeder{}
	node.Action = *core.NewAction(core.ActionOptions{Name: "Succeeder"})
	return node
}

func (node *Succeeder) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node)
}

func (node *Succeeder) Tick(tick *core.Tick) core.Status {
	return core.SUCCESS
}

func newSucceeder(properties map[string]any) (core.Node, error) {
	return NewSucceeder(), nil
}

func init() {
	core.Register("Succeeder", newSucceeder)
}
