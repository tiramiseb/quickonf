---
title: Move
---

| Instruction        | Action                                                                              | Arguments                      |
| ------------------ | ----------------------------------------------------------------------------------- | ------------------------------ |
| `move`             | Move files or directories                                                           | Map of sources to destinations |
| `force-move`       | Move files or directories, removing existing destinations                           | Map of sources to destinations |
| `migration-source` | Set the migration source path                                                       | The path                       |
| `migrate`          | Migrate files or directories from the previous home                                 | List of relative paths         |
| `migrate`          | Migrate files or directories from the previous home, removing existing destinations | List of relative paths         |

For all instructions in this module, if the source does not exist, this is a no-op.

# Migration

Let's say you have moved your previous home directory to `/home/previous` before reinstalling your system, in order to have a clean all-new environment. You still want to keep some directories or files as-is. That's what migration is about.

You first tell quickonf where your previous home is with `migration-source` and then tell it which files and directories you want to keep with `migrate`. For instance:

```yaml
- Initialization:
    - migration-source: /home/previous

- Keep my music:
    - migrate: Music

- Keep my GIMP configuration:
    - migrate: .config/GIMP
```

These instructions will move `/home/previous/Music` to `$HOME/Music` and `/home/previous/.config/GIMP` to `$HOME/.config/GIMP`. If the destination already exists, it is not overwritten.

Example:

```yaml
- My documents:
    - move:
        OLD/MyDocs: MyDocs
```
