package compatibility_test

import (
	"code.cloudfoundry.org/docker_driver_integration_tests"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

type VolumeServiceBrokerBinding struct {
	VolumeMounts []struct {
		Device       struct {
			VolumeID    string `json:"volume_id"`
			MountConfig map[string]interface{} `json:"mount_config"`
		} `json:"device"`
	} `json:"volume_mounts"`
}

var (
	integrationFixtureTemplate = docker_driver_integration_tests.LoadFixtureTemplate()
	bindingsFixture            = LoadVolumeServiceBrokerBindingsFixture()
)

func TestCompatibility(t *testing.T) {

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
