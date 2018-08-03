#!/bin/sh

get_path() {
    echo "$(grep $1 /proc/self/mountinfo | grep -v '/volumes/' | cut -f 4,9 -d" " | sed -e 's/^\([^ ]\+\)[ ]\+\([^ ]\+\)$/\2\1/')"
}

getvols() {
    _localetc="$(get_path $_etc)"
    _locallog="$(get_path $_log)"
    _localdata="$(get_path $_root)"

    if [[ "$_localetc" ]] && [[ "$_locallog" ]] && [[ "$_localdata" ]]; then
        echo "-v $_localetc:$_etc -v $_locallog:$_log -v $_localdata:$_root"
    fi
}

_vols="$(getvols)"
_ports="-p 3306:3306"

is_empty() { if [[ ! -d "$1" ]] || [[ ! "$(ls -A $1)" ]]; then echo "yes"; fi }

if [ ! "${_vols}" ]; then
    echo "ERROR: you need run Docker with the -v parameter (see documentation)"
    exit 1
fi

_run="docker run -t --rm ${_vols} ${_ports} ${_image}"

HELP=`cat <<EOT
Usage: docker run -t --rm ${_vols} ${_ports} ${_image} <command>

   init      - initialize directories if they're empty
   bootstrap - create new database (run -it)
   daemon    - run in non-detached mode
   start     - start mariadb server
   stop      - quick mariadb shutdown

EOT
`

if [[ $# -lt 1 ]] || [[ ! "${_vols}" ]]; then echo "$HELP"; exit 1; fi

hint() {
    local hint="| $* |"
    local stripped="${hint//${bold}}"
    stripped="${stripped//${normal}}"
    local edge=$(echo "$stripped" | sed -e 's/./-/g' -e 's/^./+/' -e 's/.$/+/')
    echo "$edge"
    echo "$hint"
    echo "$edge"
}

_cmd=$1
_host=$2
_datadir=/data

_start="docker run ${_vols} ${_ports} -d ${_image}"

assert_ok() { if [ "$?" = 1 ]; then hint "Abort"; exit 1; fi; }

write_systemd_file() {
    local _name="$1"
    local _map="$2"
    local _port="$3"

    local _service_file="${_etc}/docker-${_name}.service"
    local _script="${_etc}/install-systemd.sh"
    local _template_file="${_datadir}/templ/systemd.service"

    if [ "$(grep ^ID= /etc/os-release)" = 'ID=alpine' ]; then
        apk -q --no-cache add bash
    fi

    echo "$(write_template.sh $_template_file name \""${_name}"\" map \""${_map}"\" port \""${_port}"\" image \""${_image}"\")" \
        > ${_service_file}

    echo "Created ${_service_file}"

    cp ${_datadir}/templ/install.sh ${_script}
    chmod 755 ${_script}
    echo "Created ${_script}"

    if [ "$(grep ^ID= /etc/os-release)" = 'ID=alpine' ]; then
        apk --no-cache -q del --purge bash
    fi
}

run_init() {
    if [ "$(is_empty ${_etc})" ]; then
        cp -R ${_datadir}/etc/* ${_etc}/
        chown -R mysql:mysql $_etc $_root $_log

        write_systemd_file ${HOSTNAME} "${_vols}" "${_ports}"
    fi
}

get_pw() {
    local _pw=$(echo $(openssl rand -base64 128) | sed 's/[[:space:]"]//g')
    _pw=${_pw:0:60}
    echo $_pw
}

get_mem_25pc() {
    local _free=$(free -m | grep -i Mem | sed -e 's/^mem:[[:space:]]\+\([0-9]\+\).*$/\1/i')
    echo $((($_free/25)+$_free))M
}

mysqld_is_running() {
    if [ "$(mysqladmin status 2>&1 | grep -i uptime)" ]; then echo "1"; fi
}

install_client() {
    run_init

    if [ "$(is_empty ${_root})" ]; then
        hint "Geting mysql-client, openssl, tzdata"

        if [ "$(grep ^ID= /etc/os-release)" = 'ID=alpine' ]; then
            apk add -q --no-cache --virtual .deps mysql-client openssl tzdata
        else
            apt-get update -q
            apt-get install -qy mysql-client openssl tzdata
        fi

        hint "Installing default DB"
        mysql_install_db --user=mysql

        assert_ok

        if [ ! "$(mysqld_is_running)" ]; then
            hint "Starting mysql"
            mysqld --user=mysql &

            for _i in $(seq 0 9); do
                if [ "$(mysqld_is_running)" = "1" ]; then
                    break
                fi
                sleep 1
            done
        fi

        local _pw=$(get_pw)
        local _keybuf=$(get_mem_25pc)

        if [[ ! "$_pw" ]] || [[ ! "$_keybuf" ]]; then hint "Abort!"; exit 1; fi

        echo $_pw
        echo $_keybuf

        (echo -e "\nY\n${_pw}\n${_pw}\nY\nY") | mysql_secure_installation -S $_sock

        hint "Creating cfgs"
        echo -e "[mysqld]\nkey_buffer_size=${_keybuf}" > $_etc/conf.d/keybufsiz.cnf
        echo -e "[mysql]\npassword=$_pw\n[client]\npassword=$_pw" > $_etc/conf.d/passwd.cnf

        chown -R mysql:mysql $_etc $_root $_log

        echo "CREATE USER 'root'@'172.17.0.1' IDENTIFIED BY '${_pw}'" | mysql mysql
        echo "GRANT ALL PRIVILEGES ON *.* TO 'root'@'172.17.0.1'" | mysql mysql

        hint "Adding timezone info"
        mysql_tzinfo_to_sql /usr/share/zoneinfo | mysql mysql

        hint "Shutting down mysql"
        mysqladmin shutdown

        hint "DONE!"
        echo "Check install-systemd.sh in etc dir"
    fi
}

case "${_cmd}" in
    init)
        hint "initializing server"
        run_init
        ;;

    bootstrap)
        install_client
        ;;

    start)
        hint "starting server"
        mysqld --user=mysql &
        ;;

    daemon)
        run_init
        mysqld --user=mysql
        ;;

    stop)
        hint "${_cmd} server"
        pkill mysqld
        ;;

    kill)
        killall mysqld
        ;;

    *)
        echo "ERROR: Command '${_cmd}' not recognized"
        ;;
esac

