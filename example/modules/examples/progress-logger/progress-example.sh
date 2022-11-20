#!/usr/bin/bash

# Load common library
SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SELF_PATH/../../../lib/shell/bash-common.sh

function run() {
	sleep 2
}

# Main
if [[ $# -eq 0 ]]; then
    run $@
else
    while [[ $# -gt 0 ]]; do
	case $1 in
	    -s|--setup)
        echo "There is no setup"
		shift
		;;
	    -h|--help)
        echo "Help you yourself"
		shift
		;;
	    *)
		run $@
		break
		;;
	esac
    done
fi
