package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"charm.land/log/v2"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/ihfs/osfs"
	"github.com/unstoppablemango/ihfs/tarfs"
	"github.com/unstoppablemango/slacker-bot/pkg/backup"
)

var logger = log.New(os.Stdout)

func main() {
	tarPath := flag.String("tar", "", "write backup to a tar archive at the given path")
	flag.Parse()

	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithLogger(slog.New(logger)),
	)
	if err != nil {
		cli.Fail(err)
	}

	ctx := log.WithContext(context.Background(), logger)
	id := snowflake.GetEnv("DISCORD_GUILD_ID")

	b, err := backup.Create(ctx, client.Rest, id)
	if err != nil {
		cli.Fail(err)
	}

	if *tarPath != "" {
		f, err := os.Create(*tarPath)
		if err != nil {
			cli.Fail(err)
		}
		tw := tarfs.NewWriter(f)
		if err := backup.Write(ctx, b, tw); err != nil {
			cli.Fail(err)
		}
		if err := tw.Close(); err != nil {
			cli.Fail(err)
		}
		if err := f.Close(); err != nil {
			cli.Fail(err)
		}
		return
	}

	if err := backup.Write(ctx, b, osfs.New()); err != nil {
		cli.Fail(err)
	}
}
