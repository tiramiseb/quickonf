---
title: XDG
---

| Instruction    | Action             | Arguments             |
| -------------- | ------------------ | --------------------- |
| `xdg-user-dir` | Set a XDG user dir | Map from name to path |

Example:

```yaml
- Change downloads directory:
    - xdg-user-dir:
        DOWNLOAD: Downloaded
```
