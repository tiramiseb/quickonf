---
title: GNOME
---

## Install Vanilla GNOME on Ubuntu

```yaml
- Initialisation:
    - sudo-password: my-wonderful-password

- Vanilla GNOME:
    - apt:
        - gnome-session
    - update-alternatives:
        gdm3-theme.gresource: /usr/share/gnome-shell/gnome-shell-theme.gresource
```

## Some GNOME Shell preferences

Here are some of my preferences, don't hesitate to dig into the dconf database to find the ones that suit you.

```yaml
- GNOME Shell Preferences:
    - dconf:
        /apps/update-manager/show-details: "true"
        /org/gnome/desktop/interface/clock-show-date: "true"
        /org/gnome/desktop/interface/buttons-have-icons: "true"
        /org/gnome/settings-daemon/plugins/xsettings/antialiasing: "'rgba'"
        /org/gnome/settings-daemon/plugins/color/night-light-enabled: "true"
        /org/gnome/desktop/wm/preferences/action-middle-click-titlebar: "'lower'"
        /org/gnome/desktop/wm/preferences/focus-mode: "'sloppy'"
        /org/gnome/desktop/wm/preferences/button-layout: "'appmenu:close'"
        /org/gnome/desktop/screensaver/lock-enabled: "false"
        /org/gnome/desktop/screensaver/ubuntu-lock-on-suspend: "false"
        /org/gnome/desktop/media-handling/autorun-never: "true"
        /org/gnome/desktop/wm/keybindings/move-to-workspace-down: "['<Control><Shift>Down']"
        /org/gnome/desktop/wm/keybindings/move-to-workspace-left: "@as []"
        /org/gnome/desktop/wm/keybindings/move-to-workspace-right: "@as []"
        /org/gnome/desktop/wm/keybindings/move-to-workspace-up: "['<Control><Shift>Up']"
        /org/gnome/desktop/wm/keybindings/switch-to-workspace-down: "['<Control>Down']"
        /org/gnome/desktop/wm/keybindings/switch-to-workspace-left: "@as []"
        /org/gnome/desktop/wm/keybindings/switch-to-workspace-right: "@as []"
        /org/gnome/desktop/wm/keybindings/switch-to-workspace-up: "['<Control>Up']"
        /org/gnome/settings-daemon/plugins/media-keys/next: "['<Control>Right']"
        /org/gnome/settings-daemon/plugins/media-keys/play: "['<Control>KP_0']"
        /org/gnome/settings-daemon/plugins/media-keys/previous: "['<Control>Left']"
```

## GNOME Shell extensions

```yaml
- GNOME Shell extensions:
    - apt:
        - gnome-shell-extension-hide-activities
    - gnome-shell-extension:
        - places-menu@gnome-shell-extensions.gcampax.github.com
        - Hide_Activities@shay.shayel.org
    - gnome-shell-restart:
```

## Install Pop Shell

```yaml
- Pop Shell:
    - local-gnome-shell-extension-version:
        extension: pop-shell@system76.com
        store: pop-shell-current
    - github-latest:
        repository: pop-os/shell
        store: pop-shell-candidate
        pattern: pop-shell@system76.com_*.zip
        store-url: pop-shell-url
    - stop-if-older:
        current: <pop-shell-current>
        candidate: <pop-shell-candidate>
    - download:
        <pop-shell-url>: /tmp/quickonf-pop-shell.zip
    - unzip:
        /tmp/quickonf-pop-shell.zip: .local/share/gnome-shell/extensions/pop-shell@system76.com
    - remove:
        - /tmp/quickonf-pop-shell.zip
    - gnome-shell-extension:
        - pop-shell@system76.com
    - gnome-shell-restart:
```
