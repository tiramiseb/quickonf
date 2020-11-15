---
title: Unzip
---

| Instruction | Action                | Arguments                              |
| ----------- | --------------------- | -------------------------------------- |
| `unzip`     | Extract a zip archive | Map of zip files paths to destinations |

Example:

```yaml
- Extract something:
    - unzip:
        /tmp/something.zip: /tmp/quickonf-something
```
