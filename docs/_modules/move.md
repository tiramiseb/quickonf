---
title: Move
---

| Instruction  | Action                                                        | Arguments                      |
| ------------ | ------------------------------------------------------------- | ------------------------------ |
| `move`       | Move files or directories                                     | Map of sources to destinations |
| `force-move` | Move files or directories, removing destinations if necessary | Map of sources to destinations |

If the source does not exist, this is a no-op.

Example:

```yaml
- My documents:
    - move:
        OLD/MyDocs: MyDocs
```
