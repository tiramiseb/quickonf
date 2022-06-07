---
Title: Configuration language
short_title: Lang
---

# Groups

A _Quickonf_ configuration file is composed of multiple sections, called groups. The group name starts at the beginning of the line, instructions are indented.

```text
group name
  instruction A
  instruction B
  ...
  instruction G
    instruction H depending on instruction G
    instruction I depending on instruction G
```

Indentations length and type (spaces or tabs) is not strict, you may use whichever you want. Keep in mind that, when mixing indentations, tabs are considered to be equivalent to aligning to multiples of 8 spaces (ie. `<space><space><tab>` is equivalent to 8 spaces).

Some instructions have an impact on multiple other instructions, like conditional execution (`if`) or repetition (`repeat`), in this case you must use a 2nd indentation level (or even more). See below for examples.

# Commands

The most important type of instruction is the commands. There are many available commands, and probably many more to come. A command is the instruction given to _Quickonf_ for system checking and modification. The commands may need arguments and may return output values (see _Variables_ below). See the _Commands_ section for details on how to use commands.

By convention, commands modifying data only for one user are prefixed with "`user.`" and take the username as their first argument. Other commands modify system-wide values.

# Variables

Some commands return output values, that can be stored in variables if needed. Variables are named with alphanumerical characters. Syntax to store values is the following:

```text
group name
  var1 var2 = command arg1 arg2
```

Using variables after assignment is as simple as putting the variable name between "<" and ">" anywhere, like:

```text
group name
  var = command1 arg1 arg2
  command2 <var>
```

... here, the first output of the first command is used as the argument of the second command.

Except for the following global ones, variables are limited to their groups: there is no risk of conflict between variables in different groups.

## Global variables

Some global variables are set when _Quickonf_ starts:

- `hostname`: the hostname of the computer
- `oscodename`: the codename of the current operating system version
- `confdir`: path of the directory containing the configuration file (can be used to read other files in the same directory, templates etc)

## Expanding variables content

Sometimes, a variable's content refers to other variables. Take for instance a file that contains variables names. If you read the file with the `file.read` command, it may be stored in a variable, but the content of the read file is not processed for variables names. In order to process variables names references inside a variable content, simply use the `expand` instruction with the variable name as its only argument.

Example:

```text
APT sources
  aptsrc = file.read <confdir>/sources.list.tmpl
  expand aptsrc
  file.content /etc/apt/sources.list <aptsrc>
```

... here, the `sources.list.tmpl` file may contain `<oscodename>`.

# Conditional execution

Instructions can be grouped below a `if` instruction, which executes then only if the condition is true.

Supported conditions are:

- `file.absent </path/to/file>`: true if the file does not exist
- `file.present </path/to/file>`: true if the file exists
- `<value1> = <value2>`: true if values are equal (compared as strings)
- `<value1> != <value2>`: true if values are different (compared as strings)

Example:

```text
VPN
  if <hostname> = mylaptop
    nm.import openvpn <confdir>/myvpn.ovpn
```

... here, the `nm.import openvpn <confdir>/myvpn.ovpn` is only executed if the host name is "mylaptop". You may of course put multiple instructions with the 2nd indentation level.

# Repeating identical commands and arguments

Say you want to apply similar commands multiple times. Instead of writing the same prefix in multiple lines, you can use the `repeat` instruction. It makes it easier to write and to read your configuration file. It takes as arguments the command to be repeated, optionally with its first arguments.

_Quickonf_ translates it by prepending this command name and these arguments to all lines.

Example:

```text
Wifi
  repeat nm.wifi
    ssid1 psk1
    ssid2 psk2
    ssid3 psk3

My XDG dirs
  repeat user.xdg.userdir alice
    desktop DESKtop
    documents MyDocs
    download UglyStuff
```

... translates internally to:

```text
Wifi
  nm.wifi ssid1 psk1
  nm.wifi ssid2 psk2
  nm.wifi ssid3 psk3

My XDG dirs
  user.xdg.userdir alice desktop DESKtop
  user.xdg.userdir alice documents MyDocs
  user.xdg.userdir alice download UglyStuff
```

# Priority

Some instructions may need to be executed before other ones, when executing all instructions. For instance, changing the packages sources or adding an installer like snap or flatpak. For that specific case, You may use the `priority` instruction, which takes an integer as an arguments. The default priority is 0, and larger value means higher priority.

If you use the `priority` instruction in a conditional block, it impacts the whole group (meaning you may give different priorities depending on conditions).

Keep in mind that it only impacts the order of execution and display in the interface. You still can manually apply instructions of lower priority before instructions of higher priority.

For instance:

```text
Flatpak
  priority 1
  apt.install flatpak
```

# Multi-level indentation

You may use multiple level indentation, if needed by the instructions you want to use.

Example:

```text
Wifi
  nm.wifi home homessid
  if <hostname> = worklaptop
    repeat nm.wifi
      worknetwork workssid
      mobilenetwork mobilessid
```

... here, the "home" network will be known by all computers, but "worknetwork" and "mobilenetwork" are only known by your work laptop.
