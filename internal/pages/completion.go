package pages

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Completion = cli.Command{
	Name:  "completion",
	Usage: "生成自动补全脚本",
	Subcommands: []*cli.Command{
		{
			Name:  "bash",
			Usage: "生成Bash自动补全脚本",
			Action: func(_ *cli.Context) error {
				fmt.Println(`#! /bin/bash

: ${PROG:="MCST"}

# Macs have bash3 for which the bash-completion package doesn't include
# _init_completion. This is a minimal version of that function.
_cli_init_completion() {
  COMPREPLY=()
  _get_comp_words_by_ref "$@" cur prev words cword
}

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base words
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if declare -F _init_completion >/dev/null 2>&1; then
      _init_completion -n "=:" || return
    else
      _cli_init_completion -n "=:" || return
    fi
    words=("${words[@]:0:$cword}")
    if [[ "$cur" == "-"* ]]; then
      requestComp="${words[*]} ${cur} --generate-bash-completion"
    else
      requestComp="${words[*]} --generate-bash-completion"
    fi
    opts=$(eval "${requestComp}" 2>/dev/null)
    COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete $PROG
unset PROG`)
				return nil
			},
		},
		{
			Name:  "zsh",
			Usage: "生成ZShell自动补全脚本",
			Action: func(_ *cli.Context) error {
				fmt.Println(`#compdef MCST

_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

compdef _cli_zsh_autocomplete MCST`)
				return nil
			},
		},
		{
			Name:  "powershell",
			Usage: "生成PowerShell自动补全脚本",
			Action: func(_ *cli.Context) error {
				fmt.Println(`$name = "MCST"
Register-ArgumentCompleter -Native -CommandName $name -ScriptBlock {
    param($commandName, $wordToComplete, $cursorPosition)
    $other = "$wordToComplete --generate-bash-completion"
    Invoke-Expression $other | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}`)
				return nil
			},
		},
	},
}
