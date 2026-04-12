package actions

import (
	"testing"

	"github.com/behavior3/behavior3go/core"
)

func newTick() *core.Tick {
	tick := core.NewTick()
	tick.Tree = core.NewBehaviorTree()
	tick.Tree.Id = "tree1"
	tick.Blackboard = core.NewBlackboard()
	return tick
}

func TestBasicActions(t *testing.T) {
	tests := []struct {
		name     string
		nodeName string
		status   core.Status
		node     core.Node
	}{
		{name: "Error", nodeName: "Error", status: core.ERROR, node: NewError()},
		{name: "Failer", nodeName: "Failer", status: core.FAILURE, node: NewFailer()},
		{name: "Runner", nodeName: "Runner", status: core.RUNNING, node: NewRunner()},
		{name: "Succeeder", nodeName: "Succeeder", status: core.SUCCESS, node: NewSucceeder()},
	}

	for _, test := range tests {
		if test.node.GetBaseNode().Name != test.nodeName {
			t.Fatalf("%s name mismatch", test.name)
		}
		if status := test.node.Execute(newTick()); status != test.status {
			t.Fatalf("%s status mismatch", test.name)
		}
	}
}

func TestWait(t *testing.T) {
	originalNow := core.Now
	defer func() { core.Now = originalNow }()

	wait := NewWait(100)
	wait.Id = "node1"

	tick := newTick()

	core.Now = func() int64 { return 1000 }
	wait.Open(tick)
	if tick.Blackboard.Get("startTime", "tree1", "node1") != int64(1000) {
		t.Fatal("wait should persist start time")
	}

	core.Now = func() int64 { return 1099 }
	if status := wait.Tick(tick); status != core.RUNNING {
		t.Fatal("wait should still be running")
	}

	core.Now = func() int64 { return 1101 }
	if status := wait.Tick(tick); status != core.SUCCESS {
		t.Fatal("wait should succeed after duration")
	}
}
