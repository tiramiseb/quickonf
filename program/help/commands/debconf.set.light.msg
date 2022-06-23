
[38;5;228;48;5;63;1m debconf.set [0m[38;5;234m[38;5;234m [0m[38;5;234m                                                                      [0m
[0m[38;5;234m                                                                                    [0m
[38;5;234mSet a debconf parameter                                                             [0m
[38;5;234m[38;5;234m [0m[38;5;234m                                                                                   [0m
[0m[38;5;27;1m## Arguments[0m[38;5;234m[38;5;234m [0m[38;5;234m                                                                       [0m
[0m[38;5;234m                                                                                    [0m
[38;5;234m• Package where the parameter belong                                                [0m
[38;5;234m• Name of the parameter                                                             [0m
[38;5;234m• Value to apply to the parameter                                                   [0m
[38;5;234m[38;5;234m [0m[38;5;234m                                                                                   [0m
[0m[38;5;27;1m## Example[0m[38;5;234m[38;5;234m [0m[38;5;234m                                                                         [0m
[0m[38;5;234m                                                                                    [0m
[38;5;235m[0m[38;5;234m  [0m[38;5;235mInstall MS fonts[38;5;234m [0m[38;5;234m                                                                 [0m
[0m[38;5;234m  [0m[38;5;235m  debconf.set ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula true[0m
[0m[38;5;234m  [0m[38;5;235m  apt.install ttf-mscorefonts-installer[38;5;234m [0m[38;5;234m                                          [0m
[0m
