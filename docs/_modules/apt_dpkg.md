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

## dpkg

Install .deb packages. If any installation fails, it fails.

Example:

```yaml
- Install wonderfulsoft and magicapp:
    - dpkg:
        - /tmp/wonderfulsoft.deb
        - stuff/magicapp.deb
```

This instruction needs the sudo password to be set. See the [sudo]{% link _modules_/sudo.md %} module.

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

## apt

Install packages from apt repositories. If any installation fails, it fails.

Each package is installed with a different call to `apt`.

This instruction needs the sudo password to be set. See the [sudo]{% link _modules_/sudo.md %} module.

Example:

```yaml
- Install graphical tools:
    - apt:
        - gimp
        - inkscape
```

### apt-remove

Remove a deb package from the system.

Each package is removed with a different call to `apt`.

This instruction needs the sudo password to be set. See the [sudo]{% link _modules_/sudo.md %} module.

Example:

```yaml
- I do not want graphical toolss:
    - apt-remove:
        - gimp
        - inkscape
```

### apt-upgrade

Upgrade all installed packages.

This instruction needs the sudo password to be set. See the [sudo]{% link _modules_/sudo.md %} module.

Example:

```yaml
- Upgrade and clean the system:
    - apt-upgrade:
    - apt-autoremove-purge:
```

### apt-autoremove-purge

Clean unneeded packages from the system.

This instruction needs the sudo password to be set. See the [sudo]{% link _modules_/sudo.md %} module.

Example:

```yaml
- Upgrade and clean the system:
    - apt-upgrade:
    - apt-autoremove-purge:
```
