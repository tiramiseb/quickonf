---
title: Extract
---

| Instruction     | Action                   | Arguments                                 |
| --------------- | ------------------------ | ----------------------------------------- |
| `extract-tarxz` | Extract a tar.xz archive | Map of tar.xz files paths to destinations |
| `extract-zip`   | Extract a zip archive    | Map of zip files paths to destinations    |

Example:

```yaml
- Extract something:
    - unzip:
        /tmp/something.zip: /tmp/quickonf-something
```
