---
title: GitHub
---

| Instruction            | Action                                                | Arguments                       | Dry run |
| ---------------------- | ----------------------------------------------------- | ------------------------------- | ------- |
| `github-latest`        | Get the latest release for a GitHub repository        | Parameters as a map (see below) | -       |
| `github-latest-stable` | Get the latest stable release for a GitHub repository | Parameters as a map (see below) | -       |

Parameters:

- `repository` (mandatory): name of the repository
- `store`: key name in the store for the release name (version) (if not provided, the release is not stored)
- `pattern`: file pattern, if looking for a specific file
- `store-url`: key name in the store for the asset URL (if not provided, the asset URL is not stored)

Example:

```yaml
- Obsidian:
    - github-latest:
        repository: obsidianmd/obsidian-releases
        store: obsidian-release
        pattern: obsidian_*.snap
        store-url: obsidian-url
```
