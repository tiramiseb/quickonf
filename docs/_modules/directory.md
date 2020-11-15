---
title: Directory
---

| Instruction | Action             | Arguments           |
| ----------- | ------------------ | ------------------- |
| `directory` | Create directories | List of directories |

If the target already exists and is not a directory, it fails.

Example:

```yaml
- Code directories:
    - directory:
        - Code/go
        - Code/python
```
