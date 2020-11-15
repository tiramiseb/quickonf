---
title: Pip
---

| Instruction | Action                            | Arguments               |
| ----------- | --------------------------------- | ----------------------- |
| `pip`       | Install Python packages using pip | List of Python packages |

All instructions in this module need the sudo password to be set. See the [sudo]{% link _modules/sudo.md %} module.

Example:

```yaml
- Sigal:
    - pip:
        - sigal
```
