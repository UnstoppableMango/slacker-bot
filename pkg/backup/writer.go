package backup

import (
	"bytes"
	"fmt"
	"os"

	"github.com/unstoppablemango/ihfs"
	pb "github.com/unstoppablemango/slacker-bot/gen/dev/unmango/discord/backup/v1alpha1"
	"google.golang.org/protobuf/proto"
)

type Writer struct {
	fs ihfs.FS
}

// NewWriter creates a new Writer writing to the given ihfs.FS.
func NewWriter(fs ihfs.FS) *Writer {
	return &Writer{fs: fs}
}

func (w *Writer) Write(b *pb.ServerBackup) error {
	data, err := proto.Marshal(b)
	if err != nil {
		return err
	}

	err = ihfs.WriteReader(w.fs,
		fmt.Sprintf("%s.binpb", b.GetGuild().GetName()),
		bytes.NewReader(data),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	return nil
}
