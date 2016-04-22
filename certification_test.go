package volume_driver_cert_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudfoundry-incubator/volman"
	"github.com/cloudfoundry-incubator/volman/volhttp"
	"github.com/cloudfoundry-incubator/volume_driver_cert"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Certify Volman with: ", func() {
	var (
		testLogger lager.Logger
		serverPort int = 8750

		volmanProcess        ifrit.Process
		volmanRunner         ifrit.Runner
		client               volman.Manager
		certificationFixture volume_driver_cert.CertificationFixture
		err                  error
	)

	BeforeEach(func() {
		fileName := os.Getenv("FIXTURE_FILENAME")
		Expect(fileName).NotTo(Equal(""))
		certificationFixture, err = volume_driver_cert.LoadCertificationFixture(fileName)
		Expect(err).NotTo(HaveOccurred())

		testLogger = lagertest.NewTestLogger("MainTest")

		volmanRunner = certificationFixture.CreateVolmanRunner(volmanPath)
		volmanProcess = ginkgomon.Invoke(volmanRunner)

		client = volhttp.NewRemoteClient(fmt.Sprintf("http://0.0.0.0:%d", serverPort))

	})

	AfterEach(func() {
		ginkgomon.Kill(volmanProcess)
	})

	Context("after starting", func() {
		It("should not exit", func() {
			Consistently(volmanRunner).ShouldNot(Exit())
		})
	})

	Context("after starting volman server", func() {
		var (
		mountPoint volman.MountResponse
		err error
		)
		It("should return list of drivers", func() {
			drivers, err := client.ListDrivers(testLogger)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(drivers.Drivers)).ToNot(Equal(0))
		})

		Context("when a valid volume is mounted",func() {
			BeforeEach(func() {
				mountPoint, err = client.Mount(testLogger, certificationFixture.DriverName, certificationFixture.CreateConfig.Name, certificationFixture.CreateConfig.Opts)
				Expect(err).NotTo(HaveOccurred())
				Expect(mountPoint.Path).NotTo(Equal(""))
			})
			AfterEach(func() {
				err = client.Unmount(testLogger, certificationFixture.DriverName, certificationFixture.CreateConfig.Name)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should mount a volume", func() {
				matches, err := filepath.Glob(mountPoint.Path)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(matches)).To(Equal(1))
			})

			It("should be possible to write to the mountPoint", func() {
				testFile := path.Join(mountPoint.Path, "test.txt")
				err = ioutil.WriteFile(testFile, []byte("hello persi"), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = os.Remove(testFile)
				Expect(err).NotTo(HaveOccurred())

				matches, err := filepath.Glob(mountPoint.Path + "/*")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(matches)).To(Equal(0))
			})
		})
		It("should unmount a volume given same volume ID", func() {
			mountPoint, err = client.Mount(testLogger, certificationFixture.DriverName, certificationFixture.CreateConfig.Name, certificationFixture.CreateConfig.Opts)
			Expect(err).NotTo(HaveOccurred())
			err = client.Unmount(testLogger, certificationFixture.DriverName, certificationFixture.CreateConfig.Name)
			Expect(err).NotTo(HaveOccurred())

			matches, err := filepath.Glob(mountPoint.Path)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(matches)).To(Equal(0))
		})

		It("should error, given an invalid driver name", func() {
			_, err := client.Mount(testLogger, "InvalidDriver", "vol", nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Driver 'InvalidDriver' not found in list of known drivers"))
		})

	})

})

func get(path string, volmanServerPort int) (body string, status string, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%d%s", volmanServerPort, path), nil)

	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", "", err
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	return string(bodyBytes[:]), response.Status, err
}
