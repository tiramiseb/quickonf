---
title: Stop if...
---

| Instruction     | Action                                   | Arguments                       | Dry run |
| --------------- | ---------------------------------------- | ------------------------------- | ------- |
| `stop-if-exist` | Check if given files or dirs exist       | List of paths                   | -       |
| `stop-if-older` | Compare versions and maybe stop the step | Parameters as a map (see below) | -       |

## Stop if exist

Stop the whole step (ignore next instructions) if all the listed files exist.

The step is not stopped if any of the listed files is missing.

## Stop if older

Stop the whole step (ignore next instructions) if the candidate is older than or identical to the current version.

Versions must comply with semantic versioning.

Parameters:

- `current` (mandatory): current version
- `candidate` (mandatory): candidate version

Example:

```yaml
- Wonderful App:
    - dpkg-version:
        package: wonderfulapp
        store: wonderfulapp-version
    - stop-if-older:
        current: <wonderfulapp-version>
        candidate: 2.1.4
```
