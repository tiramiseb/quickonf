# Intro

Quickonf allows you to automate your environment configuration, instead of manually configuring your apps every time you reinstall your system. It is made for Linux. Or maybe other UNIX-Like system, dunno, not tested.

It is especially useful if you reinstall Ubuntu every now and then and do not want to use the same user environment over and over (because sometimes, when keeping older configuration files, some app may disfunction).

Think Ansible, but as a single binary command tailored for local desktop usage.

## How it works

You create a file, named `quickonf.yaml` by default, in your system. For instance in your home folder. Then you execute `quickonf` or `quickonf -c your-file.yaml`. It reads this file and executes all instructions in it.

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

Instructions are grouped by modules, documented here, menu is on the left (yay!).

## How to install

Direct links:

- [64 bits binary](https://github.com/tiramiseb/quickonf/releases/latest/download/quickonf)
- [32 bits binary](https://github.com/tiramiseb/quickonf/releases/latest/download/quickonf-32)

Don't forget to make the downloaded file executable!

## Flags and arguments

The `-config` or `-c` flag, associated with a file name, defines the configuration file to read. Its default value is `quickonf.yaml`.

When run with the `-list` or `-l` flag, `quickonf` lists all defined steps, without changing or setting anything.

When run with the `-dry-run` or `-r` flag, `quickonf` runs in dry-run mode, without modifying the system, allowing to test the configuration.

All non-option arguments are considered to be patterns to select specific steps you want to run (see below).

And, finally, of course, as usual, as you may expect, don't worry, `-help` or `-h` displays the flags list.

## General knowledge

### Steps selection

You can chose a set of steps to run, by giving patterns as non-option arguments. These arguments are path-compatible patterns, which will be tested against lower-case version of steps names, and which will also match if only a part of the name matched. For instance, "gnome ext" will match "Install GNOME extensions", "Configure GNOME extensions", but not "Install GNOME" nor "Configure GNOME Shell".

However, you may need some steps to run even if not explicitly selected. In that case, use `always: true` at the same level as steps instructions (see recipes for examples).

### Command output

Quickonf displays everything it does on the terminal, in green for successful operations and in red for failures.

At the end of the execution, all failures are listed in a "Report" section. If this section is empty, it means that the configuration has been successfully applied.

### Idempotence

The goal of this app is to be idempotent. It means that you can execute it as many times as you want, the result will always be the same. Of course, you may break this logic in the configuration file, that's why you sometimes need to be vigilant about the instructions you use, and that's why there are some recipes.

It will also allow you to run the command on a regular basis to ensure all applications are up-to-date, even if they are not installed with a package manager.

### Paths

Whenever an instruction needs a file or directory path, it is either relative to the home directory, or absolute. A path is never relative to the current directory.

### Migration

The [move]({% link _modules/move.md %}) module includes two instructions for migration, called `migration-source` and `migrate`.

Let's say you have moved your previous home directory to `/home/previous` before reinstalling your system, in order to have a clean all-new environment. You still want to keep some directories or files as-is. That's what migration is about.

### Failures

When an instruction fails, the whole step is stopped, but the processing of other steps continue. Generally, when defining steps, you should not rely on previous steps having run correctly. However, inside a single step, you are assured that the previous instructions have run.

### Multiple uses of an instruction

You can use the same instruction multiple times in a step. For instance, if you must create a directory then install a package then create another directory, it is possible.

### Data store

Sometimes, you want to store some data to be reused in another instruction. For this, some instructions allow storing data in a centralized key-value store. That data can then be reused by putting the data key anywhere in the keys or values, in the following form: `<key-name>`. If a key with that name exists, that string is replaced with the associated value. It is used in some recipes, don't hesitate to check them.

The following keys are set by default in the store:

- `home`: user's home directory path
- `oscodename`: Linux distribution codename (result of `lsb_release --codename --short`)
- `hostname`: System hostname

### Run as another user

As long as the current user can use sudo to execute commands as another user, any instruction may be executed as another user by adding `@<username>` at the end of its name. For instance:

```yaml
- Create Bayek:
    - user-password:
        bayek: aya

- Hello Bayek:
    - file@bayek:
        HELLO.txt: |
          Hello Bayek, you have been automatically created!
```

In order to allow all instructions to be executed by any other user, a new instance is executed as the given user, with only the given instruction in its configuration file, and its output is redirected to the current quickonf instance. It implies the `quickonf` file must be executable for the given user.
