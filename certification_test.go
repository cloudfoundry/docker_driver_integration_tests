package volume_driver_cert_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"hash/crc32"


	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	"code.cloudfoundry.org/volume_driver_cert"

	"context"

	"time"

	"code.cloudfoundry.org/voldriver"
	"code.cloudfoundry.org/voldriver/driverhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("Certify with: ", func() {
	var (
		err error

		testLogger           lager.Logger
		testContext          context.Context
		testEnv              voldriver.Env
		certificationFixture volume_driver_cert.CertificationFixture
		driverClient         voldriver.Driver
		errResponse          voldriver.ErrorResponse

		mountResponse voldriver.MountResponse
	)

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("MainTest")
		testContext = context.TODO()
		testEnv = driverhttp.NewHttpDriverEnv(testLogger, testContext)

		fileName := os.Getenv("FIXTURE_FILENAME")
		Expect(fileName).NotTo(Equal(""))

		certificationFixture, err = volume_driver_cert.LoadCertificationFixture(fileName)
		Expect(err).NotTo(HaveOccurred())
		testLogger.Info("fixture", lager.Data{"filename": fileName, "context": certificationFixture})

		driverClient, err = driverhttp.NewRemoteClient(certificationFixture.DriverAddress, certificationFixture.TLSConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("given a driver", func() {
		It("should respond with Capabilities", func() {
			resp := driverClient.Capabilities(testEnv)
			Expect(resp.Capabilities).NotTo(BeNil())
			Expect(resp.Capabilities.Scope).To(Or(Equal("local"), Equal("global")))
		})
	})

	Context("given a created volume", func() {
		BeforeEach(func() {
			errResponse = driverClient.Create(testEnv, certificationFixture.CreateConfig)
			Expect(errResponse.Err).To(Equal(""))
		})

		AfterEach(func() {
			errResponse = driverClient.Remove(testEnv, voldriver.RemoveRequest{
				Name: certificationFixture.CreateConfig.Name,
			})
			Expect(errResponse.Err).To(Equal(""))
		})

		Context("given a mounted volume", func() {
			BeforeEach(func() {
				mountResponse = driverClient.Mount(testEnv, voldriver.MountRequest{
					Name: certificationFixture.CreateConfig.Name,
				})
				Expect(mountResponse.Err).To(Equal(""))
				Expect(mountResponse.Mountpoint).NotTo(Equal(""))
			})

			AfterEach(func() {
				errResponse = driverClient.Unmount(testEnv, voldriver.UnmountRequest{
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

			It("should be be able to write and read many files", func() {
				var wg sync.WaitGroup

				numFiles := 100
				mu := sync.Mutex{}
				wg.Add(numFiles)
				startTime := time.Now()
				file_names := make(map[string]uint32)
				for i := 0; i < numFiles; i++ {
					go func() {
						defer wg.Done()
						fileName := "certtest-" + randomString(10)
						contents := "contents-" + randomString(1000000)
						mu.Lock()
						check := crc32.Checksum([]byte(contents), crc32.IEEETable)
						file_names[fileName] = check
						mu.Unlock()
						testLogger.Debug("writing-test-file", lager.Data{"mountpoint": mountResponse.Mountpoint})
						testFile := path.Join(mountResponse.Mountpoint, fileName)
						testLogger.Debug("writing-test-file", lager.Data{"filepath": testFile})
						err := ioutil.WriteFile(testFile, []byte(contents), 0644)
						Expect(err).NotTo(HaveOccurred())

						matches, err := filepath.Glob(mountResponse.Mountpoint + "/" + fileName)
						Expect(err).NotTo(HaveOccurred())
						Expect(len(matches)).To(Equal(1))
					}()
				}
				wg.Wait()

				elapsed := time.Since(startTime)
				testLogger.Info("file-write-duration", lager.Data{"duration-in-seconds":elapsed.Seconds()})
				fmt.Printf("File Write Duration: %f\n", elapsed.Seconds())

				// READ ============================
				wg.Add(len(file_names))
				startTime_read := time.Now()
				for key, val := range file_names {
					go func(key string, val uint32) {
						defer wg.Done()

						testLogger.Debug("reading-test-file", lager.Data{"mountpoint": mountResponse.Mountpoint})
						testFile := path.Join(mountResponse.Mountpoint, key)
						testLogger.Debug("reading-test-file", lager.Data{"filepath": testFile, "checksum":val})
						contents, err := ioutil.ReadFile(testFile)
						Expect(err).NotTo(HaveOccurred())
						check := crc32.Checksum(contents, crc32.IEEETable)
						Expect(check).To(Equal(val))
					}(key, val)
				}

				wg.Wait()
				elapsed_read := time.Since(startTime_read)
				testLogger.Info("file-read-duration", lager.Data{"duration-in-seconds":elapsed_read.Seconds()})
				fmt.Printf("File Read Duration: %f\n", elapsed_read.Seconds())

				// DELETE =====================
				wg.Add(len(file_names))
				startTime_delete := time.Now()
				for key := range file_names {
					go func(key string) {
						defer wg.Done()

						testLogger.Debug("deleting-test-file", lager.Data{"mountpoint": mountResponse.Mountpoint})
						testFile := path.Join(mountResponse.Mountpoint, key)
						testLogger.Debug("deleting-test-file", lager.Data{"filepath": testFile})
						err := os.Remove(testFile)
						if err != nil {testLogger.Error("error-deleting-file", err)}
						Expect(err).NotTo(HaveOccurred())
						matches, _ := filepath.Glob(mountResponse.Mountpoint + "/" + key)
						Expect(len(matches)).To(Equal(0))
					}(key)
				}

				wg.Wait()
				elapsed_del := time.Since(startTime_delete)
				testLogger.Info("delete-file-duration", lager.Data{"duration-in-seconds":elapsed_del.Seconds()})
				fmt.Printf("File Delete Duration: %f\n", elapsed_del.Seconds())

			})

			Context("when that volume is mounted again (for another container) and then unmounted", func() {
				BeforeEach(func() {
					secondMountResponse := driverClient.Mount(testEnv, voldriver.MountRequest{
						Name: certificationFixture.CreateConfig.Name,
					})
					Expect(secondMountResponse.Err).To(Equal(""))
					Expect(secondMountResponse.Mountpoint).NotTo(Equal(""))

					errResponse = driverClient.Unmount(testEnv, voldriver.UnmountRequest{
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
		errResponse = driverClient.Create(testEnv, certificationFixture.CreateConfig)
		Expect(errResponse.Err).To(Equal(""))

		mountResponse := driverClient.Mount(testEnv, voldriver.MountRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(mountResponse.Err).To(Equal(""))

		errResponse = driverClient.Unmount(testEnv, voldriver.UnmountRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(errResponse.Err).To(Equal(""))
		Expect(cellClean(mountResponse.Mountpoint)).To(Equal(true))

		errResponse = driverClient.Remove(testEnv, voldriver.RemoveRequest{
			Name: certificationFixture.CreateConfig.Name,
		})
		Expect(errResponse.Err).To(Equal(""))

	})
})

// given a mounted mountpoint, tests creation of a file on that mount point
func testFileWrite(logger lager.Logger, mountResponse voldriver.MountResponse) {
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

	err = os.Remove(testFile)
	Expect(err).NotTo(HaveOccurred())

	matches, err = filepath.Glob(path.Join(mountResponse.Mountpoint, fileName))
	Expect(err).NotTo(HaveOccurred())
	Expect(len(matches)).To(Equal(0))
}

func cellClean(mountpoint string) bool {
	matches, err := filepath.Glob(mountpoint)
	Expect(err).NotTo(HaveOccurred())
	return len(matches) == 0
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
