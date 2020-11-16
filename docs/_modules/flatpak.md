---
title: Flatpak
---

| Instruction      | Action                                           | Arguments                       |
| ---------------- | ------------------------------------------------ | ------------------------------- |
| `flatpak`        | Install flatpak packages from known repositories | List of packages names          |
| `flatpak-remote` | Add a flatpak remote repository                  | Map of repository names to URLs |

All instructions in this module need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

Example:

```yaml
- Anydesk:
    - flatpak:
        - anydesk
- Flathub:
    - flatpak-remote:
        flathub: https://dl.flathub.org/repo/flathub.flatpakrepo
```
