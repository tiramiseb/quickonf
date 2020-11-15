---
title: GNOME Shell
---

| Instruction             | Action                        | Arguments          |
| ----------------------- | ----------------------------- | ------------------ |
| `gnome-shell-extension` | Enable GNOME Shell extensions | List of extensions |
| `gnome-shell-restart`   | Restart GNOME Shell           | none               |

Example:

```yaml
- GNOME Shell extensions:
    - gnome-shell-extension:
        - native-window-placement@gnome-shell-extensions.gcampax.github.com
        - places-menu@gnome-shell-extensions.gcampax.github.com
```
