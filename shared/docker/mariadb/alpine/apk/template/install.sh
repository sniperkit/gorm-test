#!/bin/bash -

[[ "$(whoami)" != "root" ]] && echo "Usage: sudo $0" && exit 1

cls="\e(B\e[m"
function bd()       { echo "${bold}$1${normal}"; }
function yellow()   { echo -en "\e[38;5;208m$1${cls}"; }
function green()    { echo -en "\e[38;5;112m$1${cls}"; }
function red()      { echo -en "\e[38;5;196m$1${cls}"; }
function error_echo() { (>&2 echo "$(red ERROR): $1"); exit 1; }
function yesno() {
    local _yn; read -p "$1 " _yn
    local _regex='^[Yy]'
    [[ "$_yn" =~ $_regex ]] && echo "yes" || echo "no"
}

function var_substitute() {
    local _file="$1"
    [[ ! -f "$_file" ]] && error_echo "$_file not found"
    shift

    declare -A vars=()

    while [[ $# -gt 1 ]]; do
        val="$2"
        val="${val%\"}"
        vars[$1]="${val#\"}"
        shift
        shift
    done

    local _line
    local _regex='(\$\{([a-zA-Z][a-zA-Z_0-9]*)\})'

    local _old_ifs="$IFS"; IFS=
    while read -r _line; do
        while [[ "$_line" =~ $_regex ]]; do
            local _lhs="${BASH_REMATCH[1]}"
            local _name="${BASH_REMATCH[2]}"
            _line="${_line//$_lhs/\$$_name}"
            if [ ! "${vars[$_name]}" ]; then
                vars[$_name]=""
            fi
        done
    done < "$_file"

    ### get values of undefined
    for i in "${!vars[@]}"; do
        if [ ! "${vars[$i]}" ]; then
            read -p "Value for $i=" vars[$i]
        fi
    done

    local _service_file=""
    while read -r _line; do
        while [[ "$_line" =~ $_regex ]]; do
            local _LHS="${BASH_REMATCH[1]}"
            local _VAR="${BASH_REMATCH[2]}"
            if [ "${vars[$_VAR]}" ]; then
                _line="${_line//$_LHS/${vars[$_VAR]}}"
            else
                _line="${_line//$_LHS/\$-$_VAR-}"
            fi
        done
        _service_file="${_service_file}${_line}\n"
    done < "$_file"

    ### reset vars get list of disks to replace (if external)
    vars=()
    IFS=' '
    while read -r -a _arry; do
        if [ "${_arry[1]}" = "/" ]; then
            _arry[1]=
        fi
        vars[${_arry[0]}]="${_arry[1]}"
    done < <(df --type=ext4 --type=ext3 --output=source,target | grep -e '^/' | perl -e 's{^(\S+)\s+(\S+).*$}{$1 $2}' -p)
    IFS="$_old_ifs"

    for i in "${!vars[@]}"; do
        _service_file=$(echo $_service_file | perl -e "s|${i}|${vars[$i]}|g" -p)
    done

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

function get_service() {
    local _cwd="$1"
    local _service="$2"

    if [ ! "${_service}" ]; then
        local _files="$(ls -1 ${_cwd}/*.service 2>&1 | wc -l)"
        [[ $? -eq 1 ]] && _files=0

        case "$_files" in
            0)
                error_echo "Cannot find any *.service files"
                ;;
            1)
                _service="$(ls -1 *.service)"
                ;;
            *)
                _choices="$(ls -1 *.service | awk '{printf "  %d) %s\n", NR, $0}')"
                _question="$(echo -e "$(green "$_choices")\nWhich service to install? ")"
                read -p "${_question}" num
                num="^\s*$num)"
                _service=$(echo "$_choices" | grep -e $num | sed 's/^[^)]*) //')
                ;;
        esac
    fi

    echo "${_service}"
}

function install_file() {
    local _cwd="$1"
    local _service="$2"
    local _target="/lib/systemd/system/${_service}"
    local _source="${_cwd}/${_service}"

    [[ ! -f "${_source}" ]] && error_echo "No .service file to copy"

    if [[ -f "${_target}" ]] \
    && [[ "$(yesno "File $(yellow $_target) exists, overwrite?")" = "no" ]]; then
        echo $_target;
        return 1
    fi

    local _content=$(var_substitute "${_source}")
    if [ "${_content}" ]; then
        echo "${_content}" > $_target
    fi

    [[ ! -f "${_target}" ]] && error_echo "File ${_target} does not exist"

    echo $_target
}

function enable_service() {
    local _service="$1"
    local _file="$2"

    [[ ! -f "${_file}" ]] && exit 1

    if [[ "$(systemctl is-enabled ${_service} 2>&1)" = "enabled" ]] \
    && [[ "$(yesno "Service enabled. Disable $(yellow "${_service}")?")" = "yes" ]]; then
        systemctl disable ${_file}
    fi

    if [[ "$(yesno "Enable service $(yellow "${_service}")?")" = "yes" ]]; then
        systemctl enable ${_file}
        systemctl daemon-reload
        echo -e "Test service run:\n\$ systemctl start $(green "${_service}")\n"
    elif [[ "$(yesno "Remove file $(yellow "$_file")?")" = "yes" ]]; then
        rm $_file
        echo "File $(red "$_file") removed"
    fi
}

_CWD=$(cd "$(dirname $0)"; pwd)
_SERVICE="$(get_service "${_CWD}" "$1")"
_FILE="$(install_file "${_CWD}" "${_SERVICE}")"

enable_service "${_SERVICE}" "${_FILE}"

exit 0
