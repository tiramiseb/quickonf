---
title: If...
---

| Instruction          | Action                                                             | Arguments                       | Dry run |
| -------------------- | ------------------------------------------------------------------ | ------------------------------- | ------- |
| `stop-if-exist`      | Check if given files or dirs exist and maybe stop the step         | List of paths                   | -       |
| `stop-if-older`      | Compare versions and maybe stop the step                           | Parameters as a map (see below) | -       |
| `stop-if-equal`      | Compare values and maybe stop the step                             | List if values                  | -       |
| `skip-next-if-exist` | Check if given files or dirs exist and maybe skip next instruction | List of paths                   | -       |
| `skip-next-if-older` | Compare versions and maybe skip next instruction                   | Parameters as a map (see below) | -       |
| `skip-next-if-equal` | Compare values and maybe skip next instruction                     | List if values                  | -       |

## Stop if exist

Stop the whole step (ignore all fulrther instructions) if all the listed files exist.

The step is not stopped if any of the listed files is missing.

## Stop if older

Stop the whole step (ignore all fulrther instructions) if the candidate is older than or identical to the current version.

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

## Stop if equal

Stop the whole step (ignore all fulrther instructions) if all the given values are equal.

## Skip next

Instead of stopping the whole step, skip only the next instruction.
