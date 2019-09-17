package lazy_unmount_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLazyUnmount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LazyUnmount Suite")
}
