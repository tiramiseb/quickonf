Quickonf allows you to automate your environment configuration, instead of manually configuring your apps every time you reinstall your system.

It is especially useful if you reinstall Ubuntu every now and then.

Think Ansible, but as a single binary command tailored for local desktop usage.

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

## Flags and arguments

When run with the `-list` or `-l` flag, `quickonf` lists all defined steps, without changing or setting anything.

When run with the `-dry-run` or `-r` flag, `quickonf` runs in dry-run mode, without modifying the system, allowing to test the configuration.

Of course, `-help` or `-h` displays the flags list.

## General knowledge

### Command output

Quickonf displays everything it does on the terminal, in green for successful operations and in red for failures.

At the end of the execution, all failures are listed in a "Report" section. If this section is empty, it means that the configuration has been successfully applied.

### Idempotence

The goal of this app is to be idempotent. It means that you can execute it as many times as you want, the result will always be the same. Of course, you may break this logic in the configuration file, that's why you sometimes need to be vigilant about the instructions you use, and that's why there are some recipes.

It will also allow you to run the command on a regular basis to ensure all applications are up-to-date, even if they are not installed with a package manager.

### Paths

Whenever an instruction needs a file or directory path, it is either relative to the home directory, or absolute. A path is never relative to the current directory.

### Failures

When an instruction fails, the whole step is stopped, but the processing of other steps continue. Generally, when defining steps, you should not rely on previous steps having run correctly. However, inside a single step, you are assured that the previous instructions have run.

### Multiple uses of an instruction

You can use the same instruction multiple times in a step. For instance, if you must create a directory then install a package then create another directory, it is possible.

### Data store

Sometimes, you want to store some data to be reused in another instruction. For this, some instructions allow storing data in a centralized key-value store. That data can then be reused by putting the data key anywhere in the keys or values, in the following form: `<key-name>`. If a key with that name exists, that string is replaced with the associated value. It is used in some recipes, don't hesitate to check them.

The following keys are set by default in the store:

- `home`: user's home directory path
- `oscodename`: Linux distribution codename (result of `lsb_release --codename --short`)
