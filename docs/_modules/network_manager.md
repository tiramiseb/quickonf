---
title: Network Manager
---

| Instruction         | Action                       | Arguments             |
| ------------------- | ---------------------------- | --------------------- |
| `nm-import-openvpn` | Import OpenVPN configuration | List of `.ovpn` files |
| `nm-wifi`           | Configure wifi WPA-PSK keys  | Map of SSID to PSK    |

## Example

```yaml
Wifi:
  - nm-wifi:
      great-network: greatpassword
```
