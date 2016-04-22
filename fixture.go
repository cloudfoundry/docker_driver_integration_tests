package volume_driver_cert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

type CertificationFixture struct {
	VolmanDriverPath  string                  `json:"volman_driver_path"`
	DriverName        string                  `json:"driver_name"`
	CreateConfig      voldriver.CreateRequest `json:"create_config"`
}

func NewCertificationFixture(volmanDriverPath string, driverName string, createConfig voldriver.CreateRequest) *CertificationFixture {
	return &CertificationFixture{
		VolmanDriverPath:  volmanDriverPath,
		DriverName:        driverName,
		CreateConfig:      createConfig,
	}
}
func LoadCertificationFixture(fileName string) (CertificationFixture, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return CertificationFixture{}, err
	}

	certificationFixture := CertificationFixture{}
	err = json.Unmarshal(bytes, &certificationFixture)
	if err != nil {
		return CertificationFixture{}, err
	}

	// make sure that the plugins path is absolute
	if certificationFixture.VolmanDriverPath[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return CertificationFixture{}, err
		}
		certificationFixture.VolmanDriverPath = filepath.Join(usr.HomeDir, certificationFixture.VolmanDriverPath[2:])
	}
	certificationFixture.VolmanDriverPath, err = filepath.Abs(certificationFixture.VolmanDriverPath)
	if err != nil {
		return CertificationFixture{}, err
	}

	return certificationFixture, nil
}

func SaveCertificationFixture(fixture CertificationFixture, fileName string) error {
	bytes, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, bytes, 0666)
}

func (cf *CertificationFixture) CreateVolmanRunner(volmanPath string) ifrit.Runner {
	return ginkgomon.New(ginkgomon.Config{
		Name: "volman",
		Command: exec.Command(
				volmanPath,
			"-listenAddr", fmt.Sprintf("0.0.0.0:%d", 8750),
			"-driversPath", cf.VolmanDriverPath,
		),
		StartCheck: "volman.started",
	})

}
