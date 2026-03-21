package backup

import (
	"io/fs"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "github.com/unstoppablemango/slacker-bot/gen/pb/dev/unmango/discord/backup/v1alpha1"
)

// testFS is a minimal in-memory FS satisfying ihfs.WriteFileFS and ihfs.MkdirAllFS.
type testFS struct {
	files map[string][]byte
	dirs  map[string]bool
}

func newTestFS() *testFS {
	return &testFS{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

func (f *testFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrInvalid
}

func (f *testFS) WriteFile(name string, data []byte, _ fs.FileMode) error {
	f.files[name] = data
	return nil
}

func (f *testFS) MkdirAll(name string, _ fs.FileMode) error {
	f.dirs[name] = true
	return nil
}

func testBackup(name string) *pb.ServerBackup {
	return (&pb.ServerBackup_builder{
		Guild: (&pb.Guild_builder{Name: &name}).Build(),
	}).Build()
}

var _ = Describe("Writer", func() {
	var (
		tfs    *testFS
		writer *Writer
	)

	BeforeEach(func() {
		tfs = newTestFS()
		writer = NewWriter(tfs)
	})

	Describe("Write", func() {
		It("writes the proto file named after the guild", func() {
			backup := testBackup("MyGuild")

			Expect(writer.Write(backup)).To(Succeed())
			Expect(tfs.files).To(HaveKey("MyGuild.binpb"))
		})

		It("writes proto content to the file", func() {
			backup := testBackup("MyGuild")

			Expect(writer.Write(backup)).To(Succeed())
			Expect(string(tfs.files["MyGuild.binpb"])).To(Equal(backup.String()))
		})

		It("returns an error when the FS write fails", func() {
			tfs.files = nil // causes WriteFile to panic/nil map — use a failing FS instead
			failFS := &failingWriteFS{}
			w := NewWriter(failFS)

			err := w.Write(testBackup("MyGuild"))
			Expect(err).To(HaveOccurred())
		})

		Context("with URLs in the backup", func() {
			var server *httptest.Server

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("asset-content"))
				}))
				DeferCleanup(server.Close)
			})

			It("fetches and writes each URL to the FS", func() {
				url := server.URL + "/avatars/123/abc.png"
				backup := (&pb.ServerBackup_builder{
					Guild: (&pb.Guild_builder{
						Name:    strPtr("Guild"),
						IconUrl: strPtr(url),
					}).Build(),
				}).Build()

				Expect(writer.Write(backup)).To(Succeed())
				Expect(tfs.files).To(HaveKey("avatars/123/abc.png"))
				Expect(string(tfs.files["avatars/123/abc.png"])).To(Equal("asset-content"))
			})

			It("creates the directory for each URL", func() {
				url := server.URL + "/avatars/123/abc.png"
				backup := (&pb.ServerBackup_builder{
					Guild: (&pb.Guild_builder{
						Name:    strPtr("Guild"),
						IconUrl: strPtr(url),
					}).Build(),
				}).Build()

				Expect(writer.Write(backup)).To(Succeed())
				Expect(tfs.dirs).To(HaveKey("avatars/123"))
			})

			It("skips empty URL fields", func() {
				backup := testBackup("Guild")

				Expect(writer.Write(backup)).To(Succeed())
				// Only the proto file should be written
				Expect(tfs.files).To(HaveLen(1))
			})

			It("returns an error when an HTTP fetch fails", func() {
				badURL := "http://127.0.0.1:0/nope" // nothing listening
				backup := (&pb.ServerBackup_builder{
					Guild: (&pb.Guild_builder{
						Name:    strPtr("Guild"),
						IconUrl: strPtr(badURL),
					}).Build(),
				}).Build()

				Expect(writer.Write(backup)).NotTo(Succeed())
			})
		})

		Context("URL collection", func() {
			var server *httptest.Server

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("x"))
				}))
				DeferCleanup(server.Close)
			})

			DescribeTable("fetches asset URLs from all backup fields",
				func(build func(url string) *pb.ServerBackup, path string) {
					url := server.URL + "/" + path
					Expect(writer.Write(build(url))).To(Succeed())
					Expect(tfs.files).To(HaveKey(path))
				},
				Entry("guild icon", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g"), IconUrl: strPtr(u)}).Build(),
					}).Build()
				}, "guild/icon.png"),
				Entry("guild banner", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g"), BannerUrl: strPtr(u)}).Build(),
					}).Build()
				}, "guild/banner.png"),
				Entry("guild splash", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g"), SplashUrl: strPtr(u)}).Build(),
					}).Build()
				}, "guild/splash.png"),
				Entry("guild discovery splash", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g"), DiscoverySplashUrl: strPtr(u)}).Build(),
					}).Build()
				}, "guild/discovery.png"),
				Entry("role icon", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g")}).Build(),
						Roles: []*pb.Role{(&pb.Role_builder{IconUrl: strPtr(u)}).Build()},
					}).Build()
				}, "roles/icon.png"),
				Entry("user avatar", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild: (&pb.Guild_builder{Name: strPtr("g")}).Build(),
						Users: []*pb.User{(&pb.User_builder{AvatarUrl: strPtr(u)}).Build()},
					}).Build()
				}, "avatars/user.png"),
				Entry("member guild avatar", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild:   (&pb.Guild_builder{Name: strPtr("g")}).Build(),
						Members: []*pb.Member{(&pb.Member_builder{GuildAvatarUrl: strPtr(u)}).Build()},
					}).Build()
				}, "guilds/avatar.png"),
				Entry("webhook avatar", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild:    (&pb.Guild_builder{Name: strPtr("g")}).Build(),
						Webhooks: []*pb.Webhook{(&pb.Webhook_builder{AvatarUrl: strPtr(u)}).Build()},
					}).Build()
				}, "webhooks/avatar.png"),
				Entry("scheduled event image", func(u string) *pb.ServerBackup {
					return (&pb.ServerBackup_builder{
						Guild:           (&pb.Guild_builder{Name: strPtr("g")}).Build(),
						ScheduledEvents: []*pb.ScheduledEvent{(&pb.ScheduledEvent_builder{ImageUrl: strPtr(u)}).Build()},
					}).Build()
				}, "events/cover.png"),
			)
		})
	})
})

var _ = Describe("WithHTTPClient", func() {
	It("sets the HTTP client on the writer", func() {
		client := &http.Client{}
		w := NewWriter(newTestFS(), WithHTTPClient(client))
		Expect(w.http).To(BeIdenticalTo(client))
	})
})

// failingWriteFS is an ihfs.FS whose WriteFile always returns an error.
type failingWriteFS struct{}

func (f *failingWriteFS) Open(_ string) (fs.File, error) { return nil, fs.ErrInvalid }
func (f *failingWriteFS) WriteFile(_ string, _ []byte, _ fs.FileMode) error {
	return fs.ErrInvalid
}
func (f *failingWriteFS) MkdirAll(_ string, _ fs.FileMode) error { return nil }

func strPtr(s string) *string { return &s }

var _ = Describe("backupURLs", func() {
	It("returns empty strings for unset fields", func() {
		urls := backupURLs(testBackup("g"))
		Expect(urls).To(ContainElements(
			Equal(""), Equal(""), Equal(""), Equal(""),
		))
	})

	It("collects all non-empty URLs", func() {
		backup := (&pb.ServerBackup_builder{
			Guild: (&pb.Guild_builder{
				Name:    strPtr("g"),
				IconUrl: strPtr("https://example.com/icon.png"),
			}).Build(),
			Users: []*pb.User{
				(&pb.User_builder{AvatarUrl: strPtr("https://example.com/avatar.png")}).Build(),
			},
		}).Build()

		urls := backupURLs(backup)
		Expect(urls).To(ContainElements(
			"https://example.com/icon.png",
			"https://example.com/avatar.png",
		))
	})

	It("includes one entry per collection item", func() {
		backup := (&pb.ServerBackup_builder{
			Guild: (&pb.Guild_builder{Name: strPtr("g")}).Build(),
			Roles: []*pb.Role{
				(&pb.Role_builder{IconUrl: strPtr("https://example.com/r1.png")}).Build(),
				(&pb.Role_builder{IconUrl: strPtr("https://example.com/r2.png")}).Build(),
			},
		}).Build()

		urls := backupURLs(backup)
		Expect(urls).To(ContainElements(
			"https://example.com/r1.png",
			"https://example.com/r2.png",
		))
	})
})

var _ = Describe("Write (package-level)", func() {
	It("writes the proto file", func() {
		tfs := newTestFS()
		backup := testBackup("MyGuild")

		Expect(Write(backup, tfs)).To(Succeed())
		Expect(tfs.files).To(HaveKey("MyGuild.binpb"))
		Expect(string(tfs.files["MyGuild.binpb"])).To(Equal(backup.String()))
	})
})
