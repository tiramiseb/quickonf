---
title: Adduser
---

| Instruction         | Action              | Arguments                      |
| ------------------- | ------------------- | ------------------------------ |
| `add-user-to-group` | Add users to groups | Map of user name to group name |

All instructions in this module need the sudo password to be set. See the [sudo]{% link _modules/sudo.md %} module.

Example:

```yaml
- Allow me to dial out:
    - add-user-to-group:
        sebastien: dialout
```
