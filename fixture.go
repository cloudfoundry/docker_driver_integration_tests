package volume_driver_cert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/cloudfoundry-incubator/voldriver"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

type CertificationFixture struct {
	VolmanDriverPath string                  `json:"volman_driver_path"`
	DriverAddress    string                  `json:"driver_address"`
	DriverName       string                  `json:"driver_name"`
	CreateConfig     voldriver.CreateRequest `json:"create_config"`
	TLSConfig        *voldriver.TLSConfig     `json:"tls_config,omitempty"`
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

	// make sure that the paths are absolute
	usr, err := user.Current()
	if certificationFixture.VolmanDriverPath[:2] == "~/" {
		if err != nil {
			return CertificationFixture{}, err
		}
		certificationFixture.VolmanDriverPath = filepath.Join(usr.HomeDir, certificationFixture.VolmanDriverPath[2:])
	}
	if certificationFixture.DriverAddress[:2] == "~/" {
		certificationFixture.DriverAddress = filepath.Join(usr.HomeDir, certificationFixture.DriverAddress[2:])
	}
	if certificationFixture.TLSConfig != nil {
		if certificationFixture.TLSConfig.CAFile[:2] == "~/" {
			certificationFixture.TLSConfig.CAFile = filepath.Join(usr.HomeDir, certificationFixture.TLSConfig.CAFile[2:])
		}
		if certificationFixture.TLSConfig.CertFile[:2] == "~/" {
			certificationFixture.TLSConfig.CertFile = filepath.Join(usr.HomeDir, certificationFixture.TLSConfig.CertFile[2:])
		}
		if certificationFixture.TLSConfig.KeyFile[:2] == "~/" {
			certificationFixture.TLSConfig.KeyFile = filepath.Join(usr.HomeDir, certificationFixture.TLSConfig.KeyFile[2:])
		}
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
			"-volmanDriverPaths", cf.VolmanDriverPath,
		),
		StartCheck: "volman.started",
	})

}
