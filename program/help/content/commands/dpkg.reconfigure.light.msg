
[38;5;228;48;5;63;1m dpkg.reconfigure [0m[38;5;234m[38;5;234m [0m[38;5;234m                                              [0m
[0m[38;5;234m                                                                 [0m
[38;5;234mReconfigure an installed dpkg package                            [0m
[38;5;234m[38;5;234m [0m[38;5;234m                                                                [0m
[0m[38;5;27;1m## Arguments[0m[38;5;234m[38;5;234m [0m[38;5;234m                                                    [0m
[0m[38;5;234m                                                                 [0m
[38;5;234m• Name of the package                                            [0m
[38;5;234m[38;5;234m [0m[38;5;234m                                                                [0m
[0m[38;5;27;1m## Example[0m[38;5;234m[38;5;234m [0m[38;5;234m                                                      [0m
[0m[38;5;234m                                                                 [0m
[38;5;235m[0m[38;5;234m  [0m[38;5;235mlibdvd-pkg[38;5;234m [0m[38;5;234m                                                    [0m
[0m[38;5;234m  [0m[38;5;235m  debconf.set libdvd-pkg libdvd-pkg/first-install .[38;5;234m [0m[38;5;234m           [0m
[0m[38;5;234m  [0m[38;5;235m  debconf.set libdvd-pkg libdvd-pkg/post-invoke_hook-install true[0m
[0m[38;5;234m  [0m[38;5;235m  apt.install libdvd-pkg[38;5;234m [0m[38;5;234m                                      [0m
[0m[38;5;234m  [0m[38;5;235m  dpkg.reconfigure libdvd-pkg[38;5;234m [0m[38;5;234m                                 [0m
[0m
