---
title: Sudo
---

| Instruction     | Action                              | Arguments             |
| --------------- | ----------------------------------- | --------------------- |
| `sudo-password` | Set the sudo passord for future use | The password for sudo |

## Sudo password

Set the sudo password for future use, whenever access to another user is necessary. This password is then used in other modules. This instruction should be the first in the configuration file.

Example:

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```
