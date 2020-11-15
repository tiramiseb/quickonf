---
title: XDG
---

| Instruction        | Action                                 | Arguments                             |
| ------------------ | -------------------------------------- | ------------------------------------- |
| `xdg-mime-default` | Set default applications for mimetypes | Map from mimetype to application name |
| `xdg-user-dir`     | Set a XDG user dir                     | Map from name to path                 |

Application name is the `.desktop` file name without the extension.

Example:

```yaml
- Inkscape:
    - xdg-mime-default:
        image/svg+xml: inkscape

- Change downloads directory:
    - xdg-user-dir:
        DOWNLOAD: Downloaded
```
