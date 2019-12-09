# Diego Volume Driver Integration Tests
These tests are used to test volume drivers against the Diego volume manager (aka *volman*).
# Installation

Prereqs:
- [go](https://golang.org/dl/)
- ginkgo and gomega; i.e.
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```
To install integration tests in your `GOPATH`:
```
go get -t code.cloudfoundry.org/docker_driver_integration_tests
```

# Configuration

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

Note: to run tests, you'll need to be in a containing project or `GOPATH` (eg. diego_release).

## Running example SMB compatibility
```
TEST_PACKAGE=docker_driver_integration_tests/compatibility \
fly -t persi execute \
-c /Users/pivotal/workspace/smb-volume-release/scripts/ci/run_docker_driver_integration_tests.build.yml \
-j persi/smbdriver-integration \
-i smb-volume-release-concourse-tasks=/Users/pivotal/workspace/smb-volume-release \
-i docker_driver_integration_tests=/Users/pivotal/go/src/code.cloudfoundry.org/docker_driver_integration_tests \
-i smb-volume-release=/Users/pivotal/workspace/smb-volume-release \
 --privileged
```

## Running example SMB
```
TEST_PACKAGE=docker_driver_integration_tests \
fly -t persi execute \
-c /Users/pivotal/workspace/smb-volume-release/scripts/ci/run_docker_driver_integration_tests.build.yml \
-j persi/smbdriver-integration \
-i smb-volume-release-concourse-tasks=/Users/pivotal/workspace/smb-volume-release \
-i docker_driver_integration_tests=/Users/pivotal/go/src/code.cloudfoundry.org/docker_driver_integration_tests \
-i smb-volume-release=/Users/pivotal/workspace/smb-volume-release \
 --privileged
```

## Running example SMB with lazy_unmount
```
TEST_PACKAGE=docker_driver_integration_tests/lazy_unmount \
fly -t persi execute \
-c /Users/pivotal/workspace/smb-volume-release/scripts/ci/run_docker_driver_integration_tests.build.yml \
-j persi/smbdriver-integration \
-i smb-volume-release-concourse-tasks=/Users/pivotal/workspace/smb-volume-release \
-i docker_driver_integration_tests=/Users/pivotal/go/src/code.cloudfoundry.org/docker_driver_integration_tests \
-i smb-volume-release=/Users/pivotal/workspace/smb-volume-release \
 --privileged
```


## Running example NFS
```
fly -t persi execute \
-c /Users/pivotal/workspace/nfs-volume-release/scripts/ci/run_docker_driver_integration_tests.build.yml \
-j persi/nfsdriver-integration \
-i nfs-volume-release-concourse-tasks=/Users/pivotal/workspace/nfs-volume-release \
-i docker_driver_integration_tests=/Users/pivotal/go/src/code.cloudfoundry.org/docker_driver_integration_tests \
-i nfs-volume-release=/Users/pivotal/workspace/nfs-volume-release \
-i mapfs-release=/Users/pivotal/workspace/mapfs-release \
 --privileged
```


## Running example NFS lazy unmount
```
fly -t persi execute \
-c $HOME/workspace/nfsv3driver/scripts/ci/run_docker_driver_integration_tests.build.yml \
-j nfs-driver/integration \
-i docker_driver_integration_tests=$HOME/workspace/docker_driver_integration_tests \
-i nfsv3driver=$HOME/workspace/nfsv3driver \
-i mapfs=$HOME/workspace/mapfs \
 --privileged
```