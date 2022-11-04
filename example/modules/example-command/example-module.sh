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
