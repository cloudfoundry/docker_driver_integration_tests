package volume_driver_cert_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/ginkgo/config"
)

var(
	volmanPath string
  err error
)

func TestCertification(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certification Suite")
}

var _ = SynchronizedBeforeSuite(func()[]byte{
	Expect(config.GinkgoConfig.ParallelTotal).To(Equal(1),"DRIVER CERTIFICATION TESTS DO NOT RUN IN PARALLEL!!!")
	// TODO--surely this awful path can't be the only way to get to packages in "vendor?"
	buildpath, err := gexec.Build("../volume_driver_cert/vendor/github.com/cloudfoundry-incubator/volman/cmd/volman", "-race")
	Expect(err).NotTo(HaveOccurred())

	return []byte(buildpath)
},func(buildpath []byte){
	volmanPath = string(buildpath)
})

var _ = SynchronizedAfterSuite(func(){},func(){
	gexec.CleanupBuildArtifacts()
})