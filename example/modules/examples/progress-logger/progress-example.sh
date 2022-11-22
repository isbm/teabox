#!/usr/bin/bash

# Load common library
SELF_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source $SELF_PATH/../../../lib/shell/bash-common.sh

function setup() {
	api common.title "Example Of Potato Preparation"
	# Add a colored text. Because if there is no fun, what's the point?
	api common.info.set "[white]This is an example how [yellow]potato [white]can be [lime]boiled, prepared, served [white]and then served."

	# Setup things that needs to be done
	api common.list.add "{water}Boil water"
	api common.list.add "{potato}Add potatoes"
	api common.list.add "{drain}Drain water"
	api common.list.add "{smash}Smash potatoes"
	api common.list.add "{cream}Add sauer cream"
	api common.list.add "{enjoy}Enjoy!"

	# Set progress
	api common.progress.set 0 int
	api common.progress.allocate 10 int
}

function run() {
	# Complete items
	for id in water potato drain smash cream enjoy; do
		sleep 0.3
		api common.list.complete $id

		api common.progress.event "Updating $id"
		api common.progress.next
	done
	sleep 1

	# Do some silly animation for sake of just a demo
	api common.progress.event "Done!"
	for i in {1..3}
	do
		api common.info.set "[lime]Enjoy!"
		sleep 0.1
		api common.info.set "[white]Enjoy! [lime]Enjoy!"
		sleep 0.1
		api common.info.set "[white]Enjoy! [lime]Enjoy! [yellow]Enjoy!"
		sleep 0.1
		api common.info.set "[white]Enjoy! [lime]Enjoy! [yellow]Enjoy!"
		sleep 0.1
		api common.info.set "[yellow]Enjoy! [white]Enjoy! [lime]Enjoy!"
		sleep 0.1
		api common.info.set "[lime]Enjoy! [yellow]Enjoy! [white]Enjoy!"
		sleep 0.1
	done
}

# Main
if [[ $# -eq 0 ]]; then
    run $@
else
    while [[ $# -gt 0 ]]; do
	case $1 in
	    -s|--setup)
        setup
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
