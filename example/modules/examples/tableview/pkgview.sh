#!/usr/bin/bash

SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SELF_PATH/../../../lib/shell/bash-common.sh

function run() {
    PKG=$1

    # Set the title
    api logger.title "Installing $PKG as your default editor :)"
	echo "You choose $PKG"

}

if [[ $# -eq 0 ]]; then
    run "emacs-anyways.deb"
else
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--setup)
		break
		;;
	    --pkgname=*)
		run "${1#*=}"
		break
		;;
            *)
		run "emacs-default.deb"
		break
		;;
        esac
    done
fi
