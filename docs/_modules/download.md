---
title: Download
---

| Instruction | Action                           | Arguments            | Dry run     |
| ----------- | -------------------------------- | -------------------- | ----------- |
| `download`  | Download files from the Internet | Map of URLs to paths | No download |

Example:

```yaml
- Download CoreDNS 1.8.0:
    - download:
        https://github.com/coredns/coredns/releases/download/v1.8.0/coredns_1.8.0_linux_amd64.tgz: /tmp/coredns.tgz
```
