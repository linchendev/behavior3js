# Benchmarks

本文件记录 `behavior3go` 当前运行时优化完成后的基准基线，方便后续优化前后直接对照。

## 运行方式

```bash
cd ..
npm install --ignore-scripts
cd behavior3go
go test -run '^$' -bench . -benchmem ./...
```

## 当前基线

测试时间：`2026-04-19`  
测试环境：`linux/amd64`, `Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz`, `go1.25.9`

```text
BenchmarkBehaviorTreeTickSequence-12     2268369       530.3 ns/op       0 B/op       0 allocs/op
BenchmarkBlackboardSetGet-12            14467310       77.66 ns/op       7 B/op       0 allocs/op
BenchmarkBehaviorTreeLoad-12              373348      3026 ns/op      1993 B/op      30 allocs/op
BenchmarkBehaviorTreeDump-12             1227390       976.8 ns/op    1872 B/op      15 allocs/op
BenchmarkCrosslangRunnerBatch-12              4  254260351 ns/op    57484 B/op     801 allocs/op
```

## 解释

- `BehaviorTreeTickSequence`：衡量典型 `tick` 热路径的运行时与分配成本。
- `BlackboardSetGet`：衡量 tree/node 作用域下的黑板读写开销。
- `BehaviorTreeLoad`：衡量 JSON 树结构加载成本。
- `BehaviorTreeDump`：衡量运行时树导出为兼容结构的成本。
- `CrosslangRunnerBatch`：衡量 Go 发起、JS 批量执行共享夹具的总成本。

## 本轮变化

- `TickSequence` 从 `1151 ns/op, 480 B/op, 17 allocs/op` 降到 `530.3 ns/op, 0 B/op, 0 allocs/op`。
- `BehaviorTreeLoad` 从 `8315 ns/op, 2729 B/op, 58 allocs/op` 降到 `3026 ns/op, 1993 B/op, 30 allocs/op`。
- `BlackboardSetGet` 与 `Dump` 基本保持稳定，说明这轮收益主要来自 `Tick` 和 `Load` 热路径。

## 使用建议

- 后续做运行时优化时，优先比较 `TickSequence`、`BlackboardSetGet`、`Load`、`Dump` 的 `allocs/op` 是否下降。
- `CrosslangRunnerBatch` 更适合观察跨语言验证链路是否退化，不适合作为核心运行时性能指标。
- 如果硬件、Go 版本或 Node 版本变化较大，应重新生成一份基线，不直接横比旧数据。
