package core

type BaseNodeOptions struct {
	ID          string
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

type baseNodeConfig struct {
	generateID          bool
	cloneProperties     bool
	initializeParams    bool
	initializeEmptyMaps bool
}

func newBaseNode(options BaseNodeOptions, config baseNodeConfig) BaseNode {
	node := BaseNode{
		Id:          options.ID,
		Category:    options.Category,
		Name:        options.Name,
		Title:       options.Title,
		Description: options.Description,
		Properties:  options.Properties,
		Parameters:  nil,
	}

	if node.Id == "" && config.generateID {
		node.Id = createUUID()
	}
	if config.cloneProperties {
		node.Properties = copyMap(options.Properties)
	}
	if config.initializeParams {
		node.Parameters = map[string]any{}
	}

	if node.Title == "" {
		node.Title = node.Name
	}
	if node.Description == "" {
		node.Description = ""
	}
	if config.initializeEmptyMaps && node.Properties == nil {
		node.Properties = map[string]any{}
	}

	return node
}

func NewBaseNode(options BaseNodeOptions) BaseNode {
	return newBaseNode(options, baseNodeConfig{
		generateID:          true,
		cloneProperties:     true,
		initializeParams:    true,
		initializeEmptyMaps: true,
	})
}

func NewBaseNodeForLoad(options BaseNodeOptions) BaseNode {
	return newBaseNode(options, baseNodeConfig{
		generateID:          true,
		cloneProperties:     true,
		initializeParams:    true,
		initializeEmptyMaps: true,
	})
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
	return node.ExecuteNode(tick, node, node)
}

func (node *BaseNode) CloseNode(tick *Tick) {
	node.CloseNodeWithCallbacks(tick, node, node)
}

func (node *BaseNode) ExecuteNode(tick *Tick, current Node, callbacks NodeCallbacks) Status {
	baseNode := callbacks.GetBaseNode()

	tick.EnterNode(current)
	callbacks.Enter(tick)

	isOpen := false
	if tick.blackboard != nil {
		isOpen = tick.blackboard.isNodeOpen(tick.Tree.Id, baseNode.Id)
	} else {
		isOpen = boolValue(tick.Blackboard.Get("isOpen", tick.Tree.Id, baseNode.Id))
	}
	if !isOpen {
		tick.OpenNode(current)
		if tick.blackboard != nil {
			tick.blackboard.setNodeOpen(tick.Tree.Id, baseNode.Id, true)
		} else {
			tick.Blackboard.Set("isOpen", true, tick.Tree.Id, baseNode.Id)
		}
		callbacks.Open(tick)
	}

	tick.TickNode(current)
	status := callbacks.Tick(tick)

	if status != RUNNING {
		node.CloseNodeWithCallbacks(tick, current, callbacks)
	}

	tick.ExitNode(current)
	callbacks.Exit(tick)

	return status
}

func (node *BaseNode) CloseNodeWithCallbacks(tick *Tick, current Node, callbacks NodeCallbacks) {
	baseNode := callbacks.GetBaseNode()
	tick.CloseNode(current)
	if tick.blackboard != nil {
		tick.blackboard.setNodeOpen(tick.Tree.Id, baseNode.Id, false)
	} else {
		tick.Blackboard.Set("isOpen", false, tick.Tree.Id, baseNode.Id)
	}
	callbacks.Close(tick)
}
