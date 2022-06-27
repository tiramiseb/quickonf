
[38;5;228;48;5;63;1m dpkg.reconfigure [0m[38;5;252m[38;5;252m [0m[38;5;252m                                              [0m
[0m[38;5;252m                                                                 [0m
[38;5;252mReconfigure an installed dpkg package                            [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                                [0m
[0m[38;5;39;1m## Arguments[0m[38;5;252m[38;5;252m [0m[38;5;252m                                                    [0m
[0m[38;5;252m                                                                 [0m
[38;5;252m• Name of the package                                            [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                                [0m
[0m[38;5;39;1m## Example[0m[38;5;252m[38;5;252m [0m[38;5;252m                                                      [0m
[0m[38;5;252m                                                                 [0m
[38;5;251m[0m[38;5;252m  [0m[38;5;251mlibdvd-pkg[38;5;252m [0m[38;5;252m                                                    [0m
[0m[38;5;252m  [0m[38;5;251m  debconf.set libdvd-pkg libdvd-pkg/first-install .[38;5;252m [0m[38;5;252m           [0m
[0m[38;5;252m  [0m[38;5;251m  debconf.set libdvd-pkg libdvd-pkg/post-invoke_hook-install true[0m
[0m[38;5;252m  [0m[38;5;251m  apt.install libdvd-pkg[38;5;252m [0m[38;5;252m                                      [0m
[0m[38;5;252m  [0m[38;5;251m  dpkg.reconfigure libdvd-pkg[38;5;252m [0m[38;5;252m                                 [0m
[0m
