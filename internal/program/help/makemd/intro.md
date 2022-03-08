# Introduction

_Quickonf_ is a semi-automatic configuration tool for Linux desktop computers, allowing you to reinstall your Linux system and quickly get your favorite environment back.

It is especially useful if you reinstall _Ubuntu_ every now and then and do not want to use the same user environment over and over (because sometimes, when keeping older configuration files, some app may disfunction).

Think Ansible, but as a single binary command tailored for local desktop usage.

# How it works

You write the desired state of your system in a configuration file (named `quickonf.qconf` by default) and _Quickonf_ ensures the computer is configured as you wish.

Then you execute `quickonf` or `quickonf your-file.qconf`. _Quickonf_ reads this file and checks the status of the systems, then allows you to apply modifications to the system. No other argument can be provided, everything is done in an interactive interface.

# General knowledge

## Groups

The configuration is split in groups, a group could represent a part of your configuration. A group is composed of instructions, which are started one after another. As long as an instruction does not fail, the next instructions are executed. If an instruction fails, execution of all instructions in the same group is aborted.

Moreover, these groups are completely independent, and except for some special cases, you should make sure there is no dependency between two groups.

## Idempotence

The goal of this app is to be idempotent. It means that you can execute it as many times as you want, the result will always be the same. Of course, you may break this logic in the configuration file, that's why you sometimes need to be vigilant about the instructions you use, and that's why there are some recipes.

It will also allow you to run the command on a regular basis to ensure all applications are up-to-date, even if they are not installed with a package manager.

## Run as root and paths

_Quickonf_ must be run as root, (eg. using `sudo`) because a large part of a system configuration is done as root (installing packages, etc). Moreover, it is easier to execute commands and modify files as root, allowing access to any part of the system.

Therefore, most paths given as arguments of the instructions must be absolute. Some user-specific instructions allow relative paths, which are relative to the relevant user.
