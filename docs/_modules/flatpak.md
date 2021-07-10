---
title: Flatpak
---

All instructions in this module need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

| Instruction      | Action                                           | Arguments                       | Dry run    |
| ---------------- | ------------------------------------------------ | ------------------------------- | ---------- |
| `flatpak`        | Install flatpak packages from known repositories | List of packages names          | No install |
| `flatpak-remote` | Add a flatpak remote repository                  | Map of repository names to URLs | No change  |

Example:

```yaml
- Anydesk:
    - flatpak:
        - anydesk
- Flathub:
    - flatpak-remote:
        flathub: https://dl.flathub.org/repo/flathub.flatpakrepo
```
