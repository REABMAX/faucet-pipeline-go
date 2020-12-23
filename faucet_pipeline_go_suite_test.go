package faucet_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFaucetPipelineGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FaucetPipelineGo Suite")
}
