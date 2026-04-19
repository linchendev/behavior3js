package core

import "testing"

type blackboardCall struct {
	key   string
	value any
	scope []string
}

type blackboardStub struct {
	setCalls []blackboardCall
	getFunc  func(key string, scope ...string) any
}

func (blackboard *blackboardStub) Set(key string, value any, scope ...string) any {
	blackboard.setCalls = append(blackboard.setCalls, blackboardCall{
		key:   key,
		value: value,
		scope: append([]string{}, scope...),
	})
	return value
}

func (blackboard *blackboardStub) Get(key string, scope ...string) any {
	if blackboard.getFunc == nil {
		return nil
	}
	return blackboard.getFunc(key, scope...)
}

func (blackboard *blackboardStub) hasSetCall(key string, value any, scope ...string) bool {
	for _, call := range blackboard.setCalls {
		if call.key != key {
			continue
		}
		if len(call.scope) != len(scope) {
			continue
		}
		matched := true
		for index := range scope {
			if call.scope[index] != scope[index] {
				matched = false
				break
			}
		}
		if matched && call.value == value {
			return true
		}
	}
	return false
}

type testNode struct {
	BaseNode
	tickStatus  Status
	enterCount  int
	openCount   int
	tickCount   int
	closeCount  int
	exitCount   int
	executeFunc func(tick *Tick) Status
}

func newTestNode() *testNode {
	node := &testNode{tickStatus: SUCCESS}
	node.BaseNode = NewBaseNode(BaseNodeOptions{Name: "TestNode"})
	return node
}

func (node *testNode) Execute(tick *Tick) Status {
	if node.executeFunc != nil {
		return node.executeFunc(tick)
	}
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *testNode) Enter(tick *Tick) { node.enterCount++ }
func (node *testNode) Open(tick *Tick)  { node.openCount++ }
func (node *testNode) Tick(tick *Tick) Status {
	node.tickCount++
	return node.tickStatus
}
func (node *testNode) Close(tick *Tick) { node.closeCount++ }
func (node *testNode) Exit(tick *Tick)  { node.exitCount++ }

type closableNode struct {
	BaseNode
	closeCount int
}

type runningBranchNode struct {
	BaseNode
	branch Node
}

func newClosableNode(id string) *closableNode {
	node := &closableNode{}
	node.BaseNode = NewBaseNode(BaseNodeOptions{Name: id})
	node.Id = id
	return node
}

func newRunningBranchNode(id string, branch Node) *runningBranchNode {
	node := &runningBranchNode{branch: branch}
	node.BaseNode = NewBaseNode(BaseNodeOptions{Name: id})
	node.Id = id
	return node
}

func (node *closableNode) Execute(tick *Tick) Status {
	tick.EnterNode(node)
	return SUCCESS
}

func (node *closableNode) CloseNode(tick *Tick) {
	node.closeCount++
}

func (node *runningBranchNode) Execute(tick *Tick) Status {
	return node.GetBaseNode().ExecuteNode(tick, node, node)
}

func (node *runningBranchNode) Tick(tick *Tick) Status {
	tick.EnterNode(node.branch)
	return RUNNING
}

func TestBlackboard(t *testing.T) {
	blackboard := NewBlackboard()

	blackboard.Set("var1", "this is some value")
	blackboard.Set("var2", 999888)
	if blackboard.Get("var1") != "this is some value" {
		t.Fatal("global scope write mismatch")
	}
	if blackboard.Get("var2") != 999888 {
		t.Fatal("global scope integer mismatch")
	}
	if blackboard.Get("var3") != nil {
		t.Fatal("unexpected key should be nil")
	}

	blackboard.Set("var1", "value", "tree1")
	if blackboard.Get("nodeMemory", "tree1") == nil {
		t.Fatal("tree memory should initialize nodeMemory")
	}
	if blackboard.Get("openNodes", "tree1") == nil {
		t.Fatal("tree memory should initialize openNodes")
	}
	if blackboard.Get("traversalCycle", "tree1") == nil {
		t.Fatal("tree memory should initialize traversalCycle")
	}

	blackboard.Set("var1", "value 1", "tree 1")
	blackboard.Set("var2", "value 2", "tree 1", "node 1")
	blackboard.Set("var3", "value 3", "tree 1", "node 2")
	if blackboard.Get("var2", "tree 1", "node 1") != "value 2" {
		t.Fatal("node scope read mismatch")
	}
	if blackboard.Get("var3", "tree 1", "node 1") != nil {
		t.Fatal("node scopes should be isolated")
	}
}

func TestTick(t *testing.T) {
	tick := NewTick()
	if tick.NodeCount != 0 || len(tick.OpenNodes) != 0 {
		t.Fatal("tick should initialize empty")
	}

	node := newClosableNode("node1")
	tick.EnterNode(node)
	if tick.NodeCount != 1 || len(tick.OpenNodes) != 1 {
		t.Fatal("tick enter should update state")
	}

	tick.CloseNode(node)
	if len(tick.OpenNodes) != 0 {
		t.Fatal("tick close should pop open nodes")
	}
}

func TestBaseNodeExecute(t *testing.T) {
	blackboard := &blackboardStub{
		getFunc: func(key string, scope ...string) any {
			return false
		},
	}
	tree := NewBehaviorTree()
	tree.Id = "tree1"
	tick := NewTick()
	tick.Tree = tree
	tick.Blackboard = blackboard

	node := newTestNode()
	node.Id = "node1"

	status := node.Execute(tick)
	if status != SUCCESS {
		t.Fatal("base node should return tick status")
	}
	if !blackboard.hasSetCall("isOpen", true, "tree1", "node1") {
		t.Fatal("node should open before tick")
	}
	if !blackboard.hasSetCall("isOpen", false, "tree1", "node1") {
		t.Fatal("node should close after non-running tick")
	}
	if node.enterCount != 1 || node.openCount != 1 || node.tickCount != 1 || node.closeCount != 1 || node.exitCount != 1 {
		t.Fatal("base node lifecycle mismatch")
	}
	if tick.NodeCount != 1 {
		t.Fatal("tick should record entered node")
	}
}

func TestBaseNodeRunningDoesNotClose(t *testing.T) {
	blackboard := &blackboardStub{
		getFunc: func(key string, scope ...string) any {
			return false
		},
	}
	tree := NewBehaviorTree()
	tree.Id = "tree1"
	tick := NewTick()
	tick.Tree = tree
	tick.Blackboard = blackboard

	node := newTestNode()
	node.Id = "node1"
	node.tickStatus = RUNNING
	node.Execute(tick)

	if node.closeCount != 0 {
		t.Fatal("running node should not close")
	}
	if blackboard.hasSetCall("isOpen", false, "tree1", "node1") {
		t.Fatal("running node should remain open")
	}
}

func TestBehaviorTreeTick(t *testing.T) {
	tree := NewBehaviorTree()
	tree.Id = "tree1"

	blackboard := &blackboardStub{
		getFunc: func(key string, scope ...string) any {
			if key == "openNodes" && len(scope) == 1 && scope[0] == "tree1" {
				return []Node{}
			}
			return nil
		},
	}

	root := newTestNode()
	root.executeFunc = func(tick *Tick) Status {
		tick.EnterNode(newClosableNode("node1"))
		tick.EnterNode(newClosableNode("node2"))
		return SUCCESS
	}

	tree.Root = root
	tree.Tick(nil, blackboard)

	foundOpenNodes := false
	foundNodeCount := false
	for _, call := range blackboard.setCalls {
		if call.key == "openNodes" && len(call.scope) == 1 && call.scope[0] == "tree1" {
			nodes, ok := call.value.([]Node)
			if ok && len(nodes) == 2 {
				foundOpenNodes = true
			}
		}
		if call.key == "nodeCount" && call.value == 2 {
			foundNodeCount = true
		}
	}

	if !foundOpenNodes || !foundNodeCount {
		t.Fatal("tree should populate blackboard with open nodes and node count")
	}
}

func TestBehaviorTreeClosesOpenedNodes(t *testing.T) {
	tree := NewBehaviorTree()
	tree.Id = "tree1"

	node1 := newClosableNode("node1")
	node2 := newClosableNode("node2")
	node3 := newClosableNode("node3")
	node4 := newClosableNode("node4")
	node5 := newClosableNode("node5")

	root := newTestNode()
	root.executeFunc = func(tick *Tick) Status {
		tick.EnterNode(node1)
		tick.EnterNode(node2)
		tick.EnterNode(node3)
		return SUCCESS
	}
	tree.Root = root

	blackboard := &blackboardStub{
		getFunc: func(key string, scope ...string) any {
			if key == "openNodes" && len(scope) == 1 && scope[0] == "tree1" {
				return []Node{node1, node2, node3, node4, node5}
			}
			return nil
		},
	}

	tree.Tick(nil, blackboard)

	if node4.closeCount != 1 || node5.closeCount != 1 {
		t.Fatal("stale open nodes should close")
	}
	if node1.closeCount != 0 || node2.closeCount != 0 || node3.closeCount != 0 {
		t.Fatal("shared open path should remain open")
	}
}

func TestBehaviorTreePreservesRunningBranchSwitchOpenNodes(t *testing.T) {
	tree := NewBehaviorTree()
	tree.Id = "tree1"

	blackboard := NewBlackboard()
	oldBranch := newClosableNode("old")
	newBranch := newClosableNode("new")
	root := newRunningBranchNode("root", newBranch)

	tree.Root = root
	blackboard.Set("openNodes", []Node{root, oldBranch}, "tree1")

	tree.Tick(nil, blackboard)

	openNodes, ok := blackboard.Get("openNodes", "tree1").([]Node)
	if !ok {
		t.Fatal("open nodes should be stored as []Node")
	}
	if len(openNodes) != 2 {
		t.Fatalf("expected two open nodes after branch switch, got %d", len(openNodes))
	}
	if openNodes[0] != root || openNodes[1] != newBranch {
		t.Fatal("running branch switch should keep the new open branch")
	}
	if oldBranch.closeCount != 1 {
		t.Fatal("stale running branch should close exactly once")
	}
}

func TestBehaviorTreePreservesRunningPrefixWhenPreviousPathIsLonger(t *testing.T) {
	tree := NewBehaviorTree()
	tree.Id = "tree1"

	blackboard := NewBlackboard()
	sharedBranch := newClosableNode("shared")
	staleLeaf := newClosableNode("stale")
	root := newRunningBranchNode("root", sharedBranch)

	tree.Root = root
	blackboard.Set("openNodes", []Node{root, sharedBranch, staleLeaf}, "tree1")

	tree.Tick(nil, blackboard)

	openNodes, ok := blackboard.Get("openNodes", "tree1").([]Node)
	if !ok {
		t.Fatal("open nodes should be stored as []Node")
	}
	if len(openNodes) != 2 {
		t.Fatalf("expected two open nodes after prefix shrink, got %d", len(openNodes))
	}
	if openNodes[0] != root || openNodes[1] != sharedBranch {
		t.Fatal("current running prefix should remain intact when closing stale suffix")
	}
	if staleLeaf.closeCount != 1 {
		t.Fatal("stale suffix node should close exactly once")
	}
}
