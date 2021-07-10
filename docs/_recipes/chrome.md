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
- as:
    root:
      - directory:
          - /etc/opt/chrome/policies/managed
      - file:
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
- as:
    root:
      - directory:
          - /etc/chromium-browser/policies/managed
      - file:
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
