package core

type Tick struct {
	Tree       *BehaviorTree
	Debug      any
	Target     any
	Blackboard BlackboardLike
	OpenNodes  []Node
	NodeCount  int
}

func NewTick() *Tick {
	return &Tick{
		OpenNodes: []Node{},
		NodeCount: 0,
	}
}

func (tick *Tick) EnterNode(node Node) {
	tick.NodeCount++
	tick.OpenNodes = append(tick.OpenNodes, node)
}

func (tick *Tick) OpenNode(node Node) {}

func (tick *Tick) TickNode(node Node) {}

func (tick *Tick) CloseNode(node Node) {
	if len(tick.OpenNodes) == 0 {
		return
	}
	tick.OpenNodes = tick.OpenNodes[:len(tick.OpenNodes)-1]
}

func (tick *Tick) ExitNode(node Node) {}
