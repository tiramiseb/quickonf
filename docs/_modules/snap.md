---
title: Snap
---

| Instruction      | Action                                                   | Arguments                       |
| ---------------- | -------------------------------------------------------- | ------------------------------- |
| `snap`           | Install Snap packages                                    | List of packages                |
| `snap-classic`   | Install Snap packages in classic mode                    | List of packages                |
| `snap-dangerous` | Install Snap packages without verifying their signatures | List of packages                |
| `snap-edge`      | Install Snap packages from the edge channel              | List of packages                |
| `snap-refresh`   | Refresh Snap packages                                    | none                            |
| `snap-version`   | Get a Snap package version                               | Parameters as a map (see below) |

Classic mode disables security confinement.

## snap-version

Check a package version. If requested, save it to the store fur future use.

Parameters:

- `package` (mandatory): name of the package
- `store`: key name in the store (if not provided, the version is not stored)

Example:

```yaml
- Check version of Obsidian:
    - snap-version:
        package: obsidian
        store: obsidian-version
```
