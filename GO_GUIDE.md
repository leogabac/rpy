# rpy Go Guide

This guide is for someone who has never really written Go before.

The goal is not to teach all of Go. That would be a waste of your time right now.
The goal is to get you from zero to a small working CLI that does a few real things:

- builds with Go
- uses `cobra-cli`
- creates a local Python virtual environment
- runs commands inside that environment
- installs packages with `pip`

If you finish this guide, you should understand enough Go to build the first version of `rpy` without feeling like you are blindly copying code.

## 1. What you are making

Ignore the full project vision for a minute. The first version of `rpy` should be much smaller than the final idea.

You are building a command line program called `rpy` with commands like:

```bash
rpy doctor
rpy env create
rpy env path
rpy env info
rpy run python --version
rpy pip install numpy
```

That is enough to learn the important parts:

- how Go programs are structured
- how Cobra organizes commands
- how to work with files and paths
- how to run other programs from Go
- how to model state with structs and methods

Do not start with central environments, `uv`, Python downloading, or config parsing.
Those are later problems.

## 2. What Go is good at here

Go is a very practical language for CLI tools.

It is a good fit for `rpy` because it is:

- compiled to a single binary
- fast enough without trying hard
- good at filesystem and subprocess work
- simple enough that small tools stay readable

It is not an “object-oriented” language in the Java or C# sense.
There are no classes.
Instead, Go mostly uses:

- structs for data
- methods for behavior
- interfaces for shared behavior when needed

That is enough for this project.

## 3. Install the minimum tools

You need:

- Go
- `cobra-cli`
- Python

Check Go:

```bash
go version
```

Check Python:

```bash
python --version
```

Install `cobra-cli` if needed:

```bash
go install github.com/spf13/cobra-cli@latest
```

If `cobra-cli` is not found afterward, your Go bin directory is probably not on `PATH`.
That is commonly one of:

- `~/go/bin`
- `$GOPATH/bin`

## 4. What a Go project looks like

A very small Go CLI often looks like this:

```text
main.go
go.mod
cmd/
  root.go
  doctor.go
  env.go
  env_create.go
  env_path.go
  env_info.go
  run.go
  pip.go
  pip_install.go
internal/
  app/
  env/
  pyexec/
```

What these mean:

- `main.go` is the real entry point
- `go.mod` defines your module and dependencies
- `cmd/` contains Cobra command definitions
- `internal/` contains the actual program logic

This split matters.
If you put all your logic in `cmd/*.go`, the program becomes hard to test and hard to change.
The Cobra command files should mostly translate CLI input into function calls.

## 5. Your first mental model of Go

You do not need advanced Go yet. You need a small subset.

### Variables

Go has explicit types, but it can infer many of them.

```go
name := "rpy"
count := 3
enabled := true
```

Use `:=` when declaring a new variable inside a function.

Use `var` when you need the zero value or want to declare first and assign later:

```go
var err error
var path string
```

### Functions

Functions look like this:

```go
func greet(name string) string {
    return "hello " + name
}
```

The structure is:

- `func`
- function name
- parameters
- return type
- body

Some functions return multiple values:

```go
func findPython() (string, error) {
    return "python", nil
}
```

That is normal in Go.

### Errors

Go does not use exceptions for normal flow.
Most operations return an `error`, and you check it immediately.

```go
path, err := os.Getwd()
if err != nil {
    return err
}
```

This style can feel repetitive, but it keeps control flow explicit.

For this project, your default rule should be:

- if a function can fail, return an `error`
- if you call something that returns an `error`, handle it right away

### Slices

A slice is Go’s main list type:

```go
args := []string{"python", "--version"}
```

You will use slices constantly for command arguments, file lists, and package names.

### Maps

A map is a key/value table:

```go
settings := map[string]string{
    "backend": "venv",
    "mode":    "local",
}
```

You probably will not need maps much in the earliest phase.

### Structs

A struct is a named group of fields:

```go
type EnvSpec struct {
    Name   string
    Path   string
    Python string
}
```

This is how Go groups related data.
If you are coming from Python, think of it as a very lightweight object with fixed fields.

### Methods

A method is a function attached to a type:

```go
func (e EnvSpec) PythonPath() string {
    return filepath.Join(e.Path, "bin", "python")
}
```

This is the core “object-like” pattern in Go.
You put data in a struct and related behavior in methods.

### Pointers

You do not need pointer-heavy code for the first pass, but you should know the basic idea.

If a method needs to modify a struct, it usually uses a pointer receiver:

```go
func (m *Manager) SetRoot(path string) {
    m.RootDir = path
}
```

For read-only methods, value receivers are often fine:

```go
func (m Manager) LocalEnvPath() string {
    return filepath.Join(m.RootDir, ".venv")
}
```

Do not overthink this early.
For `rpy`, a practical default is:

- use value receivers for small, read-only structs
- use pointer receivers for managers or things that change internal state

### Interfaces

An interface describes behavior, not data:

```go
type Backend interface {
    CreateEnv(spec EnvSpec) error
    Install(spec EnvSpec, packages []string) error
}
```

Do not start the project with a pile of interfaces.
Add one when you actually need two interchangeable implementations, such as `venv` and `uv`.

## 6. What Cobra is doing for you

Without Cobra, you would need to manually parse:

- subcommands
- flags
- help text
- argument validation

Cobra handles that boilerplate.

You still need to write the real program logic, but Cobra gives the CLI structure.

The usual flow is:

1. define a command
2. attach it to a parent command
3. implement its `RunE` function
4. call real logic from there

`RunE` is the version you usually want because it can return an error:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    return doThing()
}
```

## 7. Scaffold the CLI

Inside the repo, initialize the module if you have not already:

```bash
go mod init github.com/you/rpy
```

Then scaffold Cobra:

```bash
cobra-cli init --pkg-name github.com/you/rpy
cobra-cli add doctor
cobra-cli add env
cobra-cli add run
cobra-cli add pip
cobra-cli add create -p envCmd
cobra-cli add path -p envCmd
cobra-cli add info -p envCmd
cobra-cli add install -p pipCmd
```

This should generate a `cmd/` tree and the root command wiring.

After that, try:

```bash
go run . --help
```

If that works, you already have a real Go CLI program, even if the commands do not do much yet.

## 8. Keep the command files thin

This is one of the most important design rules in the project.

Bad pattern:

- the Cobra command opens files
- resolves paths
- runs subprocesses
- prints business logic output
- contains 100 lines of implementation detail

Better pattern:

- the Cobra command reads arguments and flags
- it calls a method or function in `internal/`
- the real work happens there

For example:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    mgr, err := env.NewManager()
    if err != nil {
        return err
    }

    return mgr.CreateLocalEnv("python")
}
```

That is a much better direction than embedding `exec.Command(...)` directly inside every Cobra file.

## 9. Start with a few core structs

You do not need a big architecture.
You need a few sensible types.

### `EnvSpec`

Use this to describe one environment:

```go
type EnvSpec struct {
    Name    string
    Mode    string
    Path    string
    Python  string
    Backend string
}
```

Suggested meaning:

- `Name`: friendly name like `default` or `glass`
- `Mode`: `"local"` or `"central"`
- `Path`: full filesystem path to the env
- `Python`: Python executable used to create it
- `Backend`: `"venv"` for now

### `Manager`

Use a manager type to centralize path resolution and environment logic:

```go
type Manager struct {
    RootDir string
    HomeDir string
}
```

Suggested methods:

- `LocalEnvPath() string`
- `LocalEnvExists() bool`
- `CreateLocalEnv(python string) error`
- `ResolveActiveEnv() (EnvSpec, error)`

This is a good place to learn “object-like” Go.
The manager has state, and the methods operate on that state.

### `DoctorResult`

If you want `doctor` to be clean, give its checks structure:

```go
type DoctorResult struct {
    HomeDirOK   bool
    RpyHomeOK   bool
    PythonFound bool
    PythonPath  string
}
```

You can compute that in code, then print it.
This is easier to reason about than mixing checks and printing in random order.

## 10. Learn the standard library pieces you need

Most of this project can be built from the Go standard library.

### `os`

Use it for:

- current working directory
- environment variables
- checking files
- creating directories

Examples:

```go
cwd, err := os.Getwd()
home, err := os.UserHomeDir()
err = os.MkdirAll(path, 0o755)
```

### `path/filepath`

Use this instead of string-concatenating paths.

```go
venvPath := filepath.Join(cwd, ".venv")
```

This matters for portability and general correctness.

### `os/exec`

This is how you run Python, pip, and other commands.

Example:

```go
cmd := exec.Command("python", "-m", "venv", ".venv")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Stdin = os.Stdin
err := cmd.Run()
```

This package is central to `rpy`.
You will use it for:

- creating environments
- running commands inside environments
- wrapping `pip`
- checking whether Python exists

### `errors`

Use plain errors with clear messages.

```go
return errors.New("no local environment found")
```

You do not need elaborate error hierarchies yet.

### `fmt`

Use it for formatting output:

```go
fmt.Println("Path:", env.Path)
fmt.Printf("Python: %s\n", env.Python)
```

## 11. Build the first command: `doctor`

This is the best first command because it touches several useful basics without much complexity.

What `doctor` should check:

- can the program find the user home directory?
- can `~/.rpy` be created?
- is `python` available on `PATH`?

A good implementation path:

1. get the home directory with `os.UserHomeDir()`
2. build the `~/.rpy` path with `filepath.Join`
3. call `os.MkdirAll`
4. locate Python with `exec.LookPath("python")`
5. print the results

What you learn from this:

- standard library path handling
- directory creation
- external program lookup
- returning errors properly

## 12. Build `env create`

This is the first command that does something real.

Its job is simple:

- find the current directory
- decide that the env path is `./.venv`
- run `python -m venv .venv`

The core subprocess code will look roughly like this:

```go
cmd := exec.Command(python, "-m", "venv", ".venv")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Stdin = os.Stdin
return cmd.Run()
```

A better version will run in the project root and use the full path:

```go
cmd := exec.Command(python, "-m", "venv", envPath)
```

What matters conceptually:

- the Go program is not creating the venv by itself
- it is orchestrating Python correctly

That is a recurring pattern in `rpy`.

## 13. Build `env path`

This command is intentionally boring.
That is a good thing.

Its only job is to print the resolved environment path.
At first, that can simply be:

```text
<current-project>/.venv
```

This command is useful because it forces you to separate “how do I resolve the environment path?” from “how do I create it?”.

That logic belongs in your `Manager`.

## 14. Build `env info`

This should report basic information, not everything.

For the first version, show:

- mode
- path
- whether the env exists
- path to the env Python executable

Example:

```text
Mode: local
Path: /home/user/project/.venv
Exists: yes
Python: /home/user/project/.venv/bin/python
```

This is a good command for practicing small structs and small helper methods.

For example, you can have:

- one method that resolves the env spec
- one method that computes the Python path
- one method that checks whether the env exists

## 15. Build `run`

This is the most important command in the first version.

The whole point is to avoid needing shell activation.
Instead of:

```bash
source .venv/bin/activate
python train.py
```

you want:

```bash
rpy run python train.py
```

The implementation idea is:

1. resolve the active environment
2. compute its `bin` directory
3. prepend that directory to `PATH`
4. run the requested command

On Unix-like systems the env binaries live in:

```text
<env>/bin
```

On Windows they live in:

```text
<env>/Scripts
```

In Go, you can pass environment variables to a child process by setting `cmd.Env`.

Conceptually:

```go
cmd := exec.Command(args[0], args[1:]...)
cmd.Env = os.Environ()
```

Then replace or prepend `PATH`.

You will need to learn:

- how `args []string` works
- how to split command name from command arguments
- how to forward stdin/stdout/stderr
- how to preserve exit codes

Preserving exit codes matters.
If `pytest` fails inside `rpy run`, `rpy` should also exit with failure.

## 16. Build `pip install`

Do not treat `pip` as a special internal subsystem yet.
Just call the environment’s Python:

```bash
<env-python> -m pip install numpy
```

That keeps behavior explicit and predictable.

In Go, that means:

```go
cmd := exec.Command(envPython, "-m", "pip", "install", "numpy")
```

Later, when you add `pip list` and `pip freeze`, you can reuse the same pattern.

This is a good place to learn how to pass through CLI arguments.
If the user runs:

```bash
rpy pip install numpy pandas
```

then your Go code should take `args` and append them after:

```text
-m pip install
```

## 17. A practical implementation order

Follow this order and do not skip around:

1. make `go run . --help` work
2. make `rpy doctor` work
3. make `rpy env create` work
4. make `rpy env path` work
5. make `rpy env info` work
6. make `rpy run python --version` work
7. make `rpy pip install numpy` work

This order is good because each step builds on the last one without forcing you into premature design.

## 18. What code should probably live where

Use this as a rough guide.

In `cmd/`:

- define commands
- define flags
- validate simple user input
- call internal logic

In `internal/env/`:

- resolve environment paths
- create local envs
- compute env Python paths
- check whether envs exist

In `internal/pyexec/`:

- run subprocesses
- set up `PATH`
- forward standard streams
- preserve exit codes

In `internal/app/` or `internal/doctor/`:

- machine checks
- `~/.rpy` setup
- Python lookup

This split is not sacred.
It is just a clean place to start.

## 19. A few Go habits worth learning early

### Return errors instead of printing everywhere

This is better:

```go
func CreateLocalEnv() error {
    // do work
    return nil
}
```

than this:

```go
func CreateLocalEnv() {
    fmt.Println("something failed")
}
```

Let the command decide how to present the error.

### Keep functions small

If a function starts doing path resolution, environment checks, subprocess setup, and output formatting all in one place, split it.

### Prefer obvious code over clever abstractions

You are learning Go and building an early-stage CLI.
A straightforward 20-line function is better than a fancy abstraction you do not understand yet.

### Use the standard library first

Do not add dependencies for things Go already does well.

For this phase, the standard library plus Cobra is enough.

## 20. What not to worry about yet

Do not block on these topics:

- generics
- concurrency
- channels
- advanced interfaces
- reflection
- dependency injection frameworks
- perfect project architecture

None of those are required to build the first useful `rpy`.

## 21. What “object-like programming” looks like in Go

Since you asked about structs and object-like code, here is the practical answer.

In Python you might do something like this:

```python
class Manager:
    def __init__(self, root_dir):
        self.root_dir = root_dir

    def local_env_path(self):
        return self.root_dir + "/.venv"
```

In Go, the equivalent shape is:

```go
type Manager struct {
    RootDir string
}

func (m Manager) LocalEnvPath() string {
    return filepath.Join(m.RootDir, ".venv")
}
```

That is the basic pattern.

- `struct` holds data
- method holds behavior
- you create an instance and call methods on it

Example usage:

```go
mgr := Manager{RootDir: "/tmp/project"}
fmt.Println(mgr.LocalEnvPath())
```

You do not need inheritance to write clean code here.
Just keep related data and behavior together.

## 22. When you are stuck

If you get lost, reduce the problem.

Instead of “implement `rpy run`,” ask:

1. can I get the current directory?
2. can I compute `.venv/bin`?
3. can I run `python --version` from Go?
4. can I pass a custom `PATH` to the subprocess?

Most CLI work is just small steps chained together.

## 23. What success looks like

You are done with the first pass when this feels normal:

```bash
go run . doctor
go run . env create
go run . env path
go run . run python --version
go run . pip install numpy
```

At that point, you will not know all of Go, but you will know enough Go to keep moving.
That is the right target.

## 24. After this guide

Once the basics work, the next reasonable features are:

- `rpy.toml` project config
- central environments
- Python discovery
- optional `uv` backend

That is the point where handing the CLI back for completion makes sense, because the hard part then becomes design and integration rather than “what is a Go struct?”.
