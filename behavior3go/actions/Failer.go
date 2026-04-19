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
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Failer) Tick(tick *core.Tick) core.Status {
	return core.FAILURE
}

func newFailer(properties map[string]any) (core.Node, error) {
	return NewFailer(), nil
}

func loadFailer(spec core.NodeData) (core.Node, error) {
	node := &Failer{}
	node.Action = *core.NewActionForLoad(spec.Id, core.ActionOptions{Name: "Failer", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Failer", newFailer)
	core.RegisterLoadConstructor("Failer", loadFailer)
}
