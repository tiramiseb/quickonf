---
title: Symlink
---

| Instruction | Action                 | Arguments                 | Dry run     |
| ----------- | ---------------------- | ------------------------- | ----------- |
| `symlink`   | Create a symbolic link | Map from paths to targets | No creation |

Note the arguments are in not in the same order than for the `ln` command.

Example:

```yaml
- Make a symbolic link:
    - symlink:
        - tempdir: /tmp
```
