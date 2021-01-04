---
title: Cura by Dagoma
---

## The problem

Dagoma has forked the Cura software, to optimise it for their printers.

This software is especially tricky to get and update. If you want to update it manually, there are multiple steps:

- first, check the version in the software
- then, go to the download page and click on multiple links (implying multiple requests to an API)
- then, download the returned link, you receive a .zip file
- then, unzip the .zip file, it contains a README and a .deb file (for Debian/Ubuntu)
- then, try to install the .deb file
- then, install the missing dependencies
- then, delete the .zip file and the extracted content

## Automatization

Checking the Dagoma website is a matter downloading an URL.

The multiple links clicks may be replaced by a single JSON REST API request, which returns the download link in a JSON response, that must be parsed.

Then, downloading and extracting a .zip, installing the dependencies, installing the .deb, etc, is easy to automatize.

## Prerequisite

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```

## The recipe

(for Ubuntu 64b)

```yaml
- Cura by Dagoma:
    - http-get:                                                                 
        url: https://dist.dagoma3d.com/version/CuraByDagoma                     
        store: curabydago-candidate 
    - dpkg-version:
        package: curabydagoma
        store: curabydago-current
    - xdg-mime-default:
        model/stl: curabydago
    - stop-if-older:
        current: <curabydago-current>
        candidate: <curabydago-candidate>
    - json-build:
        - store=curabydago-api-req
        - product=CuraByDagoma
        - os=Linux
        - arch=x64
        - package_type=debian
    - http-post:
        url: https://dist.dagoma3d.com/api
        payload: <curabydago-api-req>
        store: curabydago-api-resp
    - json-get:
        from: <curabydago-api-resp>
        key: download-links.zip
        store: curabydago-relative-url
    - download:
        https://dist.dagoma3d.com<curabydago-relative-url>: /tmp/quickonf-curabydago.zip
    - extract-zip:
        /tmp/quickonf-curabydago.zip: /tmp/quickonf-curabydago
    - dpkg-dependencies:
        - /tmp/quickonf-curabydago/CuraByDagoma_amd64.deb
    - dpkg:
        - /tmp/quickonf-curabydago/CuraByDagoma_amd64.deb
    - remove:
        - /tmp/quickonf-curabydago/CuraByDagoma_amd64.deb
        - /tmp/quickonf-curabydago/README.md
        - /tmp/quickonf-curabydago
        - /tmp/quickonf-curabydago.zip
```
