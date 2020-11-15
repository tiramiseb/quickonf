---
title: Systemd
---

| Instruction       | Action                   | Arguments                   |
| ----------------- | ------------------------ | --------------------------- |
| `systemd-enable`  | Enable systemd services  | List of services to enable  |
| `systemd-disable` | Disable systemd services | List of services to disable |

All instructions in this module need the sudo password to be set. See the [sudo]{% link _modules/sudo.md %} module.

Example:

```yaml
- Enable Sane socket service:
    - systemd-enable:
        - saned.socket
```
