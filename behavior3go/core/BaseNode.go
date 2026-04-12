package core

type BaseNodeOptions struct {
	Category    string
	Name        string
	Title       string
	Description string
	Properties  map[string]any
}

type BaseNode struct {
	Id          string
	Category    string
	Name        string
	Title       string
	Description string
	Properties  map[string]any
	Parameters  map[string]any
}

func NewBaseNode(options BaseNodeOptions) BaseNode {
	node := BaseNode{
		Id:          createUUID(),
		Category:    options.Category,
		Name:        options.Name,
		Title:       options.Title,
		Description: options.Description,
		Properties:  copyMap(options.Properties),
		Parameters:  map[string]any{},
	}

	if node.Title == "" {
		node.Title = node.Name
	}
	if node.Description == "" {
		node.Description = ""
	}
	if node.Properties == nil {
		node.Properties = map[string]any{}
	}

	return node
}

func (node *BaseNode) GetBaseNode() *BaseNode {
	return node
}

func (node *BaseNode) Enter(tick *Tick) {}

func (node *BaseNode) Open(tick *Tick) {}

func (node *BaseNode) Tick(tick *Tick) Status {
	return 0
}

func (node *BaseNode) Close(tick *Tick) {}

func (node *BaseNode) Exit(tick *Tick) {}

func (node *BaseNode) Execute(tick *Tick) Status {
	return node.ExecuteNode(tick, node)
}

func (node *BaseNode) CloseNode(tick *Tick) {
	node.CloseNodeWithCallbacks(tick, node)
}

func (node *BaseNode) ExecuteNode(tick *Tick, callbacks NodeCallbacks) Status {
	baseNode := callbacks.GetBaseNode()

	tick.EnterNode(callbacks.(Node))
	callbacks.Enter(tick)

	if !boolValue(tick.Blackboard.Get("isOpen", tick.Tree.Id, baseNode.Id)) {
		tick.OpenNode(callbacks.(Node))
		tick.Blackboard.Set("isOpen", true, tick.Tree.Id, baseNode.Id)
		callbacks.Open(tick)
	}

	tick.TickNode(callbacks.(Node))
	status := callbacks.Tick(tick)

	if status != RUNNING {
		node.CloseNodeWithCallbacks(tick, callbacks)
	}

	tick.ExitNode(callbacks.(Node))
	callbacks.Exit(tick)

	return status
}

func (node *BaseNode) CloseNodeWithCallbacks(tick *Tick, callbacks NodeCallbacks) {
	baseNode := callbacks.GetBaseNode()
	tick.CloseNode(callbacks.(Node))
	tick.Blackboard.Set("isOpen", false, tick.Tree.Id, baseNode.Id)
	callbacks.Close(tick)
}
