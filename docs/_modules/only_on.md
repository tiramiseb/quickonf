---
title: Only on...
---

| Instruction | Action                                                  | Arguments                     | Dry run |
| ----------- | ------------------------------------------------------- | ----------------------------- | ------- |
| `only-on`   | Run the step only if hostname matches the given pattern | Pattern for hostname matching | -       |

This instruction allows use the same `quickonf.yaml` file on multiple systems, with different configs on different systems.

Example:

```yaml
- Install Vim only on some machines:
	- only-on: dell*
	- apt:
		- vim
```
