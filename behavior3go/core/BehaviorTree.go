package core

import "fmt"

type BehaviorTree struct {
	Id          string
	Title       string
	Description string
	Properties  map[string]any
	Root        Node
	Debug       any
}

func NewBehaviorTree() *BehaviorTree {
	return &BehaviorTree{
		Id:          createUUID(),
		Title:       "The behavior tree",
		Description: "Default description",
		Properties:  map[string]any{},
		Root:        nil,
		Debug:       nil,
	}
}

func (tree *BehaviorTree) Load(data *TreeData, names map[string]NodeConstructor) error {
	if data == nil {
		return nil
	}

	if data.Title != "" {
		tree.Title = data.Title
	}
	if data.Description != "" {
		tree.Description = data.Description
	}
	if data.Properties != nil {
		tree.Properties = copyMap(data.Properties)
	}

	nodes := make(map[string]Node, len(data.Nodes))

	for id, spec := range data.Nodes {
		var (
			node Node
			err  error
			ok   bool
		)

		if constructor, found := names[spec.Name]; found {
			properties := clonePropertiesForLoad(spec.Properties)
			node, err = constructor(properties)
			if err != nil {
				return err
			}
			baseNode := node.GetBaseNode()
			if spec.Id != "" {
				baseNode.Id = spec.Id
			}
			if spec.Title != "" {
				baseNode.Title = spec.Title
			}
			if spec.Description != "" {
				baseNode.Description = spec.Description
			}
			if properties != nil {
				baseNode.Properties = properties
			}
		} else {
			node, ok, err = newBuiltinNodeForLoad(spec)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("BehaviorTree.load: Invalid node name + %q.", spec.Name)
			}
		}

		baseNode := node.GetBaseNode()
		if spec.Id != "" {
			baseNode.Id = spec.Id
		}
		if spec.Title != "" {
			baseNode.Title = spec.Title
		}
		if spec.Description != "" {
			baseNode.Description = spec.Description
		}
		if len(spec.Properties) > 0 && baseNode.Properties == nil {
			baseNode.Properties = clonePropertiesForLoad(spec.Properties)
		}

		nodes[id] = node
	}

	for id, spec := range data.Nodes {
		node := nodes[id]

		switch typed := node.(type) {
		case CompositeNode:
			for _, childId := range spec.Children {
				if child, ok := nodes[childId]; ok {
					typed.AppendChild(child)
				}
			}
		case DecoratorNode:
			if spec.Child != nil {
				if child, ok := nodes[*spec.Child]; ok {
					typed.SetChild(child)
				}
			}
		}
	}

	tree.Root = nodes[data.Root]
	return nil
}

func (tree *BehaviorTree) Dump() *TreeData {
	data := &TreeData{
		Title:       tree.Title,
		Description: tree.Description,
		Properties:  copyMap(tree.Properties),
		Nodes:       map[string]NodeData{},
		CustomNodes: []CustomNodeData{},
	}

	if tree.Root == nil {
		return data
	}

	data.Root = tree.Root.GetBaseNode().Id
	customNames := map[string]bool{}
	stack := []Node{tree.Root}

	for len(stack) > 0 {
		lastIndex := len(stack) - 1
		node := stack[lastIndex]
		stack = stack[:lastIndex]

		baseNode := node.GetBaseNode()
		spec := NodeData{
			Id:          baseNode.Id,
			Name:        baseNode.Name,
			Title:       baseNode.Title,
			Description: baseNode.Description,
			Properties:  copyMap(baseNode.Properties),
			Parameters:  copyMap(baseNode.Parameters),
		}
		if spec.Properties == nil {
			spec.Properties = map[string]any{}
		}
		if spec.Parameters == nil {
			spec.Parameters = map[string]any{}
		}

		nodeName := baseNode.Name
		if !IsBuiltinNode(nodeName) && !customNames[nodeName] {
			customNames[nodeName] = true
			data.CustomNodes = append(data.CustomNodes, CustomNodeData{
				Name:     nodeName,
				Title:    baseNode.Title,
				Category: baseNode.Category,
			})
		}

		switch typed := node.(type) {
		case CompositeNode:
			children := typed.GetChildren()
			spec.Children = make([]string, 0, len(children))
			for i := len(children) - 1; i >= 0; i-- {
				spec.Children = append(spec.Children, children[i].GetBaseNode().Id)
				stack = append(stack, children[i])
			}
		case DecoratorNode:
			child := typed.GetChild()
			if child != nil {
				childId := child.GetBaseNode().Id
				spec.Child = &childId
				stack = append(stack, child)
			}
		}

		data.Nodes[baseNode.Id] = spec
	}

	return data
}

func toNodeSlice(value any) []Node {
	if value == nil {
		return []Node{}
	}

	if nodes, ok := value.([]Node); ok {
		return nodes
	}

	if nodes, ok := value.([]any); ok {
		result := make([]Node, 0, len(nodes))
		for _, item := range nodes {
			node, ok := item.(Node)
			if ok {
				result = append(result, node)
			}
		}
		return result
	}

	return []Node{}
}

func (tree *BehaviorTree) Tick(target any, blackboard BlackboardLike) Status {
	if blackboard == nil {
		panic("The blackboard parameter is obligatory and must be an instance of b3.Blackboard")
	}
	if tree.Root == nil {
		return 0
	}

	tick := acquireTick(tree, tree.Debug, target, blackboard)
	defer releaseTick(tick)

	state := tree.Root.Execute(tick)

	if tick.blackboard != nil {
		lastOpenNodes := tick.blackboard.getOpenNodes(tree.Id)
		currOpenNodes := snapshotOpenNodesForClose(lastOpenNodes, tick.OpenNodes)
		tree.closeStaleOpenNodes(tick, lastOpenNodes, currOpenNodes)
		if tree.shouldReuseOpenNodes(lastOpenNodes, currOpenNodes) {
			tick.blackboard.setNodeCount(tree.Id, tick.NodeCount)
			return state
		}
		tick.blackboard.setOpenNodes(tree.Id, cloneNodes(currOpenNodes))
		tick.blackboard.setNodeCount(tree.Id, tick.NodeCount)
		return state
	}

	lastOpenNodes := toNodeSlice(blackboard.Get("openNodes", tree.Id))
	currOpenNodes := snapshotOpenNodesForClose(lastOpenNodes, tick.OpenNodes)
	tree.closeStaleOpenNodes(tick, lastOpenNodes, currOpenNodes)
	if tree.shouldReuseOpenNodes(lastOpenNodes, currOpenNodes) {
		blackboard.Set("nodeCount", tick.NodeCount, tree.Id)
		return state
	}
	blackboard.Set("openNodes", cloneNodes(currOpenNodes), tree.Id)
	blackboard.Set("nodeCount", tick.NodeCount, tree.Id)

	return state
}

func (tree *BehaviorTree) closeStaleOpenNodes(tick *Tick, lastOpenNodes []Node, currOpenNodes []Node) {
	start := sharedOpenPrefixLength(lastOpenNodes, currOpenNodes)
	for i := len(lastOpenNodes) - 1; i >= start; i-- {
		lastOpenNodes[i].CloseNode(tick)
	}
}

func (tree *BehaviorTree) shouldReuseOpenNodes(lastOpenNodes []Node, currOpenNodes []Node) bool {
	if len(lastOpenNodes) != len(currOpenNodes) {
		return false
	}
	if len(lastOpenNodes) == 0 {
		return true
	}
	return lastOpenNodes[len(lastOpenNodes)-1] == currOpenNodes[len(currOpenNodes)-1]
}

func sharedOpenPrefixLength(lastOpenNodes []Node, currOpenNodes []Node) int {
	limit := len(lastOpenNodes)
	if len(currOpenNodes) < limit {
		limit = len(currOpenNodes)
	}

	index := 0
	for index < limit && lastOpenNodes[index] == currOpenNodes[index] {
		index++
	}
	return index
}

func snapshotOpenNodesForClose(lastOpenNodes []Node, currOpenNodes []Node) []Node {
	sharedPrefix := sharedOpenPrefixLength(lastOpenNodes, currOpenNodes)
	if sharedPrefix < len(lastOpenNodes) {
		return cloneNodes(currOpenNodes)
	}
	return currOpenNodes
}

func cloneNodes(nodes []Node) []Node {
	if len(nodes) == 0 {
		return []Node{}
	}

	cloned := make([]Node, len(nodes))
	copy(cloned, nodes)
	return cloned
}
