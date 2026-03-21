package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TarGzFS", func() {
	var (
		buf *bytes.Buffer
		tgz *TarGzFS
	)

	BeforeEach(func() {
		buf = &bytes.Buffer{}
		tgz = NewTarGzFS(buf)
	})

	Describe("WriteFile", func() {
		It("writes a tar entry with the correct name", func() {
			Expect(tgz.WriteFile("hello.txt", []byte("world"), fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			entry := firstTarEntry(buf)
			Expect(entry.Name).To(Equal("hello.txt"))
		})

		It("writes the file contents to the tar entry", func() {
			Expect(tgz.WriteFile("hello.txt", []byte("world"), fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			_, data := firstTarEntryWithData(buf)
			Expect(string(data)).To(Equal("world"))
		})

		It("sets the file mode on the tar header", func() {
			Expect(tgz.WriteFile("hello.txt", []byte("data"), 0o644)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			entry := firstTarEntry(buf)
			Expect(entry.Mode).To(Equal(int64(0o644)))
		})

		It("sets the correct size on the tar header", func() {
			data := []byte("hello world")
			Expect(tgz.WriteFile("f.txt", data, fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			entry := firstTarEntry(buf)
			Expect(entry.Size).To(Equal(int64(len(data))))
		})

		It("writes multiple files", func() {
			Expect(tgz.WriteFile("a.txt", []byte("aaa"), fs.ModePerm)).To(Succeed())
			Expect(tgz.WriteFile("b.txt", []byte("bbb"), fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			names := allTarEntryNames(buf)
			Expect(names).To(ConsistOf("a.txt", "b.txt"))
		})
	})

	Describe("MkdirAll", func() {
		It("is a no-op and returns nil", func() {
			Expect(tgz.MkdirAll("some/path", fs.ModePerm)).To(Succeed())
		})

		It("does not add any entries to the archive", func() {
			Expect(tgz.MkdirAll("some/path", fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			Expect(allTarEntryNames(buf)).To(BeEmpty())
		})
	})

	Describe("Close", func() {
		It("produces a valid gzip stream", func() {
			Expect(tgz.Close()).To(Succeed())

			_, err := gzip.NewReader(buf)
			Expect(err).NotTo(HaveOccurred())
		})

		It("produces a readable tar archive after writing", func() {
			Expect(tgz.WriteFile("x.txt", []byte("x"), fs.ModePerm)).To(Succeed())
			Expect(tgz.Close()).To(Succeed())

			Expect(allTarEntryNames(buf)).To(ContainElement("x.txt"))
		})
	})

	Describe("Open", func() {
		It("always returns fs.ErrInvalid", func() {
			_, err := tgz.Open("anything")
			Expect(err).To(MatchError(fs.ErrInvalid))
		})
	})
})

// helpers

func openTar(buf *bytes.Buffer) *tar.Reader {
	gz, err := gzip.NewReader(bytes.NewReader(buf.Bytes()))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return tar.NewReader(gz)
}

func firstTarEntry(buf *bytes.Buffer) *tar.Header {
	hdr, _, _ := nextTarEntry(openTar(buf))
	return hdr
}

func firstTarEntryWithData(buf *bytes.Buffer) (*tar.Header, []byte) {
	hdr, data, _ := nextTarEntry(openTar(buf))
	return hdr, data
}

func nextTarEntry(tr *tar.Reader) (*tar.Header, []byte, error) {
	hdr, err := tr.Next()
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(tr)
	return hdr, data, err
}

func allTarEntryNames(buf *bytes.Buffer) []string {
	tr := openTar(buf)
	var names []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		names = append(names, hdr.Name)
	}
	return names
}
