---
title: Remove
---

| Instruction | Action                      | Arguments                          | Dry run    |
| ----------- | --------------------------- | ---------------------------------- | ---------- |
| `remove`    | Remove files or directories | List of files or directories paths | No removal |

If the target does not exist, this is a no-op.

If a directory targeted by `remove` is not empty, it fails.

Example:

```yaml
- Remove old graphic software configuration:
    - remove:
        - .config/inkscape
        - .config/GIMP
```
