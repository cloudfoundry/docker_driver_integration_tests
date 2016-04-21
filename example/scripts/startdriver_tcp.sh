#!/bin/bash

cd `dirname $0`

mkdir -p ../mountdir

../exec/fakedriver -listenAddr="0.0.0.0:9776" -transport="tcp" -mountDir="../mountdir" -driversPath="../tmp_plugins_dir" &