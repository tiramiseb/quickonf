---
title: Sudo
---

| Instruction     | Action                              | Arguments             |
| --------------- | ----------------------------------- | --------------------- |
| `sudo-password` | Set the sudo passord for future use | The password for sudo |

# sudo-password

Set the sudo password for future use, whenever a root access is necessary. This password is then used in other modules. This instruction should be the first in the configuration file.

Example:

```yaml
- Initialisation:
    - sudo-password: my-wonderful-password
```
