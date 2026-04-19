package decorators

import "github.com/behavior3/behavior3go/core"

type Inverter struct {
	core.Decorator
}

func NewInverter(child core.Node) *Inverter {
	node := &Inverter{}
	node.Decorator = *core.NewDecorator(core.DecoratorOptions{
		Name:  "Inverter",
		Child: child,
	})
	return node
}

func (node *Inverter) Execute(tick *core.Tick) core.Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *Inverter) Tick(tick *core.Tick) core.Status {
	if node.Child == nil {
		return core.ERROR
	}

	status := node.Child.Execute(tick)
	if status == core.SUCCESS {
		return core.FAILURE
	}
	if status == core.FAILURE {
		return core.SUCCESS
	}
	return status
}

func newInverter(properties map[string]any) (core.Node, error) {
	return NewInverter(nil), nil
}

func loadInverter(spec core.NodeData) (core.Node, error) {
	node := &Inverter{}
	node.Decorator = *core.NewDecoratorForLoad(spec.Id, core.DecoratorOptions{Name: "Inverter", Properties: spec.Properties})
	return node, nil
}

func init() {
	core.Register("Inverter", newInverter)
	core.RegisterLoadConstructor("Inverter", loadInverter)
}
