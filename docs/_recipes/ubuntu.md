---
title: Ubuntu
---

## Install ubuntu restricted extras

The Microsoft fonts package requests the user to agree to the license. To do so automatically, a debconf variable must be set.

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password

Restricted extras:
    - debconf-set:
        package: ttf-mscorefonts-installer
        variable: msttcorefonts/accepted-mscorefonts-eula
        value: "true"
    - apt:
        - ubuntu-restricted-extras
```
