---
title: Extract
---

| Instruction     | Action                   | Arguments                                 | Dry run   |
| --------------- | ------------------------ | ----------------------------------------- | --------- |
| `extract-tarxz` | Extract a tar.xz archive | Map of tar.xz files paths to destinations | No change |
| `extract-zip`   | Extract a zip archive    | Map of zip files paths to destinations    | No change |

Example:

```yaml
- Extract something:
    - extract-zip:
        /tmp/something.zip: /tmp/quickonf-something
```
