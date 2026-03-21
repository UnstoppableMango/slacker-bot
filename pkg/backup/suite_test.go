package backup

import (
	"io"
	"testing"

	"charm.land/log/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBackup(t *testing.T) {
	log.SetDefault(log.New(io.Discard))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup Suite")
}
