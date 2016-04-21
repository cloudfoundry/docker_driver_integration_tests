#!/bin/bash

cd `dirname $0`

pkill -f fakedriver

rm ../tmp_plugins_dir/fakedriver.*

mkdir -p ../mountdir

../exec/fakedriver -listenAddr="0.0.0.0:9876" -transport="tcp-json" -mountDir="../mountdir" -driversPath="../tmp_plugins_dir" &