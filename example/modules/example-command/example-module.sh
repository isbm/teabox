#!/usr/bin/bash

SOCK="/tmp/teabox.sock"

#
# Call socket with an API call
#
function api() {
    cls=$1
    msg=$2
    typ=$3
    $(echo "$cls:$typ:$msg" | nc -w0 -U $SOCK)
    if [[ "$?" == "1" ]]; then
	exit 1
    fi
}


#
# Example how to operate loader progress and status
#
function load_form() {
    api init-reset
    api init-set-status "Loading Example Module..."
    api init-alloc-progress "3" int

    for stat in This That "And That Too"
    do
	sleep 0.1
	api init-set-status "Now Loading $stat..."
	api init-inc-progress

	# Reset fields
	api field-set-by-ord "{0}/etc|/usr/local/etc|/opt/etc"
	api field-set-by-ord "{1}ZX-80 Spectrum|Atari|S/390"
	api field-set-by-ord "{2}Borat Sagdiev"
	api field-set-by-ord "{3}false" bool
	api field-set-by-ord "{4}false" bool
    done

    sleep 0.3
}


#
# Set status of the logging output
#
function set_status() {
    api logger-status "$1"
}

#
# Set title of the logging output widget
#
function set_title() {
    api logger-title "$1"
}


function do_something() {
    location=$1
    echo "Putting something to $location..."
    sleep 1
    echo "Done"
}

function setup() {
    echo "Setup"
}

function run() {
    set_status "Checking the arguments"
    echo $@

    max_sec=5
    left_sec="$max_sec"
    for (( i=1; i<=$max_sec; i++ ))
    do
	set_status "Waiting $left_sec seconds..."
	sleep 1
	let "left_sec-=1"
    done

    set_status "Calling Something"
    do_something $1

    set_status "Calling package manager update"
    sudo apt update

    set_status "Finished!"
}

if [[ $# -eq 0 ]]; then
    run $@
else
    while [[ $# -gt 0 ]]; do
	case $1 in
	    -s|--setup)
		load_form
		shift
		;;
	    -h|--help)
		echo "Run me with --setup or directly"
		shift
		;;
	    *)
		run $@
		break
		;;
	esac
    done
fi
