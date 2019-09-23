package compatibility_test

import (
	"code.cloudfoundry.org/docker_driver_integration_tests"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"os"

	"testing"
)

type VolumeServiceBrokerBinding struct {
	Credentials struct {
	} `json:"credentials"`
	VolumeMounts []struct {
		Driver       string `json:"driver"`
		ContainerDir string `json:"container_dir"`
		Mode         string `json:"mode"`
		DeviceType   string `json:"device_type"`
		Device       struct {
			VolumeID    string `json:"volume_id"`
			MountConfig map[string]interface{} `json:"mount_config"`
		} `json:"device"`
	} `json:"volume_mounts"`
}

var (
	certificationFixtureTemplate = docker_driver_integration_tests.LoadCertificationFixtureTemplate()
	bindingsFixture              = LoadVolumeServiceBrokerBindingsFixture()
)

func TestCompatibility(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Compatibility Suite")
}

func LoadVolumeServiceBrokerBindingsFixture() []VolumeServiceBrokerBinding {
	bytes, err := ioutil.ReadFile("bindings.json")
	if err != nil {
		panic(err.Error())
	}

	bindings := []VolumeServiceBrokerBinding{}
	err = json.Unmarshal(bytes, &bindings)
	if err != nil {
		panic(err.Error())
	}

	return bindings
}
