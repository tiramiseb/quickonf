---
title: Apt and dpkg
---

Instructions in this module, except `dpkg-version`, need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

| Instruction            | Action                                 | Arguments                       | Dry run                   |
| ---------------------- | -------------------------------------- | ------------------------------- | ------------------------- |
| `dpkg`                 | Install .deb packages                  | List of packages files          | No install                |
| `dpkg-dependencies`    | Install dependencies of .deb packages  | List of packages files          | Check list but no install |
| `dpkg-reconfigure`     | Reconfigure deb packages               | List of packages names          | No configuration          |
| `dpkg-version`         | Check a package version                | Parameters as a map (see below) | -                         |
| `debconf-set`          | Set a debconf variable                 | Parameters as a map (see below) | No change                 |
| `apt`                  | Install packages from APT repositories | List of packages names          | No install                |
| `apt-source`           | Add an apt source                      | Parameters as a map (see below) | No change                 |
| `apt-remove`           | Remove installed packages              | List of packages names          | No change                 |
| `apt-upgrade`          | Upgrade all installed packages         | none                            | No change                 |
| `apt-autoremove-purge` | Clean unneeded packages                | none                            | No change                 |
| `apt-clean-cache`      | Clean the APT cache                    | none                            | No change                 |

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

## debconf-set

Set a debconf variable.

Parameters:

- `package`: name of the package
- `variable`: name of the variable
- `value`: value to set

## apt-source

Add an apt source.

Parameters:

- `id`: some identifier for the source (used for sources.list.d and key filenames)
- `sources`: content of the sources.list.d file.
- `key` (optional): content of the GPG key used to sign packages

How to use it:

```yaml
- Install a wonderful app:
    - http-get:
        url: https://www.example.com/wonderful-app/key.gpg
        store: wonderful-app-key
    - apt-source:
        id: wonderful-app
        key: <wonderful-app-key>
        sources: |
          deb https://www.example.com/wonderful-app/deb <oscodename> main
    - apt:
        - wonderful-app
```
