---
title: Symlink
---

| Instruction | Action                 | Arguments                 |
| ----------- | ---------------------- | ------------------------- |
| `symlink`   | Create a symbolic link | Map from paths to targets |

Note the arguments are in not in the same order than for the `ln` command.

Example:

```yaml
- Make a symbolic link:
    - symlink:
        - tempdir: /tmp
```
