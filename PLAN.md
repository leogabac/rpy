# PLAN.md — `rpy`

`rpy` is a small Python environment manager written in Go.

The goal is **not** to replace Conda, Pyenv, Pip, Virtualenv, or uv. The goal is to build a tiny, understandable tool that wraps the common workflows I actually use:

* manage Python versions
* create and manage virtual environments
* run commands inside environments
* install packages through `pip` or, optionally, `uv`
* support centralized environments, similar in spirit to Conda environments

The project is primarily a learning project for Go, CLI design, filesystem layout, subprocess management, and Python tooling internals.

---

## 1. Core Philosophy

`rpy` should be:

* small
* understandable
* boring
* easy to debug
* explicit about what it wraps
* useful even if it is worse than existing tools

It should avoid trying to solve hard package-management problems too early.

In particular, `rpy` should **not** initially implement:

* dependency resolution
* its own package index
* its own lockfile format
* binary package management
* Conda-style non-Python dependency solving
* cross-platform Python build management from scratch

Instead, it should compose existing tools well.

---

## 2. High-Level Identity

Name:

```text
rpy
```

Meaning:

```text
Reiko Python
```

Intended usage:

```bash
rpy py list
rpy py use 3.12
rpy env create
rpy env list
rpy env activate myenv
rpy pip install numpy
rpy run python script.py
```

---

## 3. Main Use Cases

### 3.1 Project-local virtual environments

This is the simplest and most important workflow.

Example:

```bash
cd my-project
rpy env create
rpy pip install numpy torch
rpy run python train.py
```

By default, this creates:

```text
my-project/.venv/
```

and optionally writes:

```text
my-project/rpy.toml
```

with project metadata.

---

### 3.2 Centralized environments

This is the Conda-like use case.

Instead of creating `.venv` inside each project, `rpy` can create named environments in a central directory:

```bash
rpy env create glass --python 3.12 --central
rpy env use glass
rpy run python train.py
```

Central environments live under:

```text
~/.rpy/envs/glass/
~/.rpy/envs/rl/
~/.rpy/envs/dev/
```

This allows reuse across projects, similar to:

```bash
conda activate glass
```

but implemented with regular Python virtual environments.

---

### 3.3 Python version selection

`rpy` should eventually manage or discover Python versions.

Early versions can simply detect existing Python binaries:

```bash
rpy py list
rpy py which
rpy py use /usr/bin/python3.12
```

Later versions may support installing Python versions into:

```text
~/.rpy/pythons/3.12.3/
~/.rpy/pythons/3.11.9/
```

Potential implementation options:

1. First: detect system Python and pyenv Python installations.
2. Later: download prebuilt Python distributions where available.
3. Much later: build CPython from source.

---

### 3.4 Optional `uv` backend

At some point, `rpy` may use `uv` under the hood for a specific high-speed workflow.

This should be optional and explicit.

Example:

```bash
rpy env create glass --backend uv
rpy pip install numpy --backend uv
rpy sync --backend uv
```

or configured in `rpy.toml`:

```toml
[project]
name = "glass-flow"

[env]
backend = "uv"
python = "3.12"
central = true
name = "glass"
```

Important design rule:

> `rpy` should not become a thin alias for `uv` by accident.

The default backend should remain understandable and based on standard Python tools:

```bash
python -m venv
python -m pip
```

`uv` should be used only where it gives clear value:

* faster installs
* centralized environments
* syncing from lockfiles
* isolated tool execution

---

## 4. Non-Goals

`rpy` should not try to compete with Conda early.

Non-goals for the first versions:

* no Conda package solver
* no CUDA package solving
* no system library management
* no replacement for apt/pacman/brew
* no custom PyPI resolver
* no custom wheel builder
* no multi-language environment solving

If non-Python dependencies are needed, `rpy` can later provide helper notes or hooks, but not solve them itself.

---

## 5. Storage Layout

Default home directory:

```text
~/.rpy/
```

Suggested structure:

```text
~/.rpy/
  config.toml
  pythons/
    3.12.3/
    3.11.9/
  envs/
    glass/
    rl/
    dev/
  cache/
    downloads/
    builds/
  shims/
    python
    pip
  logs/
```

Project-local config:

```text
project-root/rpy.toml
```

Example:

```toml
[project]
name = "glass-flow"

[python]
version = "3.12"

[env]
name = "glass"
mode = "central" # local | central
backend = "venv" # venv | uv
```

---

## 6. Environment Modes

`rpy` should support two environment modes.

### 6.1 Local mode

Environment path:

```text
./.venv/
```

Good for:

* project isolation
* reproducibility
* simple workflows
* editor integration

Command:

```bash
rpy env create
```

---

### 6.2 Central mode

Environment path:

```text
~/.rpy/envs/<name>/
```

Good for:

* Conda-like usage
* shared environments across multiple projects
* long-lived research environments

Command:

```bash
rpy env create glass --central
```

or:

```bash
rpy env create glass --mode central
```

---

## 7. CLI Design

The command tree should be simple.

```text
rpy
├── py
│   ├── list
│   ├── which
│   ├── use
│   └── install        # later
│
├── env
│   ├── create
│   ├── list
│   ├── remove
│   ├── use
│   ├── info
│   └── path
│
├── pip
│   ├── install
│   ├── uninstall
│   ├── freeze
│   └── list
│
├── run
│
├── sync              # later, especially for uv backend
├── init
└── doctor
```

---

## 8. Command Examples

### Initialize a project

```bash
rpy init
```

Creates:

```text
rpy.toml
```

Possibly asks:

* local or central env?
* Python version?
* backend: venv or uv?

For the first implementation, avoid interactive prompts and prefer flags/defaults.

---

### Create a local environment

```bash
rpy env create
```

Equivalent to:

```bash
python -m venv .venv
```

---

### Create a central environment

```bash
rpy env create glass --central --python 3.12
```

Creates:

```text
~/.rpy/envs/glass/
```

---

### Use a central environment in a project

```bash
rpy env use glass
```

Writes to `rpy.toml`:

```toml
[env]
name = "glass"
mode = "central"
```

---

### Run a command inside the selected environment

```bash
rpy run python train.py
rpy run pytest
rpy run jupyter lab
```

This avoids shell activation as a hard requirement.

Implementation idea:

* resolve active env
* prepend env binary directory to `PATH`
* execute the command

On Unix:

```text
<env>/bin
```

On Windows:

```text
<env>/Scripts
```

---

### Install packages

```bash
rpy pip install numpy scipy
```

Default backend:

```bash
<env-python> -m pip install numpy scipy
```

uv backend:

```bash
uv pip install --python <env-python> numpy scipy
```

or equivalent, depending on the final uv integration design.

---

## 9. Backend Abstraction

`rpy` should eventually have a small backend interface.

Conceptually:

```go
type EnvBackend interface {
    CreateEnv(spec EnvSpec) error
    Install(env Env, packages []string) error
    Uninstall(env Env, packages []string) error
    Freeze(env Env) error
    ListPackages(env Env) error
}
```

Initial backends:

```text
venv backend:
  python -m venv
  python -m pip
```

Future backend:

```text
uv backend:
  uv venv
  uv pip
  uv sync
```

Design rule:

> Keep the backend abstraction small and only add methods when a real command needs them.

---

## 10. Configuration Resolution

When running a command, `rpy` should determine the active environment in this order:

1. explicit flag
2. `rpy.toml` in current project or parent directories
3. local `.venv` in current project
4. global default environment in `~/.rpy/config.toml`
5. error with a helpful message

Example:

```bash
rpy run python script.py --env glass
```

should override project configuration.

---

## 11. Activation Strategy

Shell activation is annoying and shell-specific.

Therefore, the first implementation should prioritize:

```bash
rpy run <command>
```

instead of:

```bash
source .venv/bin/activate
```

Later, `rpy` can support activation helpers:

```bash
rpy activate glass
```

which prints shell code:

```bash
eval "$(rpy activate glass)"
```

Possible shells:

* bash
* zsh
* fish
* powershell

This should be a later feature.

---

## 12. MVP Roadmap

### Phase 0 — CLI skeleton

Goal: working Cobra CLI.

Commands:

```bash
rpy --help
rpy doctor
```

`doctor` checks:

* Go binary works
* home directory can be resolved
* `~/.rpy` can be created
* Python is available on PATH

---

### Phase 1 — Local venv creation

Goal: create `.venv` in the current project.

Commands:

```bash
rpy env create
rpy env path
rpy env info
```

Implementation:

```bash
python -m venv .venv
```

---

### Phase 2 — Run inside env

Goal: avoid activation.

Commands:

```bash
rpy run python --version
rpy run pip --version
```

Implementation:

* locate env
* modify `PATH`
* run subprocess
* forward stdin/stdout/stderr
* preserve exit codes

---

### Phase 3 — Pip wrapper

Goal: use pip through the selected environment.

Commands:

```bash
rpy pip install numpy
rpy pip list
rpy pip freeze
```

Implementation:

```bash
<env-python> -m pip install ...
```

---

### Phase 4 — Project config

Goal: persist project choices.

Commands:

```bash
rpy init
rpy env use .venv
```

File:

```text
rpy.toml
```

Stores:

* Python version preference
* env mode
* env name/path
* backend

---

### Phase 5 — Central envs

Goal: Conda-like named environments.

Commands:

```bash
rpy env create glass --central
rpy env list
rpy env use glass
rpy env remove glass
```

Storage:

```text
~/.rpy/envs/glass/
```

---

### Phase 6 — Python discovery

Goal: discover available Python interpreters.

Commands:

```bash
rpy py list
rpy py which
rpy py use 3.12
```

Sources:

* PATH
* known system locations
* pyenv installations
* `~/.rpy/pythons/`

---

### Phase 7 — Optional uv backend

Goal: allow faster env/package operations when uv is installed.

Commands:

```bash
rpy env create glass --central --backend uv
rpy pip install numpy --backend uv
rpy sync
```

Important:

* `uv` should be optional.
* `rpy doctor` should report whether uv is available.
* Error messages should explain how to install uv, but not require it.

---

### Phase 8 — Python installation

Goal: install Python versions into `~/.rpy/pythons/`.

This is intentionally later because it is more complex.

Possible strategies:

1. Use existing installed Python versions first.
2. Support pyenv integration.
3. Download standalone Python builds if available.
4. Build CPython from source only as an advanced feature.

Commands:

```bash
rpy py install 3.12.3
rpy py remove 3.12.3
```

---

## 13. Implementation Notes

### Language

Use Go.

Recommended libraries:

```text
cobra                    CLI framework
pelletier/go-toml/v2     TOML config
```

Standard library packages likely needed:

```text
os
os/exec
os/user
path/filepath
strings
errors
fmt
runtime
```

---

## 14. Suggested Go Package Layout

```text
rpy/
  go.mod
  main.go
  cmd/
    root.go
    doctor.go
    env.go
    py.go
    pip.go
    run.go
    init.go
  internal/
    config/
      config.go
      project.go
    env/
      env.go
      resolve.go
      paths.go
    backend/
      backend.go
      venv.go
      uv.go
    python/
      discover.go
      version.go
    shell/
      activation.go
    process/
      run.go
    paths/
      paths.go
```

Keep the first version simpler if needed.

Do not over-abstract before the commands exist.

---

## 15. Error Message Style

Errors should be direct and actionable.

Bad:

```text
failed to run command
```

Good:

```text
Could not create virtual environment at .venv.
Tried: python -m venv .venv
Reason: python was not found on PATH.

Try installing Python 3.12 or run:
  rpy py use /path/to/python
```

---

## 16. Testing Strategy

Start with simple tests around pure logic:

* path resolution
* config parsing
* env mode resolution
* Python executable path construction

For subprocess-heavy commands, use integration tests later.

Useful tests:

```text
- local env path resolves to ./venv or ./.venv
- central env path resolves to ~/.rpy/envs/<name>
- project config is found in parent directories
- backend selection falls back to venv
- uv backend errors clearly if uv is not installed
```

---

## 17. Design Warnings

Avoid these traps:

### 17.1 Building a package manager too early

Wrapping pip is fine.
Writing a resolver is not.

---

### 17.2 Over-designing the backend interface

Only create abstractions after at least two real backends need the same behavior.

---

### 17.3 Making shell activation required

Prefer:

```bash
rpy run python script.py
```

because it is simpler and portable.

---

### 17.4 Hiding too much magic

`rpy` should make it obvious whether it is using:

* `python -m venv`
* `python -m pip`
* `uv venv`
* `uv pip`

Debugging should be easy.

---

## 18. First Feature to Implement

The first useful feature should be:

```bash
rpy env create
```

Behavior:

1. Resolve current working directory.
2. Check whether `.venv` already exists.
3. Find `python` or `python3` on PATH.
4. Run:

```bash
python -m venv .venv
```

5. Print the created environment path.
6. Suggest:

```bash
rpy run python --version
```

This teaches:

* Cobra commands
* subprocesses in Go
* filesystem checks
* useful error handling
* basic Python environment mechanics

---

## 19. Long-Term Dream

A mature `rpy` could eventually support:

```bash
rpy init
rpy py install 3.12.3
rpy env create glass --central --backend uv
rpy pip install torch numpy scipy
rpy run python train.py
rpy sync
rpy export requirements.txt
```

But the first version should remain tiny.

The real win is not replacing existing tools.

The real win is understanding how these tools work by building a small one.
