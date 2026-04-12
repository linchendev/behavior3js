package behavior3go_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	b3 "github.com/behavior3/behavior3go"
	"github.com/behavior3/behavior3go/core"
)

type fixtureBlackboardState struct {
	Base  map[string]any            `json:"base"`
	Tree  map[string]any            `json:"tree"`
	Nodes map[string]map[string]any `json:"nodes"`
}

type fixtureTick struct {
	Now    *int64 `json:"now"`
	Target any    `json:"target"`
}

type fixtureCompare struct {
	Status     bool `json:"status"`
	TreeMemory struct {
		OpenNodes bool `json:"openNodes"`
		NodeCount bool `json:"nodeCount"`
	} `json:"treeMemory"`
	Dump bool `json:"dump"`
}

type crosslangFixture struct {
	Name       string                 `json:"name"`
	Tree       b3.TreeData            `json:"tree"`
	Blackboard fixtureBlackboardState `json:"blackboard"`
	Ticks      []fixtureTick          `json:"ticks"`
	Compare    fixtureCompare         `json:"compare"`
}

type crosslangResult struct {
	Statuses   []b3.Status `json:"statuses"`
	TreeMemory struct {
		OpenNodes []string `json:"openNodes"`
		NodeCount int      `json:"nodeCount"`
	} `json:"treeMemory"`
	Dump map[string]any `json:"dump"`
}

type expectedFixtureResult struct {
	Statuses  []b3.Status
	OpenNodes []string
	NodeCount int
}

var expectedResults = map[string]expectedFixtureResult{
	"sequence-success": {
		Statuses:  []b3.Status{b3.SUCCESS},
		OpenNodes: []string{},
		NodeCount: 3,
	},
	"priority-running": {
		Statuses:  []b3.Status{b3.RUNNING},
		OpenNodes: []string{"node-1", "node-3"},
		NodeCount: 3,
	},
	"mem-sequence-resume": {
		Statuses:  []b3.Status{b3.RUNNING, b3.SUCCESS},
		OpenNodes: []string{},
		NodeCount: 3,
	},
	"mem-priority-resume": {
		Statuses:  []b3.Status{b3.RUNNING, b3.SUCCESS},
		OpenNodes: []string{},
		NodeCount: 2,
	},
	"wait-time-window": {
		Statuses:  []b3.Status{b3.RUNNING, b3.RUNNING, b3.SUCCESS},
		OpenNodes: []string{},
		NodeCount: 1,
	},
	"max-time-timeout": {
		Statuses:  []b3.Status{b3.RUNNING, b3.FAILURE},
		OpenNodes: []string{"node-1"},
		NodeCount: 2,
	},
	"dump-roundtrip": {
		Statuses:  []b3.Status{b3.RUNNING},
		OpenNodes: []string{"node-1", "node-2", "node-4"},
		NodeCount: 3,
	},
}

func repoPaths(t *testing.T) (string, string) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve test file path")
	}

	behavior3goDir := filepath.Dir(filename)
	repoRoot := filepath.Dir(behavior3goDir)
	return repoRoot, behavior3goDir
}

func loadFixtures(t *testing.T) []struct {
	path    string
	fixture crosslangFixture
} {
	t.Helper()

	_, behavior3goDir := repoPaths(t)
	pattern := filepath.Join(behavior3goDir, "testdata", "crosslang", "*.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob fixtures: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("expected cross-language fixtures")
	}

	fixtures := make([]struct {
		path    string
		fixture crosslangFixture
	}, 0, len(paths))

	for _, path := range paths {
		payload, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read fixture %s: %v", path, err)
		}

		var fixture crosslangFixture
		if err := json.Unmarshal(payload, &fixture); err != nil {
			t.Fatalf("decode fixture %s: %v", path, err)
		}

		fixtures = append(fixtures, struct {
			path    string
			fixture crosslangFixture
		}{
			path:    path,
			fixture: fixture,
		})
	}

	return fixtures
}

func applyBlackboardState(blackboard *b3.Blackboard, treeId string, state fixtureBlackboardState) {
	for key, value := range state.Base {
		blackboard.Set(key, value)
	}

	for key, value := range state.Tree {
		blackboard.Set(key, value, treeId)
	}

	for nodeId, nodeState := range state.Nodes {
		for key, value := range nodeState {
			blackboard.Set(key, value, treeId, nodeId)
		}
	}
}

func normalizeDump(t *testing.T, value any) map[string]any {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal dump: %v", err)
	}

	var normalized map[string]any
	if err := json.Unmarshal(payload, &normalized); err != nil {
		t.Fatalf("unmarshal dump: %v", err)
	}

	return normalized
}

func nodeIDs(value any) []string {
	if value == nil {
		return []string{}
	}

	nodes, ok := value.([]b3.Node)
	if ok {
		ids := make([]string, 0, len(nodes))
		for _, node := range nodes {
			ids = append(ids, node.GetBaseNode().Id)
		}
		return ids
	}

	return []string{}
}

func nodeCount(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}

func runGoFixture(t *testing.T, fixture crosslangFixture) crosslangResult {
	t.Helper()

	originalNow := core.Now
	defer func() {
		core.Now = originalNow
	}()

	tree := b3.NewBehaviorTree()
	if err := tree.Load(&fixture.Tree, nil); err != nil {
		t.Fatalf("load go fixture %s: %v", fixture.Name, err)
	}

	blackboard := b3.NewBlackboard()
	applyBlackboardState(blackboard, tree.Id, fixture.Blackboard)

	result := crosslangResult{
		Statuses: make([]b3.Status, 0, len(fixture.Ticks)),
	}

	for _, step := range fixture.Ticks {
		if step.Now != nil {
			now := *step.Now
			core.Now = func() int64 {
				return now
			}
		} else {
			core.Now = originalNow
		}

		result.Statuses = append(result.Statuses, tree.Tick(step.Target, blackboard))
	}

	result.TreeMemory.OpenNodes = nodeIDs(blackboard.Get("openNodes", tree.Id))
	result.TreeMemory.NodeCount = nodeCount(blackboard.Get("nodeCount", tree.Id))
	result.Dump = normalizeDump(t, tree.Dump())

	return result
}

func runJSFixture(t *testing.T, fixturePath string) crosslangResult {
	t.Helper()

	repoRoot, _ := repoPaths(t)
	if _, err := os.Stat(filepath.Join(repoRoot, "node_modules", "babel-register")); err != nil {
		t.Fatalf("cross-language tests require npm install at repo root: %v", err)
	}

	command := exec.Command("node", "test/crosslang/runner.js", fixturePath)
	command.Dir = repoRoot

	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run js fixture %s: %v\n%s", fixturePath, err, string(output))
	}

	var result crosslangResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode js result %s: %v\n%s", fixturePath, err, string(output))
	}

	return result
}

func TestCrosslangFixturesGoOnly(t *testing.T) {
	for _, entry := range loadFixtures(t) {
		entry := entry
		t.Run(entry.fixture.Name, func(t *testing.T) {
			result := runGoFixture(t, entry.fixture)
			expected, ok := expectedResults[entry.fixture.Name]
			if !ok {
				t.Fatalf("missing expected result for fixture %s", entry.fixture.Name)
			}

			if !reflect.DeepEqual(result.Statuses, expected.Statuses) {
				t.Fatalf("statuses mismatch: got %v want %v", result.Statuses, expected.Statuses)
			}
			if !reflect.DeepEqual(result.TreeMemory.OpenNodes, expected.OpenNodes) {
				t.Fatalf("open nodes mismatch: got %v want %v", result.TreeMemory.OpenNodes, expected.OpenNodes)
			}
			if result.TreeMemory.NodeCount != expected.NodeCount {
				t.Fatalf("node count mismatch: got %d want %d", result.TreeMemory.NodeCount, expected.NodeCount)
			}
			if entry.fixture.Compare.Dump && len(result.Dump) == 0 {
				t.Fatal("dump should not be empty")
			}
		})
	}
}

func TestCrosslangFixturesMatchJS(t *testing.T) {
	for _, entry := range loadFixtures(t) {
		entry := entry
		t.Run(entry.fixture.Name, func(t *testing.T) {
			goResult := runGoFixture(t, entry.fixture)
			jsResult := runJSFixture(t, entry.path)

			if entry.fixture.Compare.Status && !reflect.DeepEqual(goResult.Statuses, jsResult.Statuses) {
				t.Fatalf("status mismatch: got %v want %v", goResult.Statuses, jsResult.Statuses)
			}
			if entry.fixture.Compare.TreeMemory.OpenNodes && !reflect.DeepEqual(goResult.TreeMemory.OpenNodes, jsResult.TreeMemory.OpenNodes) {
				t.Fatalf("openNodes mismatch: got %v want %v", goResult.TreeMemory.OpenNodes, jsResult.TreeMemory.OpenNodes)
			}
			if entry.fixture.Compare.TreeMemory.NodeCount && goResult.TreeMemory.NodeCount != jsResult.TreeMemory.NodeCount {
				t.Fatalf("nodeCount mismatch: got %d want %d", goResult.TreeMemory.NodeCount, jsResult.TreeMemory.NodeCount)
			}
			if entry.fixture.Compare.Dump && !reflect.DeepEqual(goResult.Dump, jsResult.Dump) {
				goPayload, _ := json.MarshalIndent(goResult.Dump, "", "  ")
				jsPayload, _ := json.MarshalIndent(jsResult.Dump, "", "  ")
				t.Fatalf("dump mismatch for %s\ngo:\n%s\njs:\n%s", entry.fixture.Name, string(goPayload), string(jsPayload))
			}
		})
	}
}

func TestCrosslangFixtureCatalog(t *testing.T) {
	fixtures := loadFixtures(t)
	if len(fixtures) < 6 {
		t.Fatalf("expected at least 6 fixtures, got %d", len(fixtures))
	}

	names := map[string]bool{}
	for _, entry := range fixtures {
		if names[entry.fixture.Name] {
			t.Fatalf("duplicate fixture %s", entry.fixture.Name)
		}
		names[entry.fixture.Name] = true
		if _, ok := expectedResults[entry.fixture.Name]; !ok {
			t.Fatalf("fixture %s missing expected Go-only result", entry.fixture.Name)
		}
	}

	for name := range expectedResults {
		if !names[name] {
			t.Fatalf("expected result %s has no fixture", name)
		}
	}

	t.Logf("loaded %d cross-language fixtures", len(fixtures))
}
