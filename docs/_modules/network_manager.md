---
title: Network Manager
---

| Instruction         | Action                       | Arguments             | Dry run   |
| ------------------- | ---------------------------- | --------------------- | --------- |
| `nm-wifi`           | Configure wifi WPA-PSK keys  | Map of SSID to PSK    | No change |
| `nm-import-openvpn` | Import OpenVPN configuration | List of `.ovpn` files | No change |

## Example

```yaml
Wifi:
  - nm-wifi:
      great-network: greatpassword
```
