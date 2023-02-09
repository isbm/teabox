#!/usr/bin/bash

SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SELF_PATH/../../../lib/shell/bash-common.sh

function setup() {
    echo "Setup phase"
}

function run() {
    echo "$@"
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
	    --do-something-when-checked)
		date >> checked.txt
		break
		;;
	    --do-something-when-is-not-checked)
		date >> unchecked.txt
		break
		;;
            *)
		run $@
		break
		;;
        esac
    done
fi
