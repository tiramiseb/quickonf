---
title: XDG
---

| Instruction        | Action                                 | Arguments                             |
| ------------------ | -------------------------------------- | ------------------------------------- |
| `xdg-autostart`    | Set autostart for applications         | List of app names or .desktop files   |
| `xdg-mime-default` | Set default applications for mimetypes | Map from mimetype to application name |
| `xdg-user-dir`     | Set a XDG user dir                     | Map from name to path                 |

For `xdg-autostart`, a ".desktop file" is the path to the `.desktop` file, especially if it is not in `/usr/share/applications`. If missing, the `.desktop` extension is automatically added.

For `xdg-autostart` and `xdg-mime-default`, application name is the `.desktop` file name without the extension.

Example:

```yaml
- Seafile autostart:
    - xdg-autostart:
        - seafile

- Inkscape:
    - xdg-mime-default:
        image/svg+xml: inkscape

- Change downloads directory:
    - xdg-user-dir:
        DOWNLOAD: Downloaded
```
