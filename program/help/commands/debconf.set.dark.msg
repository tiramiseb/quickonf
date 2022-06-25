
[38;5;228;48;5;63;1m debconf.set [0m[38;5;252m[38;5;252m [0m[38;5;252m                                                                      [0m
[0m[38;5;252m                                                                                    [0m
[38;5;252mSet a debconf parameter                                                             [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                                                   [0m
[0m[38;5;39;1m## Arguments[0m[38;5;252m[38;5;252m [0m[38;5;252m                                                                       [0m
[0m[38;5;252m                                                                                    [0m
[38;5;252m• Package where the parameter belong                                                [0m
[38;5;252m• Name of the parameter                                                             [0m
[38;5;252m• Value to apply to the parameter                                                   [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                                                   [0m
[0m[38;5;39;1m## Example[0m[38;5;252m[38;5;252m [0m[38;5;252m                                                                         [0m
[0m[38;5;252m                                                                                    [0m
[38;5;251m[0m[38;5;252m  [0m[38;5;251mInstall MS fonts[38;5;252m [0m[38;5;252m                                                                 [0m
[0m[38;5;252m  [0m[38;5;251m  debconf.set ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula true[0m
[0m[38;5;252m  [0m[38;5;251m  apt.install ttf-mscorefonts-installer[38;5;252m [0m[38;5;252m                                          [0m
[0m
