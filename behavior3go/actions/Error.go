package actions

import "github.com/behavior3/behavior3go/core"

type Error struct {
	core.Action
}

func NewError() *Error {
	node := &Error{}
	node.Action = *core.NewAction(core.ActionOptions{Name: "Error"})
	return node
}

func (node *Error) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Error) Tick(tick *core.Tick) core.Status {
	return core.ERROR
}

func newError(properties map[string]any) (core.Node, error) {
	return NewError(), nil
}

func loadError(spec core.NodeData) (core.Node, error) {
	node := &Error{}
	node.Action = *core.NewActionForLoad(spec.Id, core.ActionOptions{Name: "Error", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Error", newError)
	core.RegisterLoadConstructor("Error", loadError)
}
