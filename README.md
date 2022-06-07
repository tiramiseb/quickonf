# Quickonf: a tool to quickly (re)configure a Linux system

## Introduction

_Quickonf_ is a semi-automatic configuration tool for Linux desktop computers, allowing you to reinstall your Linux system and quickly get your favorite environment back.

It is especially useful if you reinstall _Ubuntu_ every now and then and do not want to use the same user environment over and over (because sometimes, when keeping older configuration files, some app may disfunction).

Think Ansible, but as a single binary command tailored for local desktop usage.

## How it works

You write the desired state of your system in a configuration file (named `quickonf.qconf` by default).

Then you execute `quickonf` or `quickonf your-file.qconf`. _Quickonf_ reads this file and checks the status of the systems, then allows you to apply modifications to the system. No other argument is needed, everything is done in an interactive interface.

## Want to try it?

See the documentation: <https://tiramiseb.github.io/quickonf/intro.html>.
