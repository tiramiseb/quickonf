---
title: Chrome / Chromium
---

You may enforce Chrome / Chromium extensions by creating system-wide policies.

## Prerequisite

```yaml
- Initialisation:
    - always: true
    - sudo-password: my-wonderful-password
```

## The recipe

For Chrome:

```yaml
- root-directory:
    - /etc/opt/chrome/policies/managed
- root-file:
    # cjpalhdlnbpafiamejdnhcphjbkeiagm = uBlock Origin
    # gphhapmejobijbbhgpjhcjognlahblep = GNOME Shell Integration
    /etc/opt/chrome/policies/managed/quickonf.json: |
      {
          "ExtensionInstallForcelist": [
              "cjpalhdlnbpafiamejdnhcphjbkeiagm",
              "gphhapmejobijbbhgpjhcjognlahblep"
          ]
      }
```

For Chromium:

```yaml
- root-directory:
    - /etc/chromium-browser/policies/managed
- root-file:
    # cjpalhdlnbpafiamejdnhcphjbkeiagm = uBlock Origin
    # gphhapmejobijbbhgpjhcjognlahblep = GNOME Shell Integration
    /etc/chromium-browser/policies/managed/quickonf.json: |
      {
          "ExtensionInstallForcelist": [
              "cjpalhdlnbpafiamejdnhcphjbkeiagm",
              "gphhapmejobijbbhgpjhcjognlahblep"
          ]
      }
```
