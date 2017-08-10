package ccruncher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCcruncher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ccruncher Suite")
}
