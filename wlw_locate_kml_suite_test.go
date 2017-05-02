package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWlwLocateKml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WlwLocateKml Suite")
}
