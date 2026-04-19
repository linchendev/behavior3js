package core

type BlackboardLike interface {
	Set(key string, value any, scope ...string) any
	Get(key string, scope ...string) any
}

type Node interface {
	Execute(tick *Tick) Status
	CloseNode(tick *Tick)
	GetBaseNode() *BaseNode
}

type NodeCallbacks interface {
	GetBaseNode() *BaseNode
	Enter(tick *Tick)
	Open(tick *Tick)
	Tick(tick *Tick) Status
	Close(tick *Tick)
	Exit(tick *Tick)
}

type treeMemory struct {
	memory         map[string]any
	nodeMemory     map[string]map[string]any
	openNodes      []Node
	traversalDepth int
	traversalCycle int
	nodeCount      int
}

type CompositeNode interface {
	Node
	GetChildren() []Node
	AppendChild(child Node)
}

type DecoratorNode interface {
	Node
	GetChild() Node
	SetChild(child Node)
}

type NodeConstructor func(properties map[string]any) (Node, error)

type LoadNodeConstructor func(spec NodeData) (Node, error)

type TreeData struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Root        string              `json:"root"`
	Properties  map[string]any      `json:"properties"`
	Nodes       map[string]NodeData `json:"nodes"`
	CustomNodes []CustomNodeData    `json:"custom_nodes"`
}

type NodeData struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Children    []string       `json:"children,omitempty"`
	Child       *string        `json:"child,omitempty"`
	Properties  map[string]any `json:"properties"`
	Parameters  map[string]any `json:"parameters"`
}

type CustomNodeData struct {
	Name     string `json:"name,omitempty"`
	Title    string `json:"title,omitempty"`
	Category string `json:"category,omitempty"`
}

var builtinConstructors = map[string]NodeConstructor{}
var builtinLoadConstructors = map[string]LoadNodeConstructor{}

func Register(name string, constructor NodeConstructor) {
	builtinConstructors[name] = constructor
}

func RegisterLoadConstructor(name string, constructor LoadNodeConstructor) {
	builtinLoadConstructors[name] = constructor
}

func IsBuiltinNode(name string) bool {
	_, ok := builtinConstructors[name]
	return ok
}

func NewBuiltinNode(name string, properties map[string]any) (Node, bool, error) {
	constructor, ok := builtinConstructors[name]
	if !ok {
		return nil, false, nil
	}

	node, err := constructor(properties)
	if err != nil {
		return nil, true, err
	}
	return node, true, nil
}

func newBuiltinNodeForLoad(spec NodeData) (Node, bool, error) {
	constructor, ok := builtinLoadConstructors[spec.Name]
	if ok {
		node, err := constructor(spec)
		if err != nil {
			return nil, true, err
		}
		return node, true, nil
	}

	node, found, err := NewBuiltinNode(spec.Name, clonePropertiesForLoad(spec.Properties))
	if err != nil || !found {
		return node, found, err
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
	if len(spec.Properties) > 0 {
		baseNode.Properties = clonePropertiesForLoad(spec.Properties)
	}

	return node, true, nil
}
