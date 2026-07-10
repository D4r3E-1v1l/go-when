# go-when

![Minimum Go Version](https://img.shields.io/badge/min%20go-1.22-blue)
[![LICENSE](https://img.shields.io/github/license/D4r3E-1v1l/go-when.svg)](https://github.com/D4r3E-1v1l/go-when/blob/main/LICENSE)

[English README](./README.md)

`go-when` 是一个面向 Go 的 typed decision-table helper。

它适合扁平、结果导向的决策映射，例如 value mapping、handler dispatch、fallible decision、error mapping、numeric range classification 和 mixed condition matching。

它不是为了替代 Go 的 `if` 或 `switch`，也不是完整 pattern matching 库。

## 安装

```bash
go get github.com/D4r3E-1v1l/go-when
```

## 要求

Go 1.22 或更高版本。

## 快速示例

```go
label := when.MatchAs[string](status).
	Case("pending").Then("Waiting for payment").
	Case("paid").Then("Payment received").
	Case("shipped").Then("On the way").
	Else("Unknown status")
```

## 为什么需要 go-when

Go 的 `if` 和 `switch` 很清晰，也很强大。

`go-when` 适合这种扁平、类型明确、结果导向的决策映射：

- value -> value
- value -> handler
- value -> `(value, error)`
- error -> value 或 handler
- numeric range -> category
- 混合 `Case` / `Range` / `When` / `Pattern` 条件

对于流程型、多步骤、控制流较重，或者分支里需要 `return`、`break`、`continue` 的代码，优先使用原生 `if` 或 `switch`。

## 混合条件匹配

`go-when` 的主要目标之一，是让不同类型的条件可以放在同一条 typed decision chain 里。

```go
grade := when.MatchAs[string](score).
	Case(100).Then("perfect").
	Range(when.Range(90, 100)).Then("excellent").
	Range(when.Range(60, 90)).Then("passed").
	When(isRetakeAllowed).Then("retake_allowed").
	Else("failed")
```

```go
func isRetakeAllowed(score int) bool {
	return score >= 50 && score < 60
}
```

匹配规则是 first match wins。更具体的条件应该放在更宽泛的条件前面。

## Value to Handler

`go-when` 可以把状态或命令映射成 handler。

```go
type Handler func()

handler := when.MatchAs[Handler](state).
	Case(OrderCreated).Then(requestPayment).
	Case(OrderPaid).Then(packOrder).
	Case(OrderPacked).Then(shipOrder).
	Case(OrderShipped).Then(sendTracking).
	Case(OrderCancelled).Then(cancelOrder).
	Exhaustive()

handler()
```

matcher 只表达决策：

```text
state -> handler
```

真正执行动作仍然显式：

```go
handler()
```

## Error Mapping

使用 `Err` 可以把 error 映射成 value 或 handler。

```go
handler := when.Err[ErrorHandler](err).
	Nil().Then(success).
	Is(ErrUserNotFound).Then(notFound).
	Is(ErrInvalidInput).Then(badRequest).
	Is(ErrNotAuthorized).Then(unauthorized).
	Else(internalError)

result := handler(err)
```

## Fallible Decisions

当决策需要返回 `(R, error)` 时，使用 `WithErr()`。

```go
handler, err := when.MatchAs[Handler](command).
	WithErr().
	Case(AddTodo).Then(addTodo).
	Case(DeleteTodo).Then(deleteTodo).
	Case(ExportTodo).ThenErr(nil, ErrExportDisabled).
	ElseErr(nil, ErrUnsupportedCommand)
```

这保持了 Go 原生错误风格：

```go
value, err := ...
```

`go-when` 不引入自定义 `Result[T, E]` 类型。

## 显式 Terminal

每条完整 matcher chain 都必须以 terminal method 结束。

使用 `Else` 表示接受 fallback：

```go
label := when.MatchAs[string](status).
	Case("paid").Then("Payment received").
	Else("Unknown status")
```

使用 `Exhaustive` 表示你认为所有合法 case 都已经覆盖：

```go
action := when.MatchAs[Action](state).
	Case(OrderCreated).Then(RequestPayment).
	Case(OrderPaid).Then(PackOrder).
	Case(OrderCancelled).Then(Noop).
	Exhaustive()
```

## Examples

可运行 examples 在 [`examples/`](./examples)。

```bash
go run ./examples/01_status_label
go run ./examples/02_order_state_handler
go run ./examples/03_score_grade
go run ./examples/04_error_to_result_handler
go run ./examples/05_command_dispatch_with_error
go run ./examples/06_signup_decision
go run ./examples/07_custom_pattern
```

## 文档

见 [`docs/`](./docs)。

## GoLand 插件

配套 GoLand 插件：

```text
https://github.com/D4r3E-1v1l/go-when-goland-plugin
```

插件会提供 matcher chain 结构检查和部分语义检查，例如 missing terminal、numeric overlap、unreachable numeric condition、enum exhaustive warning。

这个库不依赖插件也可以使用。

## License

MIT
