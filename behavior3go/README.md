# behavior3go

`behavior3go` 是与现有 `behavior3js` 并列的独立 Go 项目，实现了 Behavior3 核心运行时，并尽量保持与原 JavaScript 版本一致的命名、目录结构和 JSON 配置兼容性。

## 目录

- `core/`: `BehaviorTree`、`Blackboard`、`Tick`、`BaseNode` 及抽象节点类型
- `actions/`: `Error`、`Failer`、`Runner`、`Succeeder`、`Wait`
- `composites/`: `Sequence`、`Priority`、`MemSequence`、`MemPriority`
- `decorators/`: `Inverter`、`Limiter`、`MaxTime`、`Repeater`、`RepeatUntilFailure`、`RepeatUntilSuccess`

## 开发

```bash
cd ..
npm install --ignore-scripts
cd behavior3go
go test ./...
go test -run '^$' -bench . -benchmem
```

默认入口包为模块根包 `github.com/behavior3/behavior3go`，它会导出核心类型和内置节点，并自动注册默认节点用于 `BehaviorTree.Load(...)`。

跨语言兼容测试由 Go 测试发起，并通过仓库根目录下的 `test/crosslang/runner.js` 调用原 `behavior3js` 实现。共享夹具位于 `behavior3go/testdata/crosslang/`。

当前性能基线记录在 [BENCHMARKS.md](/home/chenlin/works/github/behavior3js/behavior3go/BENCHMARKS.md)。
