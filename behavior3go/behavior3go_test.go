package behavior3go_test

import (
	"testing"

	b3 "github.com/behavior3/behavior3go"
)

func stringPtr(value string) *string {
	return &value
}

func TestLoadJSONWithDefaultNodes(t *testing.T) {
	tree := b3.NewBehaviorTree()

	data := &b3.TreeData{
		Title:       "A JSON Behavior Tree",
		Description: "This description",
		Root:        "1",
		Properties: map[string]any{
			"variable": "value",
		},
		Nodes: map[string]b3.NodeData{
			"1": {
				Id:          "1",
				Name:        "Priority",
				Title:       "Root Node",
				Description: "Root Description",
				Children:    []string{"2", "3"},
				Properties: map[string]any{
					"var1": 123,
					"composite": map[string]any{
						"var2": true,
						"var3": "value",
					},
				},
			},
			"2": {
				Name:        "Inverter",
				Title:       "Node 1",
				Description: "Node 1 Description",
				Child:       stringPtr("4"),
			},
			"3": {
				Name:        "MemSequence",
				Title:       "Node 2",
				Description: "Node 2 Description",
				Children:    []string{},
			},
			"4": {
				Name:        "MaxTime",
				Title:       "Node 3",
				Description: "Node 3 Description",
				Properties: map[string]any{
					"maxTime": 1,
				},
				Parameters: map[string]any{
					"maxTime": 999,
				},
			},
		},
	}

	if err := tree.Load(data, nil); err != nil {
		t.Fatalf("load should succeed: %v", err)
	}

	root, ok := tree.Root.(*b3.Priority)
	if !ok {
		t.Fatal("root should be priority")
	}
	if tree.Title != "A JSON Behavior Tree" || tree.Description != "This description" {
		t.Fatal("tree metadata mismatch")
	}
	if root.Id != "1" || len(root.Children) != 2 {
		t.Fatal("root shape mismatch")
	}

	node1, ok := root.Children[0].(*b3.Inverter)
	if !ok {
		t.Fatal("first child should be inverter")
	}
	node2, ok := root.Children[1].(*b3.MemSequence)
	if !ok {
		t.Fatal("second child should be mem sequence")
	}
	node3, ok := node1.Child.(*b3.MaxTime)
	if !ok {
		t.Fatal("nested child should be max time")
	}

	if node1.Title != "Node 1" || node2.Title != "Node 2" || node3.Title != "Node 3" {
		t.Fatal("node metadata mismatch")
	}
	if node3.Parameters["maxTime"] == 999 {
		t.Fatal("properties should not overwrite deprecated parameters")
	}
}

func TestLoadJSONWithCustomNodes(t *testing.T) {
	tree := b3.NewBehaviorTree()

	data := &b3.TreeData{
		Title:       "A JSON Behavior Tree",
		Description: "This descriptions",
		Root:        "1",
		Nodes: map[string]b3.NodeData{
			"1": {
				Name:        "Priority",
				Title:       "Root Node",
				Description: "Root Description",
				Children:    []string{"2"},
			},
			"2": {
				Name:        "Condition",
				Title:       "Node 2",
				Description: "Node 2 Description",
			},
		},
	}

	names := map[string]b3.NodeConstructor{
		"Condition": func(properties map[string]any) (b3.Node, error) {
			return b3.NewCondition(), nil
		},
	}

	if err := tree.Load(data, names); err != nil {
		t.Fatalf("custom load should succeed: %v", err)
	}

	root, ok := tree.Root.(*b3.Priority)
	if !ok || len(root.Children) != 1 {
		t.Fatal("root should contain one custom child")
	}
	if _, ok := root.Children[0].(*b3.Condition); !ok {
		t.Fatal("custom child should resolve through names registry")
	}
}

func TestDumpJSONModel(t *testing.T) {
	tree := b3.NewBehaviorTree()
	tree.Properties = map[string]any{
		"prop": "value",
		"comp": map[string]any{
			"val1": 234,
			"val2": "value",
		},
	}

	node5 := b3.NewCondition()
	node5.Id = "node-5"
	node5.Title = "Node5"
	node5.Description = "Node 5 Description"

	node4 := b3.NewWait(0)
	node4.Id = "node-4"
	node4.Title = "Node4"
	node4.Description = "Node 4 Description"

	node3 := b3.NewMemSequence(node5)
	node3.Id = "node-3"
	node3.Title = "Node3"
	node3.Description = "Node 3 Description"

	node2 := b3.NewInverter(node4)
	node2.Id = "node-2"
	node2.Title = "Node2"
	node2.Description = "Node 2 Description"

	node1 := b3.NewPriority(node2, node3)
	node1.Id = "node-1"
	node1.Title = "Node1"
	node1.Description = "Node 1 Description"
	node1.Properties = map[string]any{"key": "value"}

	tree.Root = node1
	tree.Title = "Title in Tree"
	tree.Description = "Tree Description"

	data := tree.Dump()

	if data.Title != "Title in Tree" || data.Description != "Tree Description" || data.Root != "node-1" {
		t.Fatal("tree dump metadata mismatch")
	}
	if len(data.CustomNodes) != 1 || data.CustomNodes[0].Name != "Condition" {
		t.Fatal("dump should include custom condition node")
	}
	if data.Nodes["node-1"].Children[0] != "node-3" || data.Nodes["node-1"].Children[1] != "node-2" {
		t.Fatal("dump should preserve JS child ordering")
	}
	if data.Nodes["node-2"].Child == nil {
		t.Fatal("decorator child should be dumped")
	}
	if data.Nodes["node-4"].Child != nil || len(data.Nodes["node-4"].Children) != 0 {
		t.Fatal("action node should not dump child relationships")
	}
}

func TestLoadKeepsNodePropertiesIsolatedFromInput(t *testing.T) {
	tree := b3.NewBehaviorTree()
	data := &b3.TreeData{
		Title: "Wait Tree",
		Root:  "1",
		Nodes: map[string]b3.NodeData{
			"1": {
				Id:    "1",
				Name:  "Wait",
				Title: "Wait Node",
				Properties: map[string]any{
					"milliseconds": 100,
				},
			},
		},
	}

	if err := tree.Load(data, nil); err != nil {
		t.Fatalf("load should succeed: %v", err)
	}

	data.Nodes["1"].Properties["milliseconds"] = 250

	waitNode, ok := tree.Root.(*b3.Wait)
	if !ok {
		t.Fatal("root should be wait")
	}
	if waitNode.Properties["milliseconds"] != 100 {
		t.Fatal("loaded node properties should not be mutated through input TreeData reuse")
	}
}

func TestLoadInitializesWritableParameters(t *testing.T) {
	tree := b3.NewBehaviorTree()
	data := &b3.TreeData{
		Title: "Wait Tree",
		Root:  "1",
		Nodes: map[string]b3.NodeData{
			"1": {
				Id:    "1",
				Name:  "Wait",
				Title: "Wait Node",
				Properties: map[string]any{
					"milliseconds": 100,
				},
			},
		},
	}

	if err := tree.Load(data, nil); err != nil {
		t.Fatalf("load should succeed: %v", err)
	}

	waitNode, ok := tree.Root.(*b3.Wait)
	if !ok {
		t.Fatal("root should be wait")
	}

	waitNode.Parameters["custom"] = "value"
	if waitNode.Parameters["custom"] != "value" {
		t.Fatal("loaded node parameters should remain writable")
	}
}
