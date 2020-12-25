---
title: VS Code
---

| Instruction          | Action                       | Arguments          | Dry run    |
| -------------------- | ---------------------------- | ------------------ | ---------- |
| `vscode-extension`   | Install VS Code extensions   | List of extensions | No install |
| `vscodium-extension` | Install VS Codium extensions | List of extensions | No install |

Example:

```yaml
- Go extension for VSCode:
    - vscode-extension:
        - golang.go
```
