---
title: Apt and dpkg
---

| Instruction            | Action                                 | Arguments                       |
| ---------------------- | -------------------------------------- | ------------------------------- |
| `dpkg`                 | Install .deb packages                  | List of packages files          |
| `dpkg-version`         | Check a package version                | Parameters as a map (see below) |
| `apt`                  | Install packages from APT repositories | List of packages names          |
| `apt-remove`           | Remove installed packages              | List of packages names          |
| `apt-upgrade`          | Upgrade all installed packages         | none                            |
| `apt-autoremove-purge` | Clean unneeded packages                | none                            |
| `apt-clean-cache`      | Clean the APT cache                    | none                            |

Instructions in this module, except `dpkg-version`, need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

## dpkg-version

Check a package version. If requested, save it to the store fur future use.

Parameters:

- `package` (mandatory): name of the package
- `store`: key name in the store (if not provided, the version is not stored)

Example:

```yaml
- Check version of inkscape:
    - dpkg-version:
        package: inkscape
        store: inkscape-version
```
