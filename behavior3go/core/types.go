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

func Register(name string, constructor NodeConstructor) {
	builtinConstructors[name] = constructor
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
