---
title: Gtk Bookmarks
---

| Instruction    | Action                | Arguments                     | Dry run   |
| -------------- | --------------------- | ----------------------------- | --------- |
| `gtk-bookmark` | Set the Gtk bookmarks | List of bookmarks (see below) | No change |

A bookmark may be in one of the following forms:

- `alias=path`: the bookmark named `alias` points to the given path
- `path`: the bookmark name is the path basename

Example:

```yaml
- Gtk bookmarks:
    - gtk-bookmark:
        - temp=/tmp
        - Pictures/Phone
```

This example results in:

- a bookmark named "temp" pointing to /tmp
- a bookmark named "Phone" pointing to Pictures/Phone
