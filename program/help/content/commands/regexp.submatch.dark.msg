
[38;5;228;48;5;63;1m regexp.submatch [0m[38;5;252m[38;5;252m [0m[38;5;252m                                     [0m
[0m[38;5;252m                                                       [0m
[38;5;252mFind submatches using a regexp                         [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                      [0m
[0m[38;5;39;1m## Arguments[0m[38;5;252m[38;5;252m [0m[38;5;252m                                          [0m
[0m[38;5;252m                                                       [0m
[38;5;252m• Regexp                                               [0m
[38;5;252m• Source string                                        [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                      [0m
[0m[38;5;39;1m## Outputs[0m[38;5;252m[38;5;252m [0m[38;5;252m                                            [0m
[0m[38;5;252m                                                       [0m
[38;5;252m• First submatch                                       [0m
[38;5;252m• Second submatch                                      [0m
[38;5;252m• ...                                                  [0m
[38;5;252m[38;5;252m [0m[38;5;252m                                                      [0m
[0m[38;5;39;1m## Example[0m[38;5;252m[38;5;252m [0m[38;5;252m                                            [0m
[0m[38;5;252m                                                       [0m
[38;5;251m[0m[38;5;252m  [0m[38;5;251mFind src[38;5;252m [0m[38;5;252m                                            [0m
[0m[38;5;252m  [0m[38;5;251m  webpage = http.get.var http://www.example.com[38;5;252m [0m[38;5;252m     [0m
[0m[38;5;252m  [0m[38;5;251m  src = regexp.submatch "<script .*src=(.*)>"[38;5;252m [0m[38;5;252m       [0m
[0m[38;5;251m<webpage>[38;5;252m [0m[38;5;252m                                             [0m
[0m[38;5;252m  [0m[38;5;251m  ...[38;5;252m [0m[38;5;252m                                               [0m
[0m
