---
title: File
---

| Instruction       | Action                                             | Arguments                   | Dry run   |
| ----------------- | -------------------------------------------------- | --------------------------- | --------- |
| `file`            | Create or replace files                            | Map of path to file content | No change |
| `executable-file` | Create or replace files with executable flag       | Map of path to file content | No change |
| `restricted-file` | Create or replace files only readable by the owner | Map of path to file content | No change |

All these instructions are variants for the same action.

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
