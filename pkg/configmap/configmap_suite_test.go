package configmap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigmap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configmap Suite")
}
