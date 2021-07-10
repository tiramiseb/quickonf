---
title: As
---

Instruction in this module needs the sudo password to be set. See the [sudo]({% link _modules/sudo.md %}) module.

| Instruction | Action                               | Arguments           | Dry run                   |
| ----------- | ------------------------------------ | ------------------- | ------------------------- |
| `as`        | Execute instructions as another user | Map of instructions | See executed instructions |

As long as the current user can use sudo to execute commands as another user, any instruction may be executed as another user with the `as` instruction.

In order to allow all instructions to be executed by any other user, a new instance of quickonf is executed as the given user, with only the given instruction in its configuration file, and its output is redirected to the current quickonf instance. It implies the `quickonf` file must be executable by the given user.

This works with `root` too, of course.

Example:

```yaml
- Create Bayek:
    - user-password:
        bayek: aya

- Hello Bayek:
    - as:
        bayek:
          - file:
              SALUT.txt: |
                Salut Bayek, tu as été automatiquement créé !
          - file:
              HELLO.txt: |
                Hello Bayek, you have been automatically created!
```
