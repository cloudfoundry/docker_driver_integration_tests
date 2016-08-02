package volume_driver_cert_test

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/volume_driver_cert"

	"code.cloudfoundry.org/voldriver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("certification/fixture.go", func() {
	var (
		err                  error
		tmpDir, tmpFileName  string
		certificationFixture volume_driver_cert.CertificationFixture
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "certification")
		Expect(err).NotTo(HaveOccurred())

		tmpFile, err := ioutil.TempFile(tmpDir, "certification-fixture.json")
		Expect(err).NotTo(HaveOccurred())

		tmpFileName = tmpFile.Name()

		certificationFixture = volume_driver_cert.CertificationFixture{}
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
			certificationFixture, err = volume_driver_cert.LoadCertificationFixture(tmpFileName)
			Expect(err).NotTo(HaveOccurred())

			Expect(certificationFixture.VolmanDriverPath).To(ContainSubstring("fake-path-to-driver"))
			Expect(certificationFixture.CreateConfig.Name).To(Equal("fake-request"))
		})
	})

	Context("#SaveCertificationFixture", func() {
		BeforeEach(func() {
			certificationFixture = volume_driver_cert.CertificationFixture{
				VolmanDriverPath: "fake-path-to-driver",
				DriverName:       "fakedriver",
				CreateConfig: voldriver.CreateRequest{
					Name: "fake-request",
					Opts: map[string]interface{}{"key": "value"},
				},
			}
		})

		It("saves the fake certification fixture", func() {
			err = volume_driver_cert.SaveCertificationFixture(certificationFixture, tmpFileName)
			Expect(err).NotTo(HaveOccurred())

			bytes, err := ioutil.ReadFile(tmpFileName)
			Expect(err).ToNot(HaveOccurred())

			readFixture := volume_driver_cert.CertificationFixture{}
			json.Unmarshal(bytes, &readFixture)

			Expect(readFixture.VolmanDriverPath).To(Equal("fake-path-to-driver"))
			Expect(readFixture.CreateConfig.Name).To(Equal("fake-request"))
		})
	})

})
