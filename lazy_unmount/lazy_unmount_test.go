package lazy_unmount_test

import (
	"code.cloudfoundry.org/docker_driver_integration_tests"
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

var _ = Describe("LazyUnmount", func() {
	var (
		err error

		testLogger           lager.Logger
		testContext          context.Context
		testEnv              dockerdriver.Env
		certificationFixture docker_driver_integration_tests.CertificationFixture
		driverClient         dockerdriver.Driver
		errResponse          dockerdriver.ErrorResponse

		mountResponse dockerdriver.MountResponse
	)

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("LazyUnmountTest")
		testContext = context.TODO()
		testEnv = driverhttp.NewHttpDriverEnv(testLogger, testContext)

		fileName := os.Getenv("FIXTURE_FILENAME")
		Expect(fileName).NotTo(Equal(""))

		certificationFixture, err = docker_driver_integration_tests.LoadCertificationFixture(fileName)
		Expect(err).NotTo(HaveOccurred())
		testLogger.Info("fixture", lager.Data{"filename": fileName, "context": certificationFixture})

		driverClient, err = driverhttp.NewRemoteClient(certificationFixture.DriverAddress, certificationFixture.TLSConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("given a created volume", func() {
		BeforeEach(func() {
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

				cmd := exec.Command("bash","-c", "cat /proc/mounts | grep -E '"+mountResponse.Mountpoint+"'")
				Expect(cmdRunner(cmd)).To(Equal(0))
			})

			Context("when the nfs server is slow", func() {
				BeforeEach(func(){
					addNetworkDelay()
				})

				AfterEach(func() {
					removeNetworkDelay()
				})

				It("should unmount lazily", func() {
					block := make(chan bool)
					go func() {
						defer GinkgoRecover()
						testFileWrite(testLogger, mountResponse)
						block <- true
					}()
					Consistently(block, 2).ShouldNot(Receive())

					go func() {
						defer GinkgoRecover()
						errResponse = driverClient.Unmount(testEnv, dockerdriver.UnmountRequest{
							Name: certificationFixture.CreateConfig.Name,
						})
						Expect(errResponse.Err).To(Equal(""))
					}()

					Eventually(func() int {
						cmd := exec.Command("bash","-c", "cat /proc/mounts | grep -E '"+mountResponse.Mountpoint+"'")
						return cmdRunner(cmd)

					}, 5, 500 * time.Millisecond).Should(Equal(1))
				})
			})
		})
	})
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

func addNetworkDelay() {
	cmd := exec.Command("tc", "qdisc","add","dev","lo","root","netem","delay","2000ms")
	Expect(cmdRunner(cmd)).To(Equal(0))
}

func removeNetworkDelay() {
	cmd := exec.Command("tc","qdisc","del","dev","lo","root","netem")
	Expect(cmdRunner(cmd)).To(Equal(0))
}

func cmdRunner(cmd *exec.Cmd) int {
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 10).Should(Exit())
	return session.ExitCode()
}