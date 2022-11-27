#!/usr/bin/bash

SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SELF_PATH/../../../lib/shell/bash-common.sh

function setup() {
    api field.set.by-ord "{0}Narnia"
}

function run() {
    WORLD=$1
    TIMES=3

    # Set the title
    api logger.title "Printing $WORLD $TIMES times"

    for (( i=1; i <= $TIMES; i++ ))
    do
	echo "[black]Hello, [aqua:blue] ${WORLD} [black:lightgrey]!"

	# Set the status bar
	api logger.status "Iteration $i out of $TIMES"

	sleep 1
    done
}

if [[ $# -eq 0 ]]; then
    run "world"
else
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--setup)
		setup
		break
		;;
	    --name=*)
		run "${1#*=}"
		break
		;;
            *)
		run "The Upside Down"
		break
		;;
        esac
    done
fi
