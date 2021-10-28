---
title: Systemd
---

All instructions in this module need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

| Instruction            | Action                        | Arguments                   | Dry run   |
| ---------------------- | ----------------------------- | --------------------------- | --------- |
| `systemd-enable`       | Enable systemd services       | List of services to enable  | No change |
| `systemd-disable`      | Disable systemd services      | List of services to disable | No change |
| `systemd-user-enable`  | Enable systemd user services  | List of services to enable  | No change |
| `systemd-user-disable` | Disable systemd user services | List of services to disable | No change |

Example:

```yaml
- Enable Sane socket service:
    - systemd-enable:
        - saned.socket
```
