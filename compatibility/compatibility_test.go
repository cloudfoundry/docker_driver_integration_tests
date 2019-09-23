package compatibility_test

import (
	"code.cloudfoundry.org/dockerdriver"
	"code.cloudfoundry.org/dockerdriver/driverhttp"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var _ = Describe("Compatibility", func() {
	var (
		err error

		testLogger           lager.Logger
		testContext          context.Context
		testEnv              dockerdriver.Env
		driverClient         dockerdriver.Driver
		errResponse          dockerdriver.ErrorResponse

		mountResponse dockerdriver.MountResponse
	)

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("CompatibilityTest")
		testContext = context.TODO()
		testEnv = driverhttp.NewHttpDriverEnv(testLogger, testContext)

		driverClient, err = driverhttp.NewRemoteClient(certificationFixtureTemplate.DriverAddress, certificationFixtureTemplate.TLSConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	for _, binding := range bindingsFixture {

		cf := certificationFixtureTemplate
		cf.CreateConfig = certificationFixtureTemplate.CreateConfig
		cf.CreateConfig.Name = binding.VolumeMounts[0].Device.VolumeID
		for k, v := range binding.VolumeMounts[0].Device.MountConfig {
			if k == "source" || k == "username" || k == "password" {
				continue
			}
			cf.CreateConfig.Opts[k] = v
		}

		Context("given a created volume", func() {

			var certificationFixture = cf

			BeforeEach(func() {
				testLogger.Info("using fixture", lager.Data{"fixture": certificationFixture})
				errResponse = driverClient.Create(testEnv, certificationFixture.CreateConfig)
				Expect(errResponse.Err).To(Equal(""))
			})

			AfterEach(func() {
				errResponse = driverClient.Remove(testEnv, dockerdriver.RemoveRequest{
					Name: certificationFixture.CreateConfig.Name,
				})
				Expect(errResponse.Err).To(Equal(""))
			})

			Context("given a mounted volume", func() {
				BeforeEach(func() {
					mountResponse = driverClient.Mount(testEnv, dockerdriver.MountRequest{
						Name: certificationFixture.CreateConfig.Name,
					})
					Expect(mountResponse.Err).To(Equal(""))
					Expect(mountResponse.Mountpoint).NotTo(Equal(""))

					cmd := exec.Command("bash", "-c", "cat /proc/mounts | grep -E '"+mountResponse.Mountpoint+"'")
					Expect(cmdRunner(cmd)).To(Equal(0))
				})

				AfterEach(func() {
					errResponse = driverClient.Unmount(testEnv, dockerdriver.UnmountRequest{
						Name: certificationFixture.CreateConfig.Name,
					})
					Expect(errResponse.Err).To(Equal(""))
				})

				It("should be able to write a file", func() {
					testFileWrite(testLogger, mountResponse)
				})
			})
		})
	}
})

func testFileWrite(logger lager.Logger, mountResponse dockerdriver.MountResponse) {
	logger = logger.Session("test-file-write")
	logger.Info("start")
	defer logger.Info("end")

	fileName := "certtest-" + randomString(10)

	logger.Info("writing-test-file", lager.Data{"mountpoint": mountResponse.Mountpoint})
	testFile := path.Join(mountResponse.Mountpoint, fileName)
	logger.Info("writing-test-file", lager.Data{"filepath": testFile})
	err := ioutil.WriteFile(testFile, []byte("hello persi"), 0644)
	Expect(err).NotTo(HaveOccurred())

	matches, err := filepath.Glob(mountResponse.Mountpoint + "/" + fileName)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(1))

	bytes, err := ioutil.ReadFile(mountResponse.Mountpoint + "/" + fileName)
	Expect(err).NotTo(HaveOccurred())
	Expect(bytes).To(Equal([]byte("hello persi")))

	err = os.Remove(testFile)
	Expect(err).NotTo(HaveOccurred())

	matches, err = filepath.Glob(path.Join(mountResponse.Mountpoint, fileName))
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(0))
}

var isSeeded = false

func randomString(n int) string {
	if !isSeeded {
		rand.Seed(time.Now().UnixNano())
		isSeeded = true
	}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func cmdRunner(cmd *exec.Cmd) int {
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 10).Should(Exit())
	return session.ExitCode()
}