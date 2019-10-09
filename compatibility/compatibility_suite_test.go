package compatibility_test

import (
	"code.cloudfoundry.org/docker_driver_integration_tests"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	session *gexec.Session
)

func TestCompatibility(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "Compatibility Suite")
}

var _ = BeforeSuite(func() {
	cmd := exec.Command(os.Getenv("DRIVER_CMD"), strings.Split(os.Getenv("DRIVER_OPTS"), ",")...)

	var err error
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session.Out).Should(gbytes.Say("driver-server.server.start"))
})


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
