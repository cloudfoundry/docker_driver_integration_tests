#!/bin/bash

cd `dirname $0`

pkill -f fakedriver

mkdir -p ~/voldriver_plugins
rm ~/voldriver_plugins/fakedriver.*

mkdir -p ../mountdir

driversPath=$HOME/voldriver_plugins
../exec/fakedriver -listenAddr="127.0.0.1:9876" -transport="tcp-json" -mountDir="../mountdir" -driversPath="${driversPath}" -requireSSL=true -caFile="../certs/fakedriver_ca.crt" -certFile="../certs/fakedriver_server.crt" -keyFile="../certs/fakedriver_server.key" -clientCertFile="../certs/fakedriver_client.crt" -clientKeyFile="../certs/fakedriver_client.key" -insecureSkipVerify=true &
