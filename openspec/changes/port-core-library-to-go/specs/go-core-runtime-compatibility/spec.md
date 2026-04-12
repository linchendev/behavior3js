## ADDED Requirements

### Requirement: Go runtime mirrors core type structure in an independent project
The Go implementation SHALL live entirely under the `behavior3go/` directory as a project independent from the existing JavaScript implementation. Within `behavior3go/`, it SHALL provide core runtime types and built-in node types that preserve the current Behavior3JS conceptual structure, including `BehaviorTree`, `Blackboard`, `Tick`, `BaseNode`, `Action`, `Composite`, `Decorator`, and built-in node names such as `Sequence`, `Priority`, `Repeater`, and `Wait`. The repository layout inside `behavior3go/` MUST mirror the current source grouping closely enough that a contributor can locate the equivalent Go implementation for a JS runtime file without a translation layer.

#### Scenario: Contributor locates equivalent core runtime file
- **WHEN** a contributor compares `src/core/BehaviorTree.js` and the Go implementation
- **THEN** the repository SHALL contain a corresponding Go runtime file under `behavior3go/core` with the `BehaviorTree` type and equivalent responsibility

#### Scenario: Contributor locates equivalent built-in node file
- **WHEN** a contributor compares `src/composites/Sequence.js` and the Go implementation
- **THEN** the repository SHALL contain a corresponding Go runtime file under `behavior3go/composites` with the `Sequence` type and equivalent responsibility

### Requirement: Node execution lifecycle remains behaviorally compatible
The Go implementation SHALL preserve the current node execution lifecycle semantics: `enter` executes on every traversal, `open` executes only when the node is not already open in blackboard state, `tick` performs the node logic, `close` executes whenever the returned status is not `RUNNING`, and `exit` executes after every traversal. Status propagation across parent and child nodes MUST remain equivalent to the current Behavior3JS runtime.

#### Scenario: Running node remains open across ticks
- **WHEN** a node returns `RUNNING` during a tree tick
- **THEN** the Go runtime MUST keep that node marked open in blackboard state and MUST NOT invoke `close` for that traversal

#### Scenario: Non-running node closes during traversal
- **WHEN** a node returns `SUCCESS`, `FAILURE`, or `ERROR`
- **THEN** the Go runtime MUST invoke the equivalent close path and persist the node as closed in blackboard state

### Requirement: Blackboard scope semantics remain compatible
The Go implementation SHALL preserve the current blackboard memory model with three scopes: global scope, tree scope, and node-within-tree scope. Reads and writes in one scope MUST NOT leak into another scope except where the current Behavior3JS semantics already allow shared access.

#### Scenario: Tree-scoped values remain isolated
- **WHEN** one value is written under `tree 1` scope and another under `tree 2` scope using the same key
- **THEN** reads from each tree scope MUST return only the value written for that tree

#### Scenario: Node-scoped values remain isolated within a tree
- **WHEN** two different node scopes write values under the same tree scope
- **THEN** a read using one node scope MUST NOT return the value written by the other node scope

### Requirement: Runtime optimizations do not change observable results
Any Go-specific optimization SHALL be non-breaking. For the same tree definition, target input, and blackboard state, the Go implementation MUST produce the same returned status, blackboard state transitions, and open-node tracking outcomes as the current JavaScript runtime.

#### Scenario: Optimized traversal preserves tree result
- **WHEN** a Go implementation uses internal optimizations such as slice preallocation or explicit error types
- **THEN** the externally observable tick result and blackboard updates MUST remain equivalent to the JavaScript baseline for the same scenario
