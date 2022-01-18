# Groups

A _Quickonf_ configuration file is composed of multiple sections, called groups. The group name starts at the beginning of the line, instructions are indented.

```text
group name
  instruction
  instruction
  ...
  instruction
```

# Commands

The most important type of instruction is the commands. There are many available commands, and probably many more to come. A command is the instruction given to _Quickonf_ for system checking and modification. The commands may need arguments and may return output values (see _Variables_ below). See the _Commands_ section for details on how to use commands.

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

## Global variables

Some global variables are set when _Quickonf_ starts:

- `hostname`: the hostname of the computer
- `oscodename`: the codename of the current operating system version
- `confdir`: path of the directory containing the configuration file (can be used to read other files in the same directory, templates etc)

## Expanding variables content

Sometimes, a variable's content refers to other variables. Take for instance a file that contains variables names. If you read the file with the `file.read` command, it may be stored in a variable, but the content of the read file is not processed for variables names. In order to process variables names references inside a variable content, simply use the `expand` instruction with the variable name as its only argument.

For instance:

```text
APT sources
  aptsrc = file.read <confdir>/sources.list.tmpl
  expand aptsrc
  file.content /etc/apt/sources.list <aptsrc>
```

... here, the `sources.list.tmpl` file may contain `<oscodename>`.

# Conditional execution

Instructions can be grouped below a `if` instruction, which executes then only if the condition is true.

Currently, the `if` instruction only accepts three arguments:

- a first value
- an operator
- a second value

The operator is one of:

- `=`: condition is true if both values are equal
- `!=`: condition is true of values are different

Of course, this makes sense only when using variables, set by other commands.
