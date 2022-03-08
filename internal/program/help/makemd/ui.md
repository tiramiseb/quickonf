# User interface

The interactive user interface is composed of two panes.

The left pane lists all groups (according to your configuration file), while the right pane shows details on the selected group.

# Filters and details

By default, groups are filtered so that only groups in need of a change are displayed. You can toggle this filter with the "F" key or the "Filter checks" button in the title bar.

If you need more details about actions that would be executed, you can toggle them with the "D" key or the "More details" button in the title bar.

# Colors

## Checks

- _grey_: waiting to be checked or being checked
- _green_: everything is already applied
- _red_: a check failed, giving up
- _blue_: there is something to change
- _yellow_: currently applying

## Details

- _green_: nothing needs to be applied
- _blue_: something must be applied
- _yellow_: currently applying
- _red_: applying change failed (next instructions are aborted)

# Keyboard usage

| key            | action                                             |
| -------------- | -------------------------------------------------- |
| left           | switch to the checks pane                          |
| right          | switch to the details pane                         |
| up (checks)    | select the previous group                          |
| down (checks)  | select the next group                              |
| up (details)   | scroll details up                                  |
| down (details) | scroll details down                                |
| f              | filter the groups (hide/show the succeeded groups) |
| d              | show more details                                  |
| h              | display the help                                   |
| q, esc, ctrl-c | quit _Quickonf_                                    |

The _page up_, _page down_, _home_ and _end_ keys may also be used to navigate in checks and details.

# Mouse usage

The buttons in the top-right corner are the following:

| button         | action                                             |
| -------------- | -------------------------------------------------- |
| Filter checks  | filter the groups (hide/show the succeeded groups) |
| More details   | show more details about changes in the right pane  |
| Help           | show or hide the help section                      |
| Quit           | quit _Quickonf_                                    |
