With _Quickonf_, you are able to configure your Linux computers, and keep the configuration you want, using a simgle configuration file.

You may execute it whenever you want, to check your system is configured as needed or to fix some configuration you should not have changed.

After reinstalling your system or anytime you want, you simply have to run _Quickonf_ and voilà, your system is (re)configured exactly how you had it before (as long as your configuration is described in your _Quickonf_ configuration file, of course)!

## Example overview

Let's look at some extracts of a potential _Quickonf_ configuration file...

{{<example>}}

```plain
APT sources
	priority 1
	aptsrc = file.read <confdir>/apt-sources.list
	expand aptsrc
	file.content /etc/apt/sources.list <aptsrc>
	apt.upgrade
```

This group allows configuring the APT sources and upgrading packages.

The `apt-sources.list` file contains:

```plain
# Sources from Quickonf
deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename> main restricted universe multiverse
deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename>-updates main restricted universe multiverse
deb http://fr.archive.ubuntu.com/ubuntu/ <oscodename>-backports main restricted universe multiverse
deb http://security.ubuntu.com/ubuntu/ <oscodename>-security main restricted universe multiverse
deb http://archive.canonical.com/ubuntu <oscodename> partner
```

- `priority 1` means the "APT sources" group will be ordered before others
- `file.read` reads the content of the file `apt-sources.list` which is placed in the same directory as the _Quickonf_ configuration, and puts its content into variable `aptsrc`
- `expand aptsrc` replaces occurrences of variables names with their values in variable `aptsrc`
- `file.content` puts the content of variable `aptsrc` into file `/etc/apt/sources.list`
- `apt.upgrade` updates the packages list and upgrades all needed packages

{{</example>}}

{{<example>}}

```plain
Basic tools
	repeat apt.install
		baobab
		gnome-tweaks
		openssh-server
		net-tools
```

This groups installs the `baobab`, `gnome-tweaks`, `openssh-server` and `net-tools` packages, by repeating the `apt.install` command over them.

{{</example>}}

{{<example>}}

```plain
Vanilla GNOME
	apt.install gnome-session
	alternatives.set gdm-theme.gresource /usr/share/gnome-shell/gnome-shell-theme.gresource
```

This group installs the `gnome-session` package and changes the _alternative_ for `gdm-theme` to the regular GNOME Shell theme, providing the vanilla experience of GNOME on Ubuntu.

{{</example>}}

{{<example>}}

```plain
BitScope DSO
	webpage = http.get.var http://my.bitscope.com/download/?p=download&f=APDX
	version = regexp.submatch "input type=\"checkbox\" name=\"row[]\" .*value=\"[0-9]+\".*</td>[[:space:]]*<td>bitscope-dso_(.*)_amd64.deb" <webpage>
	current = dpkg.version bitscope-dso
	if <version> != <current>
		tmpdest = tempdir
		http.get.file http://bitscope.com/download/files/bitscope-dso_<version>_amd64.deb <tmpdest>/bitscope.deb
		dpkg.install <tmpdest>/bitscope.deb
```

This groupe installs or upgrades the `BitScope DSO` package, which is not available in a repository.

- `http.get.var` reads the given webpage and puts its content in variable `webpage`.
- `regexp.submatch` looks for the version of the latest available package for the `bitscope-dso_*_amd64.deb` file in the `webpage` variable, and extracts it into variable `version`.
- `dpkg.version` checks the current version of package `bitscope-dso` and stores it into variable `current`
- `if` compares the `version` and `current` variables, and executes the next lines only if they are different
- `tempdir` creates a temporary directory (which will be removed when closing _Quickonf_) and puts its path into `tmpdest`
- `http.get.file` downloads the package to `bitscope.deb` in the temporary directory
- `dpkg.install` installs the `bitscope.deb` package

{{</example>}}
