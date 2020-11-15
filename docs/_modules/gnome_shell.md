---
title: GNOME Shell
---

| Instruction                           | Action                                                         | Arguments                       |
| ------------------------------------- | -------------------------------------------------------------- | ------------------------------- |
| `gnome-shell-extension`               | Enable GNOME Shell extensions                                  | List of extensions              |
| `local-gnome-shell-extension-version` | Check the version of a locally-installed GNOME Shell extension | Parameters as a map (see below) |
| `gnome-shell-restart`                 | Restart GNOME Shell                                            | none                            |

Parameters for version checking:

- `extension` (mandatory): name of the extension
- `store`: key name in the store (if not provided, the version is not stored)

Example:

```yaml
- GNOME Shell extensions:
    - gnome-shell-extension:
        - native-window-placement@gnome-shell-extensions.gcampax.github.com
        - places-menu@gnome-shell-extensions.gcampax.github.com
    - local-gnome-shell-extension-version:
        extension: places-menu@gnome-shell-extensions.gcampax.github.com
```
