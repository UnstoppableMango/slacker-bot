package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/ihfs/osfs"
	"github.com/unstoppablemango/slacker-bot/pkg/backup"
)

var logger = log.New(os.Stdout)

func main() {
	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithLogger(slog.New(logger)),
	)
	if err != nil {
		cli.Fail(err)
	}

	ctx := log.WithContext(context.Background(), logger)
	id := snowflake.GetEnv("DISCORD_GUILD_ID")
	fsys := osfs.New()
	b, err := backup.Create(ctx, client.Rest, id)

	if err := backup.Write(b, fsys); err != nil {
		cli.Fail(err)
	}
}
