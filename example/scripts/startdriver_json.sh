#!/bin/bash

cd `dirname $0`

mkdir -p ../mountdir

../exec/fakedriver -listenAddr="0.0.0.0:9876" -transport="tcp-json" -mountDir="../mountdir" -driversPath="../tmp_plugins_dir" &