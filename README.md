# Quickonf: a tool to quickly (re)configure a Linux system

[![Download latest release](https://img.shields.io/github/v/release/tiramiseb/quickonf?display_name=tag&sort=semver&label=Download&style=for-the-badge)](https://github.com/tiramiseb/quickonf/releases/latest/download/quickonf)
[![Read the documentation](https://img.shields.io/badge/Read-documentation-blue?style=for-the-badge)](https://tiramiseb.github.io/quickonf/intro.html)

## Introduction

_Quickonf_ is a semi-automatic configuration tool for Linux desktop computers, allowing you to reinstall your Linux system and quickly get your favorite environment back.

It could for instance be useful if you reinstall _Ubuntu_ every now and then and do not want to use the same user environment over and over (because sometimes, when keeping older configuration files, some app may dysfunction).

Think Ansible, but as a single binary command tailored for local desktop usage.

## How it works

You write the desired state of your system in a configuration file (named `quickonf.qconf` by default).

Then you execute `quickonf` or `quickonf your-file.qconf`. _Quickonf_ reads this file and checks the status of the systems, then allows you to apply modifications to the system. No other argument is needed, everything is done in an interactive interface.

---

[![License](https://img.shields.io/github/license/tiramiseb/quickonf)](https://github.com/tiramiseb/quickonf/blob/main/LICENSE)
![Issues](https://img.shields.io/github/issues/tiramiseb/quickonf)
![Commits](https://img.shields.io/github/commits-since/tiramiseb/quickonf/latest/main?sort=semver)
![Last commit](https://img.shields.io/github/last-commit/tiramiseb/quickonf/main)
