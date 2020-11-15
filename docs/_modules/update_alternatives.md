---
title: Update Alternatives
---

| Instruction           | Action                  | Arguments                  |
| --------------------- | ----------------------- | -------------------------- |
| `update-alternatives` | Change default commands | Map from command to target |

All instructions in this module need the sudo password to be set. See the [sudo]{% link _modules/sudo.md %} module.

Example:

```yaml
- Use vanilla GNOME Shell GDM theme:
    - update-alternatives:
        gdm3-theme.gresource: /usr/share/gnome-shell/gnome-shell-theme.gresource
```
