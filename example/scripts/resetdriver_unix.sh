#!/bin/bash

cd `dirname $0`

rm ../tmp_plugins_dir/fakedriver.*

mkdir -p ../mountdir

../exec/fakedriver -listenAddr="../tmp_plugins_dir/fakedriver.sock" -transport="unix" -mountDir="../mountdir" -driversPath="../tmp_plugins_dir" &
