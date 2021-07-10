---
title: Packages repositories
---

## Prerequisite

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```

## Change the main apt repositories

```yaml
- APT sources:
    - as:
        root:
          - file:
            /etc/apt/sources.list: |
              # Sources from Quickonf
              #deb http://archive.ubuntu.com/ubuntu/   <oscodename>           main restricted universe multiverse
              #deb http://archive.ubuntu.com/ubuntu/   <oscodename>-updates   main restricted universe multiverse
              #deb http://archive.ubuntu.com/ubuntu/   <oscodename>-backports main restricted universe multiverse
              deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename>           main restricted universe multiverse
              deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename>-updates   main restricted universe multiverse
              deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename>-backports main restricted universe multiverse
              deb http://security.ubuntu.com/ubuntu/   <oscodename>-security  main restricted universe multiverse
              deb http://archive.canonical.com/ubuntu  <oscodename>           partner
          - apt-upgrade:
```

## Disable automatic apt updates

```yaml
- No automatic apt updates:
    - as:
        root:
          - file:
              /etc/apt/apt.conf.d/10periodic: |
                APT::Periodic::Update-Package-Lists "0";
                APT::Periodic::Download-Upgradeable-Packages "0";
                APT::Periodic::AutocleanInterval "0";
```

## Upgrade Snap packages

```yaml
- Snap packages:
    - snap-refresh:
```

## Install Flatpak

```yaml
- Flatpak:
    - apt:
        - flatpak
```
