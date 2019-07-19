package docker_driver_integration_tests_test

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/dockerdriver"
	"code.cloudfoundry.org/docker_driver_integration_tests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("certification/fixture.go", func() {
	var (
		err                  error
		tmpDir, tmpFileName  string
		certificationFixture docker_driver_integration_tests.CertificationFixture
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "certification")
		Expect(err).NotTo(HaveOccurred())

		tmpFile, err := ioutil.TempFile(tmpDir, "certification-fixture.json")
		Expect(err).NotTo(HaveOccurred())

		tmpFileName = tmpFile.Name()
		tmpFile.Close()

		certificationFixture = docker_driver_integration_tests.CertificationFixture{}
	})

	AfterEach(func() {
		err = os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("#LoadCertificationFixture", func() {
		BeforeEach(func() {
			certificationFixtureContent := `{
 						"volman_driver_path": "fake-path-to-driver",
  					"driver_address": "http://fakedriver_address",
  					"driver_name": "fakedriver",
						"create_config": {
						    "Name": "fake-request",
						    "Opts": {"key":"value"}
 						},
						"tls_config": {
								"InsecureSkipVerify": true,
								"CAFile": "fakedriver_ca.crt",
								"CertFile":"fakedriver_client.crt",
								"KeyFile":"fakedriver_client.key"
							}
						}`

			err = ioutil.WriteFile(tmpFileName, []byte(certificationFixtureContent), 0666)
			Expect(err).NotTo(HaveOccurred())
		})

		It("loads the fake certification fixture", func() {
			certificationFixture, err = docker_driver_integration_tests.LoadCertificationFixture(tmpFileName)
			Expect(err).NotTo(HaveOccurred())

			Expect(certificationFixture.VolmanDriverPath).To(ContainSubstring("fake-path-to-driver"))
			Expect(certificationFixture.CreateConfig.Name).To(Equal("fake-request"))
		})
	})

	Context("#SaveCertificationFixture", func() {
		BeforeEach(func() {
			certificationFixture = docker_driver_integration_tests.CertificationFixture{
				VolmanDriverPath: "fake-path-to-driver",
				DriverName:       "fakedriver",
				CreateConfig: dockerdriver.CreateRequest{
					Name: "fake-request",
					Opts: map[string]interface{}{"key": "value"},
				},
			}
		})

		It("saves the fake certification fixture", func() {
			err = docker_driver_integration_tests.SaveCertificationFixture(certificationFixture, tmpFileName)
			Expect(err).NotTo(HaveOccurred())

			bytes, err := ioutil.ReadFile(tmpFileName)
			Expect(err).ToNot(HaveOccurred())

			readFixture := docker_driver_integration_tests.CertificationFixture{}
			json.Unmarshal(bytes, &readFixture)

			Expect(readFixture.VolmanDriverPath).To(Equal("fake-path-to-driver"))
			Expect(readFixture.CreateConfig.Name).To(Equal("fake-request"))
		})
	})

})
