package core

type Blackboard struct {
	baseMemory map[string]any
	treeMemory map[string]*treeMemory
}

func NewBlackboard() *Blackboard {
	return &Blackboard{
		baseMemory: map[string]any{},
		treeMemory: map[string]*treeMemory{},
	}
}

func (blackboard *Blackboard) getTreeMemory(treeScope string) *treeMemory {
	memory, ok := blackboard.treeMemory[treeScope]
	if !ok {
		memory = &treeMemory{
			memory:         map[string]any{},
			nodeMemory:     map[string]map[string]any{},
			openNodes:      []Node{},
			traversalDepth: 0,
			traversalCycle: 0,
		}
		blackboard.treeMemory[treeScope] = memory
	}

	return memory
}

func (blackboard *Blackboard) getNodeMemory(treeMemory *treeMemory, nodeScope string) map[string]any {
	memory, ok := treeMemory.nodeMemory[nodeScope]
	if !ok {
		memory = map[string]any{}
		treeMemory.nodeMemory[nodeScope] = memory
	}

	return memory
}

func (blackboard *Blackboard) isNodeOpen(treeScope string, nodeScope string) bool {
	treeMemory := blackboard.getTreeMemory(treeScope)
	nodeMemory := blackboard.getNodeMemory(treeMemory, nodeScope)
	return boolValue(nodeMemory["isOpen"])
}

func (blackboard *Blackboard) setNodeOpen(treeScope string, nodeScope string, value bool) {
	treeMemory := blackboard.getTreeMemory(treeScope)
	nodeMemory := blackboard.getNodeMemory(treeMemory, nodeScope)
	nodeMemory["isOpen"] = value
}

func (blackboard *Blackboard) getOpenNodes(treeScope string) []Node {
	return blackboard.getTreeMemory(treeScope).openNodes
}

func (blackboard *Blackboard) setOpenNodes(treeScope string, openNodes []Node) {
	blackboard.getTreeMemory(treeScope).openNodes = openNodes
}

func (blackboard *Blackboard) setNodeCount(treeScope string, count int) {
	treeMemory := blackboard.getTreeMemory(treeScope)
	treeMemory.nodeCount = count
	treeMemory.memory["nodeCount"] = count
}

func (blackboard *Blackboard) Set(key string, value any, scope ...string) any {
	treeScope := ""
	nodeScope := ""
	if len(scope) > 0 {
		treeScope = scope[0]
	}
	if len(scope) > 1 {
		nodeScope = scope[1]
	}

	if treeScope == "" {
		blackboard.baseMemory[key] = value
		return value
	}

	treeMemory := blackboard.getTreeMemory(treeScope)
	if nodeScope != "" {
		nodeMemory := blackboard.getNodeMemory(treeMemory, nodeScope)
		nodeMemory[key] = value
		return value
	}

	switch key {
	case "nodeMemory":
		if typed, ok := value.(map[string]map[string]any); ok {
			treeMemory.nodeMemory = typed
			return value
		}
	case "openNodes":
		treeMemory.openNodes = toNodeSlice(value)
		return value
	case "traversalDepth":
		if typed, ok := toInt(value); ok {
			treeMemory.traversalDepth = typed
			return value
		}
	case "traversalCycle":
		if typed, ok := toInt(value); ok {
			treeMemory.traversalCycle = typed
			return value
		}
	case "nodeCount":
		if typed, ok := toInt(value); ok {
			treeMemory.nodeCount = typed
			treeMemory.memory[key] = typed
			return value
		}
	}

	treeMemory.memory[key] = value
	return value
}

func (blackboard *Blackboard) Get(key string, scope ...string) any {
	treeScope := ""
	nodeScope := ""
	if len(scope) > 0 {
		treeScope = scope[0]
	}
	if len(scope) > 1 {
		nodeScope = scope[1]
	}

	if treeScope == "" {
		return blackboard.baseMemory[key]
	}

	treeMemory := blackboard.getTreeMemory(treeScope)
	if nodeScope != "" {
		nodeMemory := blackboard.getNodeMemory(treeMemory, nodeScope)
		return nodeMemory[key]
	}

	switch key {
	case "nodeMemory":
		return treeMemory.nodeMemory
	case "openNodes":
		return treeMemory.openNodes
	case "traversalDepth":
		return treeMemory.traversalDepth
	case "traversalCycle":
		return treeMemory.traversalCycle
	case "nodeCount":
		return treeMemory.nodeCount
	default:
		return treeMemory.memory[key]
	}
}
