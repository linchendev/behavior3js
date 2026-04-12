## ADDED Requirements

### Requirement: Go runtime loads existing behavior tree JSON without schema changes
The Go implementation SHALL accept the existing behavior tree JSON structure used by Behavior3JS, including `title`, `description`, `root`, `properties`, the `nodes` dictionary, per-node `id`, `name`, `title`, `description`, `properties`, and relationship fields such as `children` and `child`. Existing configuration files MUST be loadable without field renaming, preprocessing, or format conversion.

#### Scenario: Existing tree definition loads directly
- **WHEN** a user provides a JSON tree definition produced for the current Behavior3JS runtime
- **THEN** the Go runtime MUST parse and load that definition without requiring schema translation

#### Scenario: Unknown node name fails during load
- **WHEN** the JSON tree definition references a node `name` that is not registered as a built-in node or custom node
- **THEN** the Go runtime MUST fail the load operation with an explicit invalid-node error

### Requirement: Load flow preserves current construction semantics
The Go implementation SHALL preserve the current loading sequence semantics of `BehaviorTree.load`: build a node table from the `nodes` dictionary, resolve each node type by name, assign persisted node metadata and properties, connect composite children and decorator child references, and finally set the tree root from the `root` identifier.

#### Scenario: Composite children are connected after node creation
- **WHEN** a JSON tree contains a composite node with child identifiers
- **THEN** the Go runtime MUST create all referenced nodes before resolving the parent-child connections

#### Scenario: Decorator child is connected after node creation
- **WHEN** a JSON tree contains a decorator node with a `child` identifier
- **THEN** the Go runtime MUST resolve the child reference after the decorator node and child node both exist in the node table

### Requirement: Dump output remains semantically compatible
The Go implementation SHALL provide a dump or serialization path that emits a behavior tree document compatible with the current Behavior3JS structure, including tree metadata, `root`, `properties`, `nodes`, and `custom_nodes` for custom node definitions. A tree loaded from a compatible JSON document and then dumped by the Go runtime MUST preserve the same behavioral meaning.

#### Scenario: Dump preserves root and node metadata
- **WHEN** a loaded tree is dumped by the Go runtime
- **THEN** the resulting document MUST include the same root node identifier and the persisted node metadata needed to reconstruct the tree

#### Scenario: Dump includes custom node definitions
- **WHEN** a tree contains nodes that are not part of the built-in decorator, composite, or action registries
- **THEN** the dumped document MUST include corresponding `custom_nodes` entries describing those custom node types

### Requirement: Custom node registration remains name-based
The Go implementation SHALL support a custom node registration mechanism that preserves the current name-based loading contract. A custom node referenced by `name` in JSON MUST be resolved through a caller-provided registry equivalent in purpose to the JavaScript `names` argument.

#### Scenario: Custom node resolves from provided registry
- **WHEN** a JSON tree references a custom node name and the caller provides a registration entry for that name
- **THEN** the Go runtime MUST instantiate that custom node during load using the registered constructor
