package env

import (
	"fmt"
	"strings"
)

func PrintShellInit(shell string) error {
	script, err := ShellInitScript(shell)
	if err != nil {
		return err
	}

	fmt.Println(script)
	return nil
}

func PrintShellInstallHint(shell string) error {
	hint, err := ShellInstallHint(shell)
	if err != nil {
		return err
	}

	fmt.Println(hint)
	return nil
}

func ShellInitScript(shell string) (string, error) {
	switch normalizeShell(shell) {
	case "sh":
		return bashShellInit(), nil
	case "zsh":
		return zshShellInit(), nil
	case "fish":
		return fishShellInit(), nil
	default:
		return "", fmt.Errorf("unsupported shell %q; supported shells are bash, zsh, and fish", shell)
	}
}

func bashShellInit() string {
	return strings.TrimSpace(`
# rpy shell integration
rpy() {
  if [ "$#" -ge 3 ] && [ "$1" = "env" ] && [ "$2" = "use" ]; then
    local _rpy_use_activate=""
    local _rpy_use_args=()
    local arg
    for arg in "${@:3}"; do
      if [ "$arg" = "--activate" ]; then
        _rpy_use_activate=1
        continue
      fi
      _rpy_use_args+=("$arg")
    done
    if [ -n "$_rpy_use_activate" ]; then
      command rpy env use "${_rpy_use_args[@]}" || return $?
      local _rpy_activate_path
      _rpy_activate_path="$(command rpy env activate --shell bash --path)" || return $?
      if [ -n "${VIRTUAL_ENV:-}" ] && [ "${VIRTUAL_ENV}" != "$(dirname "$(dirname "$_rpy_activate_path")")" ]; then
        deactivate >/dev/null 2>&1 || true
      fi
      . "$_rpy_activate_path"
      export RPY_AUTO_VENV="$(dirname "$(dirname "$_rpy_activate_path")")"
      return $?
    fi
  fi

  if [ "$#" -ge 2 ] && [ "$1" = "env" ] && [ "$2" = "activate" ]; then
    shift 2
    local _rpy_activate_path
    _rpy_activate_path="$(command rpy env activate --shell bash --path "$@")" || return $?
    if [ -n "${VIRTUAL_ENV:-}" ] && [ "${VIRTUAL_ENV}" != "$(dirname "$(dirname "$_rpy_activate_path")")" ]; then
      deactivate >/dev/null 2>&1 || true
    fi
    . "$_rpy_activate_path"
    export RPY_AUTO_VENV="$(dirname "$(dirname "$_rpy_activate_path")")"
    return $?
  fi

  command rpy "$@"
}

_rpy_autoenv() {
  local current="$PWD/.venv"

  if [ -n "${RPY_AUTO_VENV:-}" ] && [ "${VIRTUAL_ENV:-}" = "$RPY_AUTO_VENV" ] && [ "$current" != "$RPY_AUTO_VENV" ]; then
    deactivate >/dev/null 2>&1 || true
    unset RPY_AUTO_VENV
  fi

  if [ -z "${VIRTUAL_ENV:-}" ] && [ -f "$current/bin/activate" ]; then
    . "$current/bin/activate"
    export RPY_AUTO_VENV="$current"
  fi
}

case ";${PROMPT_COMMAND:-};" in
  *";_rpy_autoenv;"*) ;;
  *) PROMPT_COMMAND="_rpy_autoenv${PROMPT_COMMAND:+;$PROMPT_COMMAND}" ;;
esac

_rpy_autoenv
`)
}

func fishShellInit() string {
	return strings.TrimSpace(`
# rpy shell integration
function rpy
    if test (count $argv) -ge 3; and test "$argv[1]" = "env"; and test "$argv[2]" = "use"
        set -l use_activate 0
        set -l passthrough
        for arg in $argv[3..-1]
            if test "$arg" = "--activate"
                set use_activate 1
                continue
            end
            set passthrough $passthrough $arg
        end
        if test $use_activate -eq 1
            command rpy env use $passthrough
            or return $status
            set -l activate_path (command rpy env activate --shell fish --path)
            or return $status
            set -l activate_root (dirname (dirname "$activate_path"))
            if set -q VIRTUAL_ENV; and test "$VIRTUAL_ENV" != "$activate_root"
                deactivate >/dev/null 2>/dev/null
            end
            source "$activate_path"
            set -gx RPY_AUTO_VENV "$activate_root"
            return $status
        end
    end

    if test (count $argv) -ge 2; and test "$argv[1]" = "env"; and test "$argv[2]" = "activate"
        set -l passthrough $argv[3..-1]
        set -l activate_path (command rpy env activate --shell fish --path $passthrough)
        or return $status
        set -l activate_root (dirname (dirname "$activate_path"))
        if set -q VIRTUAL_ENV; and test "$VIRTUAL_ENV" != "$activate_root"
            deactivate >/dev/null 2>/dev/null
        end
        source "$activate_path"
        set -gx RPY_AUTO_VENV "$activate_root"
        return $status
    end

    command rpy $argv
end

function __rpy_autoenv --on-variable PWD
    set -l current "$PWD/.venv"

    if set -q RPY_AUTO_VENV; and test "$VIRTUAL_ENV" = "$RPY_AUTO_VENV"; and test "$current" != "$RPY_AUTO_VENV"
        deactivate >/dev/null 2>/dev/null
        set -e RPY_AUTO_VENV
    end

    if not set -q VIRTUAL_ENV; and test -f "$current/bin/activate.fish"
        source "$current/bin/activate.fish"
        set -gx RPY_AUTO_VENV "$current"
    end
end

__rpy_autoenv
`)
}

func zshShellInit() string {
	return strings.TrimSpace(`
# rpy shell integration
function rpy() {
  if [ "$#" -ge 3 ] && [ "$1" = "env" ] && [ "$2" = "use" ]; then
    local _rpy_use_activate=""
    local _rpy_use_args=()
    local arg
    for arg in "${@:3}"; do
      if [ "$arg" = "--activate" ]; then
        _rpy_use_activate=1
        continue
      fi
      _rpy_use_args+=("$arg")
    done
    if [ -n "$_rpy_use_activate" ]; then
      command rpy env use "${_rpy_use_args[@]}" || return $?
      local _rpy_activate_path
      _rpy_activate_path="$(command rpy env activate --shell zsh --path)" || return $?
      if [ -n "${VIRTUAL_ENV:-}" ] && [ "${VIRTUAL_ENV}" != "$(dirname "$(dirname "$_rpy_activate_path")")" ]; then
        deactivate >/dev/null 2>&1 || true
      fi
      . "$_rpy_activate_path"
      export RPY_AUTO_VENV="$(dirname "$(dirname "$_rpy_activate_path")")"
      return $?
    fi
  fi

  if [ "$#" -ge 2 ] && [ "$1" = "env" ] && [ "$2" = "activate" ]; then
    shift 2
    local _rpy_activate_path
    _rpy_activate_path="$(command rpy env activate --shell zsh --path "$@")" || return $?
    if [ -n "${VIRTUAL_ENV:-}" ] && [ "${VIRTUAL_ENV}" != "$(dirname "$(dirname "$_rpy_activate_path")")" ]; then
      deactivate >/dev/null 2>&1 || true
    fi
    . "$_rpy_activate_path"
    export RPY_AUTO_VENV="$(dirname "$(dirname "$_rpy_activate_path")")"
    return $?
  fi

  command rpy "$@"
}

_rpy_autoenv() {
  local current="$PWD/.venv"

  if [ -n "${RPY_AUTO_VENV:-}" ] && [ "${VIRTUAL_ENV:-}" = "$RPY_AUTO_VENV" ] && [ "$current" != "$RPY_AUTO_VENV" ]; then
    deactivate >/dev/null 2>&1 || true
    unset RPY_AUTO_VENV
  fi

  if [ -z "${VIRTUAL_ENV:-}" ] && [ -f "$current/bin/activate" ]; then
    . "$current/bin/activate"
    export RPY_AUTO_VENV="$current"
  fi
}

autoload -Uz add-zsh-hook
if ! (( ${chpwd_functions[(I)_rpy_autoenv]} )); then
  add-zsh-hook chpwd _rpy_autoenv
fi

_rpy_autoenv
`)
}

func ShellInstallPath(shell string) (string, error) {
	switch normalizeShell(shell) {
	case "sh":
		return "~/.bashrc", nil
	case "zsh":
		return "~/.zshrc", nil
	case "fish":
		return "~/.config/fish/config.fish", nil
	default:
		return "", fmt.Errorf("unsupported shell %q; supported shells are bash, zsh, and fish", shell)
	}
}

func ShellInstallLine(shell string) (string, error) {
	switch normalizeShell(shell) {
	case "sh":
		return `eval "$(rpy shell-init --shell bash)"`, nil
	case "zsh":
		return `eval "$(rpy shell-init --shell zsh)"`, nil
	case "fish":
		return `rpy shell-init --shell fish | source`, nil
	default:
		return "", fmt.Errorf("unsupported shell %q; supported shells are bash, zsh, and fish", shell)
	}
}

func ShellInstallHint(shell string) (string, error) {
	path, err := ShellInstallPath(shell)
	if err != nil {
		return "", err
	}

	line, err := ShellInstallLine(shell)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Add this line to %s:\n\n%s", path, line), nil
}
