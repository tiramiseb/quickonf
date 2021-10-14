---
title: Git
---

| Instruction       | Action                               | Arguments                           | Dry run           |
| ----------------- | ------------------------------------ | ----------------------------------- | ----------------- |
| `git-config`      | Set a git parameter                  | Map of parameter name to value      | No config         |
| `git-clone-pull`  | Clone or pull a repository           | Map of repository URL to path       | No clone nor pull |
| `git-commit-hash` | Store the hash of the current commit | Map of repository path to store key | -                 |

Example:

```yaml
- Git:
    - git-config:
        user.name: Some Name
        user.email: some.name@example.com
```
