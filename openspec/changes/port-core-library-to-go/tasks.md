## 1. Go 模块骨架

- [x] 1.1 在仓库中新增 `behavior3go/` 独立项目目录和 `go.mod`，建立与 `src/` 对照的 `core`、`actions`、`composites`、`decorators` 目录
- [x] 1.2 添加与 JS 常量、状态值和基础辅助函数对应的 Go 文件，明确命名映射和导出策略
- [x] 1.3 建立 `behavior3go/` 的最小可编译入口与包导出方案，确保其作为独立项目可被外部引用

## 2. 核心运行时迁移

- [x] 2.1 实现 Go 版 `Blackboard`，保持全局、tree、node 三层作用域语义一致
- [x] 2.2 实现 Go 版 `Tick` 和 `BaseNode` 生命周期包装逻辑，保持 `enter/open/tick/close/exit` 语义一致
- [x] 2.3 实现 Go 版 `BehaviorTree` 的 `Tick` 路径，保持 open-node 跟踪和关闭逻辑与 JS 基线一致
- [x] 2.4 实现 `Action`、`Composite`、`Decorator`、`Condition` 抽象基类，保持类型职责与名称对照

## 3. JSON 兼容与加载流程

- [x] 3.1 定义与现有树配置兼容的 Go 数据结构，覆盖树级字段、节点字段、关系字段和 `custom_nodes`
- [x] 3.2 实现 Go 版 `Load`，保持“先建节点表、再连边、最后设置根节点”的加载流程
- [x] 3.3 实现内置节点注册表和自定义节点注册机制，对齐 JS 的名称解析方式
- [x] 3.4 实现 Go 版 `Dump`，确保输出结构与现有 Behavior3JS JSON 语义兼容

## 4. 内置节点迁移

- [x] 4.1 在 `behavior3go/actions` 下迁移内置节点并保持类型名、文件位置和返回值语义一致
- [x] 4.2 在 `behavior3go/composites` 下迁移内置节点并保持子节点遍历顺序与 JS 一致
- [x] 4.3 在 `behavior3go/decorators` 下迁移内置节点并保持状态转换与黑板交互一致

## 5. 兼容性验证与非破坏性优化

- [x] 5.1 为 `behavior3go/` 的核心运行时、黑板作用域、JSON 加载与 dump 建立 Go 单元测试
- [x] 5.2 增加 JS/Go 共享夹具或对照测试，验证相同配置在两端得到一致结果
- [ ] 5.3 在不改变观察行为的前提下加入内部优化，并通过回归测试证明行为未变
- [x] 5.4 更新 README 或迁移文档，说明 `behavior3js` 与 `behavior3go` 的项目边界、兼容范围和使用方式
