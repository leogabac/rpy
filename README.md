# rpy

`rpy` is a small Python environment manager written in Go.

I am trying to glue together the subset of Python workflow I use. And I wanted to learn Go.

> [!NOTE]
> `rpy` stands for "Reiko's Python"

## Why

This project came out of a pretty ordinary frustration: my Python setup had too many moving parts.
For a long time I used `pyenv` for Python versions, `venv` or `virtualenv` for environments, and `pip` for packages.
That works, but it means one workflow is split across several tools, each with its own rules and quirks.

`conda` solves more of that in one place, but it is heavier than I want for day-to-day work.
It is still useful on HPC systems, where you often do not get to choose the machine.
`uv` is close to ideal in speed and ergonomics, but it does not yet give me the shared-environment pattern I want for some research projects.

So `rpy` is my attempt to build a smaller tool around the pieces I use.

## What `rpy` does

`rpy` currently manages three things:

1. Python runtimes under `~/.rpy/pythons/<version>`
2. A local project environment at `./.venv`
3. Shared named environments under `~/.rpy/envs/<name>`

Each project uses exactly one environment at a time:

- `local`: `./.venv`
- `shared:<name>`: `~/.rpy/envs/<name>`

The project selection is stored in a small `.rpy-env` file when you choose a shared env.

## Requirements

This tool is deliberately Linux-first.

Managed Python installs are built from upstream CPython source instead of relying on interpreter discovery magic, prebuilt embedded runtimes, or trying to install system packages for you.

That means `rpy` expects you to already have system build dependencies.

## Install

Build the binary from the repo root:

```sh
go build -o rpy .
```

You can then move `rpy` somewhere on your `PATH`, for example:

```sh
install -m 0755 rpy ~/.local/bin/rpy
```

## Shell Integration

`rpy` cannot directly mutate your current shell session, so activation commands are normally printed as shell snippets.

For one-off use:

```sh
eval "$(rpy env activate)"
```

For a usable interactive workflow, load shell integration once in your shell startup.

For `bash`:

```sh
eval "$(rpy shell-init --shell bash)"
```

For `zsh`:

```sh
eval "$(rpy shell-init --shell zsh)"
```

For `fish`:

```fish
rpy shell-init --shell fish | source
```

With shell integration loaded:

- `rpy env activate` activates the current project environment in-place
- `rpy env use <name> --activate` switches and activates in one step
- entering a project directory auto-activates the environment selected for that project

To see the line you should add to your startup file:

```sh
rpy shell-init --shell bash --install
```

## Python Runtime Installation

Install a managed Python runtime:

```sh
rpy py install 3.12.3
```

This downloads `Python-3.12.3.tgz` from `python.org`, then runs:

- `./configure`
- `make`
- `make install`

The installed runtime ends up under:

```sh
~/.rpy/pythons/3.12.3
```

List known interpreters:

```sh
rpy py list
```

Resolve one interpreter:

```sh
rpy py which 3.12.3
rpy py which python3
rpy py which /usr/bin/python3
```

Remove a managed runtime:

```sh
rpy py remove 3.12.3
```

## Local Environment Workflow

Create a local project environment:

```sh
rpy env create
```

Choose the Python explicitly:

```sh
rpy env create --python 3.12.3
rpy env create --python python3
rpy env create --python /usr/bin/python3
```

Recreate an existing local environment:

```sh
rpy env create --force
```

Remove it:

```sh
rpy env remove local
```

By default, if no shared environment is selected, the project uses `./.venv`.

## Shared Environment Workflow

Create a shared environment:

```sh
rpy env create glass
```

Create it with a specific interpreter:

```sh
rpy env create glass --python 3.12.3
```

Create it and immediately select it for the current project:

```sh
rpy env create glass --python 3.12.3 --use
```

Recreate an existing shared environment:

```sh
rpy env create glass --python 3.12.3 --force
```

List the local and shared environments visible to the current project:

```sh
rpy env list
```

Switch the current project to a shared environment:

```sh
rpy env use glass
```

Switch back to the local environment:

```sh
rpy env use local
```

Switch and activate in one step:

```sh
rpy env use glass --activate
rpy env use local --activate
```

`rpy env use <name>` persists the selection for the project by writing `.rpy-env`.
`rpy env use <name> --activate` is different: it is temporary for the current shell session and does not write `.rpy-env`.

Remove a shared environment:

```sh
rpy env remove glass
```

## Inspect the Current Environment

Show what the current project is using:

```sh
rpy env info
```

Print the current environment root path:

```sh
rpy env path
```

Print the activation command:

```sh
rpy env activate
```

Print only the activation script path:

```sh
rpy env activate --path
```

## Run Commands Inside the Current Environment

Run a command with the environment on `PATH`:

```sh
rpy run python --version
rpy run pytest
```

Run `pip` through the current environment interpreter:

```sh
rpy pip install numpy
rpy pip list
rpy pip freeze
```

There are also convenience subcommands:

```sh
rpy pip install requests
rpy pip list
rpy freeze
```

## Example Workflows

### Personal Project With Local `.venv`

```sh
rpy py install 3.12.3
rpy env create --python 3.12.3
rpy env use local --activate
rpy pip install -U pip
rpy pip install -r requirements.txt
rpy run python --version
```

### Shared Research Environment Across Multiple Projects

Create the shared environment once:

```sh
rpy py install 3.12.3
rpy env create research --python 3.12.3
```

Use it inside one project:

```sh
rpy env use research --activate
rpy pip install numpy pandas matplotlib
```

Use it inside another project:

```sh
cd ../other-project
rpy env use research --activate
```

## Doctor

Run the built-in checks:

```sh
rpy doctor
```

It currently checks:

- home directory resolution
- `~/.rpy` availability
- Python on `PATH`
- whether shell integration is detected in the current shell
- whether the current project's selected environment exists
- whether `make` and `cc` are available for Python source builds

## Current Limitations

This is still an MVP.

Notably:

- there is no full dependency sync workflow yet
- there is no lockfile or environment spec format yet
- shared env removal does not scan other projects for stale references
- Python installation is Linux-first and source-build oriented
- build dependencies are documented, not managed

That is intentional for now. The goal is a small workflow tool.
