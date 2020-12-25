# Developers documentation

Hey dev! Pull Requests are welcome! You will find some insights about Quickonf development here!

## Instructions

Instructions are where all the magic takes place.

Each instruction must be a function with the following signature:

```go
func(in interface{}, out output.Output) error
```

`in` is the data from the YAML file, and `out` is the output writer. If the instruction fails at any time, simply return an error, no need to use `out.Error` manually.

The very first step is to show which instruction is running, with:

```go
out.InstructionTitle("Describe the instruction here")
```

Then, you must extract the input data, either as a string, as a slice of strings or as a map of strings to strings, with on of the following calls:

```go
data, err := helper.String(in)
data, err := helper.SliceString(in)
data, err := helper.MapStringString(in)
```

These are the only three possibilites for now, and it should be kept as single as possible: let's try to not have more formats!

(Don't forget to test `err`!)

Then, either you loop on the data or you read it by key, depending on what you need.

Before each modification of the system, you must check the value of the `Dryrun` global variable and act consequently. Please note some helpers (see below) already check it, so you may not need to check it before. Instructions with no impact on the system (like parsing, data transformation, etc) do not need to check it.

In the end, the general signature of an instruction is something like:

```go
// FooBar does foo to bar
func FooBar(in interface{}, out output.Output) error {
    out.InstructionTitle("Do foo to bar")
    data, err := helper.SliceString(in)
    if err != nil {
        return err
    }
    for _, bar := range data {
        if Dryrun {
            out.Infof("Would do foo to %s", bar)
            continue
        }
        out.Infof("Doing foo to %s", bar)
        # Do foo to bar
    }
    return nil
}
```

Please keep in mind that an instruction must be global and reusable. Except for some very specific cases, instructions must not be directly related to some software.

On the other side, instructions must not be too generic. For instance, an instruction to simply execute the command that would be given in the YAML file is forbidden.

## Modules

Modules are simply files in the `internal/modules` directory. Just a way to group instructions together.

Group all related instructions in a single file in this directory and don't forget to call the `Register` function for each instruction. For instance:

```go
func init() {
    Register("foo-bar", FooBar)
}
```

Generally, the instruction as used in the YAML file is written in lowercase with words separated by a dash, and the name of the function is in CamelCase, with the same words.

Of course, don't hesitate to take inspiration in other modules/instructions!

## Documentation

Each module must have a corresponding documentation, in `docs/_modules`, named after the module filename. It must begin with a YAML metadata header with only a `title` field, and its content must begin with a table made of four columns, describing instructions available in the module:

- Instruction: the name of the instruction as used in the YAML file
- Action: short description of the instruction
- Arguments: short description of the arguments
- Dry run: specifics about dry-run mode (if there is no difference, the simply write "-")

Example:

```markdown
---
title: Foo Bar
---

| Instruction | Action         | Arguments    | Dry run   |
| ----------- | -------------- | ------------ | --------- |
| `foo-bar`   | Do foo to bars | List of bars | No change |
```

If you need to give some explanation, add some chapters below the table (with level-2 titles).

Don't hesitate to check other documentation pages!

## Output

The output allows giving information to the user during the Quickonf execution. The default (and currently only) output is printing on standard output and storing in a report, which is written at the end of the execution.

The output system offers multiple functions...

The ones you must not use in instructions code:

- `StepTitle`: used at a higher level, print the name of the step (DO NOT use it in instructions)
- `Error`: used at a higher level, print an error returned by an instruction (DO NOT use it in instructions)
- `Report`: write the final report

The ones you may use to print messages:

- `InstructionTitle`: used only once per instruction, at its very beginning
- `Info` and `Infof`: write informational messages
- `Success` and `Successf`: inform the user that a modification was successful
- `Alert` and `Alertf`: write an alert message (use them sparingly)

The ones you may use to make the user wait:

- `ShowLoader` and `HideLoader`: show and hide a loading bar (be careful about hiding, make sure to call HideLoader before returning an error, for instance)
- `ShowPercentage` and `HidePercentage`: if you know the percentage of what the user is waiting for, use these
- `ShowXonY` and `HideXonY`: if there is a total count and you know (and you can update) the current position, use these

## Helpers

Because some actions may be common to multiple instructions, there are some helpers in `internal/helpers`. You may of course create new helpers if you think it will help.

Currently available helpers are:

- download:

  - `DownloadFileWithPercent(url, path string, out output.Output) error`: download a file from the given URL and store it at the given path, showing a percentage loader (already tests Dryrun, no need to check it before calling this helper)
  - `DownloadFile(url, path string) error`: download a file from the given URL and store it at the given path, silently (don't use it for large files, please) (already tests Dryrun, no need to check it before calling this helper)
  - `DownloadJSON(url string, destination interface{}) error`: download JSON data from the given URL and store it in the given interface
  - `Download(url string) ([]byte, error)`: download data from the given URL and return it
  - `Post(url string, payload []byte) ([]byte, error)`: execute a POST method request to the given URL with the given payload (content-type is automatically detected) and return the resulting data

- exec (for these helpers, please note the current environment is passed to the executed commands, along with `LANG=C` to avoid locale-specific values):

  - `Exec(env []string, cmd string, args ...string) ([]byte, error)`: execute the given command, with the given arguments, providing the given environment variables, and return its stdout (or stderr as a part of the error)
  - `ExecSudo(env []string, args ...string) ([]byte, error)`: execute the command (given as the first argument) as root, with the other arguments, providing the given environment variables, and return its stdout (or stderr as a part of the error)

- git:

  - `GitClone(repo, dest string, depth int, out output.Output) error`: clone the given repository to the given destination, with the given depth, while showing a loader if `out` is not nil (already tests Dryrun, no need to check it before calling this helper)

- path:

  - `Path(str string) string`: transform the given path as an absolute path. The relative paths are considered to be relative to the user's home directory. Please use this helper everywhere you ask for a path.

- symlink:

  - `Symlink(path, target string) (status SymlinkStatus, err error)`: create a symbolic link, making `path` point to `target` (already tests Dryrun, no need to check it before calling this helper)

- zip:

  `UnzipFile(zipfilepath, dest string, out output.Output) error`: extract the given .zip file to the given destination, writing information and showing a X/Y progress bar if `out` is not nil (already tests Dryrun, no need to check it before calling this helper)

## Recipes

Recipes are the place where you may share configuration for some specific situation or software. Instructions should not be specific, but recipes are exactly made for it!
