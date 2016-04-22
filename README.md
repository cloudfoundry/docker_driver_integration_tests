# Diego Volume Driver Certification Tests
These tests are used to certify volume drivers against the Diego volume manager (aka *volman*).
# Installation

Prereqs:
- [go](https://storage.googleapis.com/golang/go1.4.3.darwin-amd64.pkg)

```
git clone git@github.com:cloudfoundry-incubator/volume_driver_cert.git
```
or
```
go get git@github.com:cloudfoundry-incubator/volume_driver_cert.git
```

# Certification

- Make sure that your driver is running (you can see the start/stop scripts in [example](example/).
- Create a fixture file that contains connection information for your driver

```
{
  "volman_driver_path": "~/voldriver_plugins",
  "driver_name": "fakedriver",
  "create_config": {
    "Name": "fake-volume",
    "Opts": {"volume_id":"fake-volume"}
  }
}
```
- Run ginkgo not in parallel mode.  (If you use -p, the tests will fail.)

```
ginkgo -r
```

Note: to run tests, you'll need to be in a containing project or gopath (eg. diego_release).

## Examples
There are sample scripts and fixture files in the [example](example/) folder that run certification tests against the volman test driver.
These scripts start and stop the driver and test 3 different communication protocols, but your certification test need not be as complicated.
