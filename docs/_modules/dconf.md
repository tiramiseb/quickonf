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
        /org/gnome/desktop/interface/buttons-have-icons: true
```

Values may be:

- `true` or `false` for booleans
- any integer or float number
- an integer prefixed with `uint32:`
- any string
- a list of any other value (which is translated to an inline array)
