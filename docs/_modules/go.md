---
title: Go
---

| Instruction  | Action                        | Arguments                      |
| ------------ | ----------------------------- | ------------------------------ |
| `go-env`     | Set go environment parameters | Map of parameter name to value |
| `go-package` | Install Go packages           | List of packages               |

Example:

```yaml
- Go:
    - go-env:
        GOPATH: <home>/Code/go
    - go-package:
        - github.com/golang/mock/mockgen@latest
```
