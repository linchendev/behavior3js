package actions

import "github.com/behavior3/behavior3go/core"

type Failer struct {
	core.Action
}

func NewFailer() *Failer {
	node := &Failer{}
	node.Action = *core.NewAction(core.ActionOptions{Name: "Failer"})
	return node
}

func (node *Failer) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node)
}

func (node *Failer) Tick(tick *core.Tick) core.Status {
	return core.FAILURE
}

func newFailer(properties map[string]any) (core.Node, error) {
	return NewFailer(), nil
}

func init() {
	core.Register("Failer", newFailer)
}
