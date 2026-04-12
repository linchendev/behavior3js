package composites

import (
	"testing"

	"github.com/behavior3/behavior3go/core"
)

type stubNode struct {
	base       core.BaseNode
	statuses   []core.Status
	callCount  int
	closeCount int
}

func newStubNode(statuses ...core.Status) *stubNode {
	node := &stubNode{statuses: statuses}
	node.base = core.NewBaseNode(core.BaseNodeOptions{Name: "Stub"})
	return node
}

func (node *stubNode) Execute(tick *core.Tick) core.Status {
	index := node.callCount
	node.callCount++
	if index >= len(node.statuses) {
		return node.statuses[len(node.statuses)-1]
	}
	return node.statuses[index]
}

func (node *stubNode) CloseNode(tick *core.Tick) {
	node.closeCount++
}

func (node *stubNode) GetBaseNode() *core.BaseNode {
	return &node.base
}

func newTick() *core.Tick {
	tick := core.NewTick()
	tick.Tree = core.NewBehaviorTree()
	tick.Tree.Id = "tree1"
	tick.Blackboard = core.NewBlackboard()
	return tick
}

func TestSequenceAndPriority(t *testing.T) {
	sequence := NewSequence(
		newStubNode(core.SUCCESS),
		newStubNode(core.SUCCESS),
		newStubNode(core.FAILURE),
	)
	if sequence.Name != "Sequence" {
		t.Fatal("sequence name mismatch")
	}
	if status := sequence.Tick(newTick()); status != core.FAILURE {
		t.Fatal("sequence should stop on first non-success")
	}

	priority := NewPriority(
		newStubNode(core.FAILURE),
		newStubNode(core.RUNNING),
		newStubNode(core.SUCCESS),
	)
	if priority.Name != "Priority" {
		t.Fatal("priority name mismatch")
	}
	if status := priority.Tick(newTick()); status != core.RUNNING {
		t.Fatal("priority should stop on first non-failure")
	}
}

func TestMemPriority(t *testing.T) {
	node1 := newStubNode(core.FAILURE)
	node2 := newStubNode(core.FAILURE)
	node3 := newStubNode(core.RUNNING, core.FAILURE)
	node4 := newStubNode(core.SUCCESS)

	memPriority := NewMemPriority(node1, node2, node3, node4)
	memPriority.Id = "node1"

	tick := newTick()
	if status := memPriority.Execute(tick); status != core.RUNNING {
		t.Fatal("mem priority should return running on first tick")
	}

	if tick.Blackboard.Get("runningChild", "tree1", "node1") != 2 {
		t.Fatal("mem priority should remember running child")
	}

	if status := memPriority.Execute(tick); status != core.SUCCESS {
		t.Fatal("mem priority should resume from remembered child")
	}

	if node1.callCount != 1 || node2.callCount != 1 || node3.callCount != 2 || node4.callCount != 1 {
		t.Fatal("mem priority resume path mismatch")
	}
}

func TestMemSequence(t *testing.T) {
	node1 := newStubNode(core.SUCCESS)
	node2 := newStubNode(core.SUCCESS)
	node3 := newStubNode(core.RUNNING, core.SUCCESS)
	node4 := newStubNode(core.FAILURE)

	memSequence := NewMemSequence(node1, node2, node3, node4)
	memSequence.Id = "node1"

	tick := newTick()
	if status := memSequence.Execute(tick); status != core.RUNNING {
		t.Fatal("mem sequence should return running on first tick")
	}

	if tick.Blackboard.Get("runningChild", "tree1", "node1") != 2 {
		t.Fatal("mem sequence should remember running child")
	}

	if status := memSequence.Execute(tick); status != core.FAILURE {
		t.Fatal("mem sequence should resume from remembered child")
	}

	if node1.callCount != 1 || node2.callCount != 1 || node3.callCount != 2 || node4.callCount != 1 {
		t.Fatal("mem sequence resume path mismatch")
	}
}
