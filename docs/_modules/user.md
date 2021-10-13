---
title: User
---

All instructions in this module need the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

| Instruction     | Action                                             | Arguments                      | Dry run   |
| --------------- | -------------------------------------------------- | ------------------------------ | --------- |
| `user-password` | Make sure a user exists and has the given password | Map of user name to password   | No change |
| `user-in-group` | Make sure a user is in a group                     | Map of user name to group name | No change |

When creating a user, their shell is set to /bin/bash and their home is created.

Example:

```yaml
- John Doe must be able to user a modem:
    - user-password:
        johndoe: ilovejane
    - user-in-group:
        johndoe: dialout
```
