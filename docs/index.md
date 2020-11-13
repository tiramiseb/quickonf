Quickonf allows you to automate your environment configuration, instead of manually configure your apps every time you reinstall your system.

It is especially useful if you reinstall Ubuntu every now and then.

## How it works

You create a `quickonf.yaml` file in your system. For instance in your home folder. Then you execute `quickonf`. It reads this file and executes all instructions in it.

In this file, you create steps, with the following format:

```yaml
# Some comments, if you need them
- Step title:
  - instruction:
      parameters

- Another step title:
  - another instruction:
      parameters

[...]
```

Yes, you guessed it, it's written in YAML.

Instructions are grouped by modules, documented here (yay!).

## Files paths

Whenever an instruction needs a file path, it is either relative to the home directory, or absolute. A path is never relative to the current directory.
