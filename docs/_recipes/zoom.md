---
title: Zoom
---

## The problem

The Zoom .deb package is downloaded on the product's download page, not in a repository.

Its latest version may only be read in the download page. The name of the file to download is generic.

Moreover, the version does not comply with SemVer, because it contains 4 numbers.

## Prerequisite

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```

## The recipe

(For Ubuntu 64b)

```yaml
- Zoom:
    - parse-web-page:
        url: https://zoom.us/download?os=linux
        regexp: packageVersionX64 = 'Version (?P<mainver>[0-9]+\.[0-9]+)\.[0-9]+ \((?P<subver>[0-9]+)\.[0-9]+\)
        store-mainver: zoom-candidate-main
        store-subver: zoom-candidate-sub
    - dpkg-version:
        package: zoom
        store: zoom-current-full
    - regexp-replace:
        from: <zoom-current-full>
        regexp: ([0-9]+)\.([0-9]+)\.([0-9]+)\.[0-9]+
        replace: $1.$2.$3
        store: zoom-current
    - stop-if-older:
        current: <zoom-current>
        candidate: <zoom-candidate-main>.<zoom-candidate-sub>
    - download:
        https://zoom.us/client/latest/zoom_amd64.deb: /tmp/quickonf-zoom.deb
    - dpkg:
        - /tmp/quickonf-zoom.deb
    - remove:
        - /tmp/quickonf-zoom.deb
```
