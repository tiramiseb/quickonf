---
title: Dconf
---

| Instruction | Action                                | Arguments             |
| ----------- | ------------------------------------- | --------------------- |
| `dconf`     | Set a parameter in the dconf database | Map of keys to values |

Example:

```yaml
- Display the icons on the buttons:
    - dconf:
        /org/gnome/desktop/interface/buttons-have-icons: "true"
```
