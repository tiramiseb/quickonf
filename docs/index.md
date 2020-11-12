Quickonf allows you to automate your environment configuration, instead of manually configure your apps every time you reinstall your system.

It is especially useful if you reinstall Ubuntu at every release.

## How it works

You create a `quickonf.yaml` file in your system. For instance in your home folder. Then you execute `quickonf`. It reads this file and executes all instructions in it.

In this file, you create steps, with the following format:

```yaml
# Some comments, if you need them
- Step title:
  - instruction:
      parameters

- Another step title:
  - another instruction:
      parameters

[...]
```

Yes, you guessed it, it's written in YAML.

Details about available instructions are listed in this documentation.

(well, they will be...)
