# GoLand Plugin

A companion GoLand plugin is available:

```text
https://github.com/D4r3E-1v1l/go-when-goland-plugin
```

The library works without the plugin.

The plugin is intended to improve the editing experience for `go-when` chains.

## Planned or supported checks

The plugin may provide inspections such as:

- missing terminal methods
- incomplete matcher branches
- invalid condition/action/terminal combinations
- duplicate terminals
- terminal not last
- numeric condition overlap
- unreachable numeric conditions
- enum exhaustive warnings

## Why a plugin helps

Some `go-when` rules are easier to check with editor tooling than with the Go type system alone.

Examples:

- a chain should end with `Else`, `ElseDo`, or `Exhaustive`
- numeric ranges should not accidentally overlap
- a broad range can make a later exact value unreachable
- an `Exhaustive()` enum chain may miss a newly added const case

## The plugin is optional

`go-when` is a normal Go library.

You can use it without installing the GoLand plugin.

The plugin only adds static analysis and editor feedback.
