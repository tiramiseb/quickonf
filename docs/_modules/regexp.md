---
title: RegExp
---

| Instruction      | Action                                          | Arguments                       | Dry run |
| ---------------- | ----------------------------------------------- | ------------------------------- | ------- |
| `regexp-replace` | Replace value according to a regular expression | Parameters as a map (see below) | -       |

## Regexp replace

Mandatory parameters for `regexp-replace` are:

- `from`: the value to transform
- `regexp`: the regular expression to apply on the value
- `repl` the replacement for the regexp

Optional parameter for `regexp-replace` is:

- `store`: key name in the store for the result (if not provided, the result is not stored)

See the [Zoom]({% link _recipes/zoom.md %}) recipe for an example.
