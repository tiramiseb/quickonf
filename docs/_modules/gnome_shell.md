---
title: GNOME Shell
---

| Instruction                           | Action                                                         | Arguments                       | Dry run    |
| ------------------------------------- | -------------------------------------------------------------- | ------------------------------- | ---------- |
| `install-gnome-shell-extension`       | Install GNOME Shell extensions from extensions.gnome.org       | List of extensions              | No install |
| `enable-gnome-shell-extension`        | Enable GNOME Shell extensions                                  | List of extensions              | No change  |
| `local-gnome-shell-extension-version` | Check the version of a locally-installed GNOME Shell extension | Parameters as a map (see below) | -          |
| `gnome-shell-restart`                 | Restart GNOME Shell                                            | none                            | No restart |

Extensions names are their UUIDs, like "places-menu@gnome-shell-extensions.gcampax.github.com".

Parameters for version checking:

- `extension` (mandatory): name of the extension
- `store`: key name in the store (if not provided, the version is not stored)

Example:

```yaml
- GNOME Shell extensions:
    - install-gnome-shell-extension:
        - gsconnect@andyholmes.github.io
    - enable-gnome-shell-extension:
        - gsconnect@andyholmes.github.io
        - native-window-placement@gnome-shell-extensions.gcampax.github.com
        - places-menu@gnome-shell-extensions.gcampax.github.com
    - local-gnome-shell-extension-version:
        extension: places-menu@gnome-shell-extensions.gcampax.github.com
```
