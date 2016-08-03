# Diego Volume Driver Certification Tests
These tests are used to certify volume drivers against the Diego volume manager (aka *volman*).
# Installation

Prereqs:
- [go](https://golang.org/dl/)
- ginkgo and gomega; i.e.
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```
To install cert tests in your `GOPATH`:
```
go get -t code.cloudfoundry.org/volume_driver_cert
```

# Certification

- Make sure that your driver is running (you can see the start/stop scripts in [example](example/).
- Create a fixture file that contains connection information for your driver

```
{
  "volman_driver_path": "~/voldriver_plugins",
  "driver_address": "http://127.0.0.1:9786",
  "driver_name": "fakedriver",
  "create_config": {
    "Name": "fake-volume",
    "Opts": {"volume_id":"fake-volume"}
  }
}
```
NB: Optionally, you can supply a TLS Config as follows:-
```
{
  "volman_driver_path": "~/voldriver_plugins",
  ...
  "tls_config": {
    "InsecureSkipVerify": true,
    "CAFile": "localdriver_ca.crt",
    "CertFile":"localdriver_client.crt",
    "KeyFile":"localdriver_client.key"
  }
}
```
- Run ginkgo not in parallel mode.  (If you use -p, the tests will fail.)

```
ginkgo
```

Note: to run tests, you'll need to be in a containing project or `GOPATH` (eg. diego_release).

# For example

Our support example driver is a [local volume driver](https://github.com/cloudfoundry-incubator/local-volume-release/). We run certifications in our CI to certify it against both Volman and Docker.

The definitions of those tasks are in local-volume-release/scripts. They can be used to create a Concourse pipeline or run the certifications locally. We'll focus on the Volman certifications here.

## Running example
We used a start/stop script to manage our driver (local-volume-release/scripts/startdriver* and local-volume-release/scripts/stopdriver.sh).

We created stock fixture files (local-volume-release/scripts/fixtures/*) and certs for the encrypted tests (local-volume-release/scripts/certs).

Finally, we encapsulated the running of the various types of the driver (json plain/tls, unix sockets, tcp) in a script (local-volume-release/scripts/run-certification-tests).

To run (with all the prereqs met):
```
  local-volume-release/scripts/run-certification-tests
```


