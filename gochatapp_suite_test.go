package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGochatapp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gochatapp Suite")
}
