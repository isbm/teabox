#!/usr/bin/bash

SOCK="/tmp/example-callback.sock"

#
# Set status of the logging output
#
function set_status() {
    m=$1
    $(echo "LOGGER-STATUS::$m" | nc -w0 -U $SOCK)
    if [[ "$?" == "1" ]]; then
	exit 1
    fi
}

#
# Set title of the logging output widget
#
function set_title() {
    m=$1
    $(echo "LOGGER-TITLE::$m" | nc -w0 -U $SOCK)
    if [[ "$?" == "1" ]]; then
	exit 1
    fi
}

function load_form() {
    echo "Loading form"
    sleep 0.3
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
