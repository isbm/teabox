#!/usr/bin/bash

SOCK="/tmp/example-callback.sock"

#
# Set status of the logging output
#
function set_status() {
    m=$1
    $(echo "LOGGER-STATUS::$m" | nc -w0 -U $SOCK)
}

#
# Set title of the logging output widget
#
function set_title() {
    m=$1
    $(echo "LOGGER-TITLE::$m" | nc -w0 -U $SOCK)
}


function do_something() {
    location=$1
    echo "Putting something to $location..."
    sleep 1
    echo "Done"
}

set_status "Calling Something"
do_something $1

set_status "Calling package manager update"
sudo apt update

set_status "Finished!"
