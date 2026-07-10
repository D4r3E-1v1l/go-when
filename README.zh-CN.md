# go-when

[English README](./README.md)

`go-when` 是一个面向 Go 的 typed decision matcher。

它主要用于 Go 服务端业务代码里的值映射、错误映射、区间分类和 action dispatch。

它不是完整 pattern matching 库，也不是为了替代 Go 的 `if` 或 `switch`。

## 安装

```bash
go get github.com/D4r3E-1v1l/go-when
```

## 要求

Go 1.18 或更高版本。

`go-when` 使用了 Go generics。

## 快速示例

```go
result := when.MatchAs[string](code).
    Case(200).Then("ok").
    Case(404).Then("not_found").
    Else("unknown")
```

## 为什么需要 go-when

Go 的 `if` 和 `switch` 很简单、很强大，但服务端业务代码里经常出现重复的 decision mapping：

* status code -> response
* enum/state -> action
* error -> HTTP/gRPC response
* numeric range -> level
* operation -> handler
* operation -> handler, error

`go-when` 专注于这些实际业务映射场景。

## MatchAs

`MatchAs` 用于 comparable value。

```go
action := when.MatchAs[Action](phase).
    Case(Pending).Then(Create).
    Case(Running).Then(Sync).
    Case(Failed).Then(Cleanup).
    Else(Noop)
```

## Range

`Range` 用于数值区间分类。

```go
level := when.MatchAs[string](score).
    Range(when.Range(0, 60)).Then("low").
    Range(when.Range(60, 90)).Then("medium").
    Range(when.From(90)).Then("high").
    Else("unknown")
```

## Error Mapping

`Err` 用于把 error 映射成业务值。

```go
resp := when.Err[HTTPResp](err).
    Is(ErrNotFound).Then(notFoundResp).
    Is(ErrPermissionDenied).Then(forbiddenResp).
    Contains("timeout").Then(timeoutResp).
    Else(internalResp)
```

## Fallible Matcher

如果 matcher 需要返回 `(R, error)`，使用 `WithErr()`。

```go
handler, err := when.MatchAs[Handler](op).
    WithErr().
    Case(Create, Update).Then(writeHandler).
    Case(Delete).Then(deleteHandler).
    ElseErr(nil, ErrUnsupportedOperation)
```

## 显式 Terminal

每条完整 matcher chain 都必须以 terminal method 结束：

* `Else(...)`
* `ElseDo(...)`
* `Exhaustive()`

如果你希望提供 fallback，使用 `Else` 或 `ElseDo`。

如果你认为所有合法 case 都已经覆盖，使用 `Exhaustive`。

## GoLand 插件

配套 GoLand 插件：

```text
https://github.com/D4r3E-1v1l/go-when-goland-plugin
```

插件支持检查：

* 缺少 terminal method
* 不完整 matcher branch
* matcher type 不支持的 condition/action/terminal
* 重复 terminal
* numeric condition overlap
* unreachable numeric condition
* enum exhaustive warning

## 文档

更多内容见 [`docs/`](./docs)。

## License

MIT
