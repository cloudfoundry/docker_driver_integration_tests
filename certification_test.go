package volume_driver_cert_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	"github.com/cloudfoundry-incubator/volume_driver_cert"

	"github.com/cloudfoundry-incubator/voldriver"
	"github.com/cloudfoundry-incubator/voldriver/driverhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Certify with: ", func() {
	var (
		err error

		testLogger           lager.Logger
		certificationFixture volume_driver_cert.CertificationFixture
		driverClient         voldriver.Driver
		errResponse          voldriver.ErrorResponse

		mountResponse voldriver.MountResponse
	)

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("MainTest")

		fileName := os.Getenv("FIXTURE_FILENAME")
		Expect(fileName).NotTo(Equal(""))

		certificationFixture, err = volume_driver_cert.LoadCertificationFixture(fileName)
		Expect(err).NotTo(HaveOccurred())
		testLogger.Info("fixture", lager.Data{"filename": fileName, "context": certificationFixture})

		driverClient, err = driverhttp.NewRemoteClient(certificationFixture.DriverAddress, certificationFixture.TLSConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("given a created volume", func() {
		BeforeEach(func() {
			errResponse = driverClient.Create(testLogger, certificationFixture.CreateConfig)
			Expect(errResponse.Err).To(Equal(""))
		})

		AfterEach(func() {
			errResponse = driverClient.Remove(testLogger, voldriver.RemoveRequest{
				Name: certificationFixture.CreateConfig.Name,
			})
			Expect(errResponse.Err).To(Equal(""))
		})

		Context("given a mounted volume", func() {
			BeforeEach(func() {
				mountResponse = driverClient.Mount(testLogger, voldriver.MountRequest{
					Name: certificationFixture.CreateConfig.Name,
				})
				Expect(mountResponse.Err).To(Equal(""))
				Expect(mountResponse.Mountpoint).NotTo(Equal(""))
			})

			AfterEach(func() {
				errResponse = driverClient.Unmount(testLogger, voldriver.UnmountRequest{
					Name: certificationFixture.CreateConfig.Name,
				})
				Expect(errResponse.Err).To(Equal(""))
			})

			It("should mount that volume", func() {
				matches, err := filepath.Glob(mountResponse.Mountpoint)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(matches)).To(Equal(1))
			})

			It("should write to that volume", func() {
				testFileWrite(testLogger, mountResponse)
			})

			Context("when that volume is mounted again (for another container) and then unmounted", func() {
				BeforeEach(func() {
					secondMountResponse := driverClient.Mount(testLogger, voldriver.MountRequest{
						Name: certificationFixture.CreateConfig.Name,
					})
					Expect(secondMountResponse.Err).To(Equal(""))
					Expect(secondMountResponse.Mountpoint).NotTo(Equal(""))

					errResponse = driverClient.Unmount(testLogger, voldriver.UnmountRequest{
						Name: certificationFixture.CreateConfig.Name,
					})
					Expect(errResponse.Err).To(Equal(""))
				})

				It("should still write to that volume", func() {
					testFileWrite(testLogger, mountResponse)
				})
			})
		})

		Context("given an unmounted volume", func() {
			// the It should unmount a volume given same volume ID test should be here!
		})
	})

	It("should unmount a volume given same volume ID", func() {
		errResponse = driverClient.Create(testLogger, certificationFixture.CreateConfig)
		Expect(errResponse.Err).To(Equal(""))

		mountResponse := driverClient.Mount(testLogger, voldriver.MountRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(mountResponse.Err).To(Equal(""))

		errResponse = driverClient.Unmount(testLogger, voldriver.UnmountRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(errResponse.Err).To(Equal(""))

		errResponse = driverClient.Remove(testLogger, voldriver.RemoveRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(errResponse.Err).To(Equal(""))

		Expect(cellClean(mountResponse.Mountpoint)).To(Equal(true))
	})
})

// given a mounted mountpoint, tests creation of a file on that mount point
func testFileWrite(logger lager.Logger, mountResponse voldriver.MountResponse) {
	logger = logger.Session("test-file-write")
	logger.Info("start")
	defer logger.Info("end")

	logger.Info("writing-test-file", lager.Data{"mountpoint": mountResponse.Mountpoint})
	testFile := path.Join(mountResponse.Mountpoint, "test.txt")
	logger.Info("writing-test-file", lager.Data{"filepath": testFile})
	err := ioutil.WriteFile(testFile, []byte("hello persi"), 0644)
	Expect(err).NotTo(HaveOccurred())

	matches, err := filepath.Glob(mountResponse.Mountpoint + "/test.txt")
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(1))

	err = os.Remove(testFile)
	Expect(err).NotTo(HaveOccurred())

	matches, err = filepath.Glob(mountResponse.Mountpoint + "/test.txt")
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(0))
}

func cellClean(mountpoint string) bool {
	matches, err := filepath.Glob(mountpoint)
	Expect(err).NotTo(HaveOccurred())
	return len(matches) == 0
}