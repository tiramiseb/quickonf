---
title: Only on...
---

| Instruction | Action                                              | Arguments                     | Dry run |
| ----------- | --------------------------------------------------- | ----------------------------- | ------- |
| `only-on`   | Continue only if hostname matches the given pattern | Pattern for hostname matching | -       |

This instruction allows use the same `quickonf.yaml` file on multiple systems, with different configs on different systems: if the hostname does not match the pattern, the step is stopped and further instructions are not executed.

Example:

```yaml
- Install Vim only on some machines:
	- only-on: dell*
	- apt:
		- vim
```
