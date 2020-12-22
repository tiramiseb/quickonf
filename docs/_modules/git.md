---
title: Git
---

| Instruction      | Action                     | Arguments                      |
| ---------------- | -------------------------- | ------------------------------ |
| `git-config`     | Set a git parameter        | Map of parameter name to value |
| `git-clone-pull` | Clone or pull a repository | Map of repository URL to path  |

Example:

```yaml
- Git:
    - git-config:
        user.name: Some Name
        user.email: some.name@example.com
```
