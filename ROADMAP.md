# rpy Roadmap

This roadmap is intentionally small. It starts with the first useful CLI and
stops before the more advanced environment-manager features.

## Phase 0: Cobra skeleton

Goal: a working `rpy` binary with help text and a few commands.

Commands:

- `rpy --help`
- `rpy doctor`

Tasks:

- initialize the module
- scaffold the CLI with `cobra-cli`
- add the root command
- add `doctor`

Done when:

- the binary builds
- help text works
- `doctor` prints useful status checks

## Phase 1: Local env creation

Goal: create a project-local virtual environment.

Commands:

- `rpy env create`
- `rpy env path`
- `rpy env info`

Implementation:

- create `.venv` with `python -m venv .venv`
- resolve the env path
- report basic metadata about the env

Done when:

- a project directory gets `.venv/`
- `env path` prints the expected path
- `env info` confirms the env exists

## Phase 2: Run inside the env

Goal: run commands without shell activation.

Commands:

- `rpy run python --version`
- `rpy run pip --version`

Implementation:

- detect the active environment
- prepend its `bin` or `Scripts` directory to `PATH`
- forward stdin/stdout/stderr
- preserve the child process exit code

Done when:

- commands run in the venv
- failures are passed through correctly

## Phase 3: Pip wrapper

Goal: install packages through the selected env.

Commands:

- `rpy pip install numpy`
- `rpy pip list`
- `rpy pip freeze`

Implementation:

- invoke `<env-python> -m pip ...`
- keep CLI flags small and explicit

Done when:

- packages install into the chosen environment
- package listing works

## Phase 4: Project config

Goal: store project choices in `rpy.toml`.

Commands:

- `rpy init`
- `rpy env use .venv`

Stored values:

- project name
- env mode
- env name or path
- Python version preference
- backend choice

Done when:

- the project config can be written and read
- commands can resolve settings from it

## Phase 5: Central envs

Goal: support named reusable environments under `~/.rpy/envs/`.

Commands:

- `rpy env create glass --central`
- `rpy env list`
- `rpy env use glass`
- `rpy env remove glass`

Done when:

- named envs are created in the shared directory
- projects can point at them

## Phase 6: Python discovery

Goal: find available Python interpreters.

Commands:

- `rpy py list`
- `rpy py which`
- `rpy py use 3.12`

Sources:

- `PATH`
- system install locations
- `pyenv`
- `~/.rpy/pythons/`

Done when:

- `rpy` can report usable Python binaries

## Phase 7: Optional uv backend

Goal: make `uv` an optional backend, not the default assumption.

Commands:

- `rpy env create glass --central --backend uv`
- `rpy pip install numpy --backend uv`
- `rpy sync`

Done when:

- `uv` is only used when requested
- `doctor` can report whether `uv` exists

## Phase 8: Python installation

Goal: install Python versions into `~/.rpy/pythons/`.

This is a later phase because it adds more platform-specific complexity.

Done when:

- `rpy` can manage Python versions without relying on the system install alone

