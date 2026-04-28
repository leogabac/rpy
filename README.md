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

## Current MVP

The current CLI is intentionally narrow:

- `rpy env create`
- `rpy env path`
- `rpy env info`
- `rpy env activate`
- `rpy py list`
- `rpy py which 3.12`
- `rpy run ...`
- `rpy pip ...`

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
