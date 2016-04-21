#!/bin/bash

cd `dirname $0`

pkill -f fakedriver

rm ../tmp_plugins_dir/fakedriver.*

mkdir -p ../mountdir

../exec/fakedriver -listenAddr="0.0.0.0:9776" -transport="tcp" -mountDir="../mountdir" -driversPath="../tmp_plugins_dir" &
