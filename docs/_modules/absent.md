---
title: Absent
---

| Instruction | Action                      | Arguments                    |
| ----------- | --------------------------- | ---------------------------- |
| `absent`    | Delete files or directories | List of files or directories |

If a directory targeted by `absent` is not empty, it fails.

Example:

```yaml
- Custom personal directories tree:
    - absent:
        - Templates
        - Public
```
