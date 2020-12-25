---
title: Dconf
---

| Instruction | Action                                | Arguments             | Dry run   |
| ----------- | ------------------------------------- | --------------------- | --------- |
| `dconf`     | Set a parameter in the dconf database | Map of keys to values | No change |

Example:

```yaml
- Display the icons on the buttons:
    - dconf:
        /org/gnome/desktop/interface/buttons-have-icons: "true"
```
