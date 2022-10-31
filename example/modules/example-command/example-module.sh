#!/usr/bin/bash

function do_something() {
    location=$1
    echo "Putting something to $location..."
    sleep 1
    echo "Done"
}

do_something $1
sudo apt update
