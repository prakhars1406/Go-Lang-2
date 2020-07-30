package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAccountService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AccountService Suite")
}
