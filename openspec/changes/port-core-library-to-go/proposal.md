## Why

当前仓库提供的是 JavaScript 版本的 Behavior3 运行时，核心行为树执行、黑板内存和 JSON 配置加载流程已经比较稳定，但运行时和构建链路都偏旧。将核心库改写为 Go，可以为服务端 AI、工具链和嵌入式场景提供更稳定的部署方式，同时保留现有配置资产和使用心智，降低迁移成本。

## What Changes

- 新增 Go 版本的核心运行时，并将全部实现放入 `behavior3go/` 目录，作为与现有 JavaScript 项目并列的独立项目，覆盖 `BehaviorTree`、`Blackboard`、`Tick`、`BaseNode` 及内置 `Action`、`Composite`、`Decorator` 节点。
- 保持与现有 JSON 配置格式兼容，包括树定义、节点定义、自定义属性以及 `load`/`dump` 的数据组织方式。
- 在 Go 实现中尽量保留现有命名、目录组织和执行语义，便于对照阅读与逐步迁移。
- 在不改变外部行为的前提下优化运行时实现，例如减少不必要的分配、明确错误路径和提升可测试性。
- 不移除现有 JavaScript 实现；本次变更以新增 `behavior3go/` 独立项目和兼容性约束为主，不引入破坏性配置变更。

## Capabilities

### New Capabilities
- `go-core-runtime-compatibility`: 提供与现有 Behavior3JS 核心运行时一致的 Go 实现，包括节点生命周期、状态传播和黑板作用域语义。
- `behavior-tree-json-compatibility`: 提供与现有树配置格式和加载流程一致的 Go 侧序列化与反序列化能力，确保现有配置文件可直接迁移使用。

### Modified Capabilities

None.

## Impact

- 新增 `behavior3go/` 独立项目目录、Go 模块定义和对应测试。
- 需要对现有核心行为做兼容性梳理，尤其是 `BehaviorTree.load`、`BehaviorTree.tick`、`Blackboard` 作用域和内置节点返回值语义。
- 文档需要补充 JS 与 Go 版本的关系、兼容范围和迁移边界。
