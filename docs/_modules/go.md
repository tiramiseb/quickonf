---
title: Go
---

| Instruction  | Action                        | Arguments                      | Dry run    |
| ------------ | ----------------------------- | ------------------------------ | ---------- |
| `go-env`     | Set go environment parameters | Map of parameter name to value | No change  |
| `go-package` | Install Go packages           | List of packages               | No install |

Example:

```yaml
- Go:
    - go-env:
        GOPATH: <home>/Code/go
    - go-package:
        - github.com/golang/mock/mockgen@latest
```
