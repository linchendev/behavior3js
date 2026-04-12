package core

import "sync"

type Tick struct {
	Tree       *BehaviorTree
	Debug      any
	Target     any
	Blackboard BlackboardLike
	OpenNodes  []Node
	NodeCount  int
}

var tickPool = sync.Pool{
	New: func() any {
		return &Tick{
			OpenNodes: make([]Node, 0, 16),
		}
	},
}

func NewTick() *Tick {
	return &Tick{
		OpenNodes: make([]Node, 0, 16),
		NodeCount: 0,
	}
}

func acquireTick(tree *BehaviorTree, debug any, target any, blackboard BlackboardLike) *Tick {
	tick := tickPool.Get().(*Tick)
	tick.Reset(tree, debug, target, blackboard)
	return tick
}

func releaseTick(tick *Tick) {
	tick.Reset(nil, nil, nil, nil)
	tickPool.Put(tick)
}

func (tick *Tick) Reset(tree *BehaviorTree, debug any, target any, blackboard BlackboardLike) {
	tick.Tree = tree
	tick.Debug = debug
	tick.Target = target
	tick.Blackboard = blackboard
	tick.OpenNodes = tick.OpenNodes[:0]
	tick.NodeCount = 0
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
