---
title: Directory
---

| Instruction | Action             | Arguments           | Dry run     |
| ----------- | ------------------ | ------------------- | ----------- |
| `directory` | Create directories | List of directories | No creation |

If the target already exists and is not a directory, it fails.

Example:

```yaml
- Code directories:
    - directory:
        - Code/go
        - Code/python
```
