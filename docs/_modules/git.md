---
title: Git
---

| Instruction  | Action              | Arguments                      |
| ------------ | ------------------- | ------------------------------ |
| `git-config` | Set a git parameter | Map of parameter name to value |

Example:

```yaml
- Git:
    - git-config:
        user.name: Some Name
        user.email: some.name@example.com
```
