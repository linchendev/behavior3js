package core

type Blackboard struct {
	baseMemory map[string]any
	treeMemory map[string]map[string]any
}

func NewBlackboard() *Blackboard {
	return &Blackboard{
		baseMemory: map[string]any{},
		treeMemory: map[string]map[string]any{},
	}
}

func (blackboard *Blackboard) getTreeMemory(treeScope string) map[string]any {
	memory, ok := blackboard.treeMemory[treeScope]
	if !ok {
		memory = map[string]any{
			"nodeMemory":     map[string]map[string]any{},
			"openNodes":      []Node{},
			"traversalDepth": 0,
			"traversalCycle": 0,
		}
		blackboard.treeMemory[treeScope] = memory
	}

	return memory
}

func (blackboard *Blackboard) getNodeMemory(treeMemory map[string]any, nodeScope string) map[string]any {
	nodeMemory := treeMemory["nodeMemory"].(map[string]map[string]any)
	memory, ok := nodeMemory[nodeScope]
	if !ok {
		memory = map[string]any{}
		nodeMemory[nodeScope] = memory
	}

	return memory
}

func (blackboard *Blackboard) getMemory(treeScope string, nodeScope string) map[string]any {
	memory := blackboard.baseMemory

	if treeScope != "" {
		memory = blackboard.getTreeMemory(treeScope)
		if nodeScope != "" {
			memory = blackboard.getNodeMemory(memory, nodeScope)
		}
	}

	return memory
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

	memory := blackboard.getMemory(treeScope, nodeScope)
	memory[key] = value
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

	memory := blackboard.getMemory(treeScope, nodeScope)
	return memory[key]
}
