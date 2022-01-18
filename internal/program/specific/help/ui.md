The interactive user interface is composed of two panes. The left pane lists all checks done (according to your configuration file), group by group, the right pane lists actions that would be applied. Each group may be expanded to see its details.

# Colors

Colored background on the whole line represent a group title. Instructions, on the other side, are prefixed with the instruction name (with a colored background), followed by details.

## Checks

Groups titles:

- _green_: everything is already applied
- _blue_: there is something to change
- _yellow_: currently checking
- _red_: a check failed, giving up
- _grey_: waiting to be checked

Instructions prefixes:

- _green_: instruction is already applied
- _blue_: something may be applied
- _red_: check failed

## Applies

Groups titles:

- _grey_: changes are waiting to be applied
- _yellow_: changes are being applied
- _green_: changes have been applied
- _red_: applying changes failed

Instructions prefixes:

- _blue_: something may be applied
- _green_: instruction has been applied
- _red_: applying change failed (next instructions are aborted)

# Keyboard usage

| key            | action                                             |
| -------------- | -------------------------------------------------- |
| left           | switch to the checks pane                          |
| right          | switch to the applies pane                         |
| up             | select the previous group                          |
| down           | select the next group                              |
| space, t       | toggle (expand) the current group                  |
| enter, x       | recheck or apply the current group                 |
| f              | filter the groups (hide/show the succeeded groups) |
| h              | display the help                                   |
| q, esc, ctrl-c | quit _Quickonf_                                    |

# Mouse usage

Clicking on a group selects it. If the group is already selected, it is toggled.

The buttons in the top-right corner are the following:

| button  | action                                             |
| ------- | -------------------------------------------------- |
| Recheck | re-execute the current group checks                |
| Apply   | apply the current group actions                    |
| Toggle  | toggle the current group                           |
| Filter  | filter the groups (hide/show the succeeded groups) |
| Help    | show or hide the help section                      |
| Quit    | quit _Quickonf_                                    |
