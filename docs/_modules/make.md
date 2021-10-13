---
title: Make
---

| Instruction | Action                   | Arguments                       | Dry run   |
| ----------- | ------------------------ | ------------------------------- | --------- |
| `make`      | Execute the make command | Map of directories to arguments | No change |

For each element of the map, the make command is executed in the given directory, with the given arguments

Example:

```yaml
- Make:
    /tmp/foobar: build install
```

This example executes `make build install` in the `/tmp/foobar` directory
