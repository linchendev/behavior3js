package behavior3go_test

import (
	"os"
	"path/filepath"
	"testing"

	b3 "github.com/behavior3/behavior3go"
)

func benchmarkSequenceTree() *b3.BehaviorTree {
	tree := b3.NewBehaviorTree()
	tree.Root = b3.NewSequence(
		b3.NewSucceeder(),
		b3.NewSucceeder(),
		b3.NewSucceeder(),
	)
	return tree
}

func benchmarkLoadTreeData() *b3.TreeData {
	return &b3.TreeData{
		Title:       "Benchmark Tree",
		Description: "Used for benchmarks",
		Root:        "node-1",
		Properties: map[string]any{
			"key": "value",
		},
		Nodes: map[string]b3.NodeData{
			"node-1": {
				Id:       "node-1",
				Name:     "Sequence",
				Title:    "Root",
				Children: []string{"node-2", "node-3", "node-4"},
			},
			"node-2": {
				Id:    "node-2",
				Name:  "Succeeder",
				Title: "A",
			},
			"node-3": {
				Id:    "node-3",
				Name:  "Succeeder",
				Title: "B",
			},
			"node-4": {
				Id:    "node-4",
				Name:  "Succeeder",
				Title: "C",
			},
		},
		CustomNodes: []b3.CustomNodeData{},
	}
}

func BenchmarkBehaviorTreeTickSequence(b *testing.B) {
	tree := benchmarkSequenceTree()
	blackboard := b3.NewBlackboard()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Tick(nil, blackboard)
	}
}

func BenchmarkBlackboardSetGet(b *testing.B) {
	blackboard := b3.NewBlackboard()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blackboard.Set("value", i, "tree-1", "node-1")
		_ = blackboard.Get("value", "tree-1", "node-1")
	}
}

func BenchmarkBehaviorTreeLoad(b *testing.B) {
	data := benchmarkLoadTreeData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := b3.NewBehaviorTree()
		if err := tree.Load(data, nil); err != nil {
			b.Fatalf("load failed: %v", err)
		}
	}
}

func BenchmarkBehaviorTreeDump(b *testing.B) {
	tree := benchmarkSequenceTree()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree.Dump()
	}
}

func BenchmarkCrosslangRunnerBatch(b *testing.B) {
	fixtures, err := loadFixturesRaw()
	if err != nil {
		b.Fatalf("load fixtures: %v", err)
	}
	repoRoot, _, err := repoPathsRaw()
	if err != nil {
		b.Fatalf("repo paths: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "node_modules", "babel-register")); err != nil {
		b.Skip("node_modules missing, run npm install --ignore-scripts first")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := runJSFixturesBatchRaw(repoRoot, fixtures); err != nil {
			b.Fatalf("run js fixture batch: %v", err)
		}
	}
}
