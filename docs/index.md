Quickonf allows you to automate your environment configuration, instead of manually configure your apps every time you reinstall your system.

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

## General knowledge

### Paths

Whenever an instruction needs a file or directory path, it is either relative to the home directory, or absolute. A path is never relative to the current directory.

### Failures

When an instruction fails, the whole step is stopped, but the processing of other steps continue. Generally, when defining steps, you should not rely on previous steps having run correctly. However, inside a single step, you are assured that the previous instructions have been run.

### Multiple uses of an instruction

You can use the same instruction multiple times in a step. For instance, if you must create a directory then install a package then create a new directory, it is possible.

### Data store

Sometimes, you want to store some data to be reused in another instruction. For this, some instructions allow storing data in a centralized key-value store. That data can then be reused by putting the data key anywhere in the keys or values, in the following form: `<key-name>`. If a key with that name exists, that string is replaced with the associated value.

By default, the `home` key contains the user's home directory path.
