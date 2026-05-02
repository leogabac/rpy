# rpy

`rpy` is a small Python environment manager written in Go.

I am not trying to build the next Conda, or replace pyenv, pip, virtualenv, or uv.
I am trying to glue some of my common Python workflows.

And I wanted to learn Go.

## Why

This project came out of a pretty ordinary frustration: my Python setup had too many moving parts.
For a long time I used `pyenv` for Python versions, `venv` or `virtualenv` for environments, and `pip` for packages.
That works, but it means one workflow is split across several tools, each with its own rules and quirks.

`conda` solves more of that in one place, but it is heavier than I want for day-to-day work.
It is still useful on HPC systems, where you often do not get to choose the machine.
`uv` is close to ideal in speed and ergonomics, but it does not yet give me the shared-environment pattern I want for some research projects.

So `rpy` is my attempt to build a smaller tool around the pieces I use.

The project is deliberately Linux-first. Managed Python installs are built from upstream CPython source instead of relying on interpreter discovery or bundled system dependency management.

## Current MVP

The current CLI is intentionally narrow:

- `rpy env create`
- `rpy env create glass --use`
- `rpy env use local`
- `rpy env list`
- `rpy env path`
- `rpy env info`
- `rpy env activate`
- `rpy py install 3.12.3`
- `rpy py list`
- `rpy py which 3.12`
- `rpy run ...`
- `rpy pip ...`

`rpy py install` downloads `Python-<version>.tgz` from `python.org`, runs `./configure`, `make`, and `make install`, and places the result under `~/.rpy/pythons/<version>/`.

Build dependencies are not installed for you. That is intentional. Set your system up the way you prefer, then let `rpy` compile against it.

On Debian or Ubuntu you will usually want something close to:

```sh
sudo apt install build-essential libssl-dev zlib1g-dev \
  libbz2-dev libreadline-dev libsqlite3-dev libffi-dev \
  liblzma-dev tk-dev uuid-dev
```

If configure or build steps fail, install the missing system libraries yourself and retry, similar to the `pyenv` workflow.

Activation works by printing a shell snippet because a subprocess cannot mutate
your current shell session directly:

```sh
eval "$(rpy env activate)"
```

For a better interactive workflow, load the shell integration once in your
shell startup:

```sh
eval "$(rpy shell-init --shell bash)"
```

For fish:

```fish
rpy shell-init --shell fish | source
```

That integration does two things:

- makes `rpy env activate` work in the current shell
- auto-activates `.venv` when you enter a project directory

## Environment Model

Each project uses exactly one environment at a time:

- `local`: `./.venv`
- `shared:<name>`: `~/.rpy/envs/<name>`

Create a local env:

```sh
rpy env create
```

Create and select a shared env:

```sh
rpy env create glass --python 3.12.0 --use
```

Switch back to the local env:

```sh
rpy env use local
```

See what the project can use and what is currently selected:

```sh
rpy env list
rpy env info
```
