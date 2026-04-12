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
	if names == nil {
		names = map[string]NodeConstructor{}
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

	nodes := map[string]Node{}

	for id, spec := range data.Nodes {
		var (
			node Node
			err  error
			ok   bool
		)

		if constructor, found := names[spec.Name]; found {
			node, err = constructor(copyMap(spec.Properties))
			if err != nil {
				return err
			}
		} else {
			node, ok, err = NewBuiltinNode(spec.Name, copyMap(spec.Properties))
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
		if spec.Properties != nil {
			baseNode.Properties = copyMap(spec.Properties)
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

	tick := NewTick()
	tick.Debug = tree.Debug
	tick.Target = target
	tick.Blackboard = blackboard
	tick.Tree = tree

	state := tree.Root.Execute(tick)

	lastOpenNodes := toNodeSlice(blackboard.Get("openNodes", tree.Id))
	currOpenNodes := append([]Node{}, tick.OpenNodes...)

	start := 0
	limit := len(lastOpenNodes)
	if len(currOpenNodes) < limit {
		limit = len(currOpenNodes)
	}

	for i := 0; i < limit; i++ {
		start = i + 1
		if lastOpenNodes[i] != currOpenNodes[i] {
			break
		}
	}

	for i := len(lastOpenNodes) - 1; i >= start; i-- {
		lastOpenNodes[i].CloseNode(tick)
	}

	blackboard.Set("openNodes", currOpenNodes, tree.Id)
	blackboard.Set("nodeCount", tick.NodeCount, tree.Id)

	return state
}
