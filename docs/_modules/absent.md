---
title: Absent
---

This module provides one instruction, `absent`, which makes sure files or directories are absent from the system. **It does not remove non-empty directories.**

## Example

```yaml
- Custom personal directories tree:
    - absent:
        - Templates
        - Public
```
