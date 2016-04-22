#!/bin/bash

set -x

cd `dirname $0`

pkill -f fakedriver

mkdir -p ~/voldriver_plugins
rm ~/voldriver_plugins/fakedriver.*

mkdir -p ../mountdir

# temporarily create a sock file in order to find an absolute path for it
touch ~/voldriver_plugins/fakedriver.sock
listenAddr=$(realpath ~/voldriver_plugins/fakedriver.sock)
rm ~/voldriver_plugins/fakedriver.sock

../exec/fakedriver -listenAddr="${listenAddr}" -transport="unix" -mountDir="../mountdir" &
