---
title: File
---

| Instruction            | Action                                             | Arguments                   |
| ---------------------- | -------------------------------------------------- | --------------------------- |
| `file`                 | Create or replace files                            | Map of path to file content |
| `executable-file`      | Create or replace files with executable flag       | Map of path to file content |
| `restricted-file`      | Create or replace files only readable by the owner | Map of path to file content |
| `root-file`            | Create or replace files as root                    | Map of path to file content |
| `executable-root-file` | Create or replace executable files as root         | Map of path to file content |
| `restricted-root-file` | Create or replace files only readable by root      | Map of path to file content |

All these instructions are variants for the same action.

Instructions for root need the sudo password to be set. See the [sudo]{% link _modules/sudo.md %} module.

Example:

```yaml
- SSH key:
    - restricted-file:
        .ssh/id_rsa: |
          -----BEGIN RSA PRIVATE KEY-----
          blah blah blah
          -----END RSA PRIVATE KEY-----
    - file:
        .ssh/id_rsa.pub: |
          ssh-rsa AAAAblah blah blah
```
