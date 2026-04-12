package decorators

import (
	"testing"

	"github.com/behavior3/behavior3go/core"
)

type stubNode struct {
	base      core.BaseNode
	statuses  []core.Status
	callCount int
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

func (node *stubNode) CloseNode(tick *core.Tick) {}

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

func TestInverter(t *testing.T) {
	inverter := NewInverter(newStubNode(core.SUCCESS))
	if inverter.Name != "Inverter" {
		t.Fatal("inverter name mismatch")
	}
	if status := inverter.Execute(newTick()); status != core.FAILURE {
		t.Fatal("inverter should flip success to failure")
	}

	inverter.Child = newStubNode(core.FAILURE)
	if status := inverter.Execute(newTick()); status != core.SUCCESS {
		t.Fatal("inverter should flip failure to success")
	}
}

func TestLimiter(t *testing.T) {
	child := newStubNode(core.SUCCESS)
	limiter := NewLimiter(1, child)
	limiter.Id = "node1"

	tick := newTick()
	if status := limiter.Execute(tick); status != core.SUCCESS {
		t.Fatal("limiter should allow first execution")
	}
	tick.Blackboard.Set("i", 1, "tree1", "node1")
	if status := limiter.Tick(tick); status != core.FAILURE {
		t.Fatal("limiter should stop after max loop")
	}
	if child.callCount != 1 {
		t.Fatal("limiter should not execute child after limit")
	}
}

func TestMaxTime(t *testing.T) {
	originalNow := core.Now
	defer func() { core.Now = originalNow }()

	child := newStubNode(core.RUNNING)
	maxTime := NewMaxTime(15, child)
	maxTime.Id = "node1"
	tick := newTick()

	core.Now = func() int64 { return 1000 }
	maxTime.Open(tick)

	core.Now = func() int64 { return 1014 }
	if status := maxTime.Tick(tick); status != core.RUNNING {
		t.Fatal("max time should still be running")
	}

	core.Now = func() int64 { return 1016 }
	if status := maxTime.Tick(tick); status != core.FAILURE {
		t.Fatal("max time should fail after timeout")
	}
}

func TestRepeaterFamily(t *testing.T) {
	repeater := NewRepeater(3, newStubNode(core.SUCCESS))
	repeater.Id = "node1"
	if status := repeater.Execute(newTick()); status != core.SUCCESS {
		t.Fatal("repeater should finish with child status")
	}

	repeatUntilFailure := NewRepeatUntilFailure(10, newStubNode(core.SUCCESS, core.SUCCESS, core.FAILURE))
	if status := repeatUntilFailure.Execute(newTick()); status != core.FAILURE {
		t.Fatal("repeat until failure should stop at failure")
	}

	repeatUntilSuccess := NewRepeatUntilSuccess(10, newStubNode(core.FAILURE, core.FAILURE, core.SUCCESS))
	if status := repeatUntilSuccess.Execute(newTick()); status != core.SUCCESS {
		t.Fatal("repeat until success should stop at success")
	}
}
