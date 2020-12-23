---
title: TeamViewer
---

## The problem

The TeamViewer .deb package is downloaded on the product's download page, not in a repository.

Moreover, its latest version may only be read in the download page. The name of the file to download is generic.

## Prerequisite

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```

## The recipe

(For Ubuntu 64b)

```yaml
- TeamViewer:
    - parse-web-page:
        url: https://www.teamviewer.com/download/linux/
        regexp: \*\.deb package (?P<ver>[0-9]+\.[0-9]+\.[0-9]+)
        store-ver: teamviewer-candidate
    - dpkg-version:
        package: teamviewer
        store: teamviewer-current
    - stop-if-older:
        current: <teamviewer-current>
        candidate: <teamviewer-candidate>
    - download:
        https://download.teamviewer.com/download/linux/teamviewer_amd64.deb: /tmp/quickonf-teamviewer.deb
    - dpkg:
        - /tmp/quickonf-teamviewer.deb
    - remove:
        - /tmp/quickonf-teamviewer.deb
```
