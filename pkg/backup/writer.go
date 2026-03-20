package backup

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/disgoorg/disgo/rest"
	"github.com/unmango/go/fopt"
	"github.com/unstoppablemango/ihfs"
	pb "github.com/unstoppablemango/slacker-bot/gen/pb/dev/unmango/discord/backup/v1alpha1"
)

type WriterOpt func(*Writer)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type Writer struct {
	fs   ihfs.FS
	rest rest.Rest
	http HttpClient
}

// NewWriter creates a new Writer making requests using
// the given [rest.Rest] and writing to the given ihfs.FS.
func NewWriter(fs ihfs.FS, rest rest.Rest, options ...WriterOpt) *Writer {
	w := &Writer{
		fs:   fs,
		rest: rest,
		http: http.DefaultClient,
	}
	fopt.ApplyAll(w, options)

	return w
}

// Write writes the given backup to the ihfs.FS of the Writer.
//
// It will write the backup as a text protobuf with the name of the guild as the filename.
// It will follow all external links in the backup using the Writer's HTTP client
// and write those to the ihfs.FS as well.
func (w *Writer) Write(b *pb.ServerBackup) error {
	err := ihfs.WriteReader(w.fs,
		fmt.Sprintf("%s.binpb", b.GetGuild().GetName()),
		strings.NewReader(b.String()),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	for _, u := range backupURLs(b) {
		if err := w.persistURL(u); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) persistURL(raw string) error {
	if raw == "" {
		return nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing url %q: %w", raw, err)
	}

	res, err := w.http.Get(u.String())
	if err != nil {
		return fmt.Errorf("fetching %q: %w", u, err)
	}
	defer res.Body.Close()

	path := strings.TrimLeft(u.Path, "/")
	if err = ihfs.MkdirAll(w.fs, filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("creating directory for %q: %w", path, err)
	}

	return ihfs.WriteReader(w.fs, path, res.Body, os.ModePerm)
}

func backupURLs(b *pb.ServerBackup) []string {
	g := b.GetGuild()
	urls := []string{
		g.GetIconUrl(),
		g.GetBannerUrl(),
		g.GetSplashUrl(),
		g.GetDiscoverySplashUrl(),
	}

	for _, r := range b.GetRoles() {
		urls = append(urls, r.GetIconUrl())
	}

	for _, u := range b.GetUsers() {
		urls = append(urls, u.GetAvatarUrl())
	}

	for _, m := range b.GetMembers() {
		urls = append(urls, m.GetGuildAvatarUrl())
	}

	for _, wh := range b.GetWebhooks() {
		urls = append(urls, wh.GetAvatarUrl())
	}

	for _, e := range b.GetScheduledEvents() {
		urls = append(urls, e.GetImageUrl())
	}

	return urls
}

func WithHTTPClient(client *http.Client) WriterOpt {
	return func(w *Writer) {
		w.http = client
	}
}
