#!/bin/bash -

[[ $# -lt 3 ]] && echo "$0 <template-file> name \"value\"..." && exit 1

function vars_substitute() {
    local _file="$1"
    local _line
    local _regex='(\$\{([a-zA-Z][a-zA-Z_0-9]*)\})'

    local _old_ifs="$IFS"; IFS=
    while read -r _line; do
        while [[ "$_line" =~ $_regex ]]; do
            local _lhs="${BASH_REMATCH[1]}"
            local _name="${BASH_REMATCH[2]}"
            _line="${_line//$_lhs/\$$_name}"
            if [ ! "${VARS[$_name]}" ]; then
                VARS[$_name]=""
            fi
        done
    done < "$_file"

    local _service_file=""
    while read -r _line; do
        while [[ "$_line" =~ $_regex ]]; do
            local _lhs="${BASH_REMATCH[1]}"
            local _var="${BASH_REMATCH[2]}"
            if [ "${VARS[$_var]}" ]; then
                _line="${_line//$_lhs/${VARS[$_var]}}"
            else
                _line="${_line//$_lhs/\$\-$_var\-}"
            fi
        done
        _service_file="${_service_file}${_line}\n"
    done < "$_file"

    IFS="$_old_ifs"

    local _content=""
    local _regex='(\$-([a-zA-Z][a-zA-Z_0-9]*)-)'
    while read -r _line; do
        while [[ "$_line" =~ $_regex ]]; do
            local _LHS="${BASH_REMATCH[1]}"
            local _VAR="${BASH_REMATCH[2]}"
            _line="${_line//$_LHS/\$\{$_VAR\}}"
        done
        _content="${_content}${_line}\n"
    done < <(echo $_service_file)

    echo -e "$_content"
}

_file="$1"
[[ ! -f "$_file" ]] && echo "$_file not found" && exit 1

shift

declare -A VARS=()

while [[ $# -gt 1 ]]; do
    val="$2"
    val="${val%\"}"
    VARS[$1]="${val#\"}"
    shift
    shift
done

echo "$(vars_substitute $_file)"
