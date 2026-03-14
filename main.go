package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/unmango/go/cli"
)

func main() {
	log := log.New(os.Stdout)
	client, err := disgo.New(os.Getenv("DISCORD_BOT_TOKEN"),
		bot.WithLogger(slog.New(log)),
	)
	if err != nil {
		cli.Fail(err)
	}

	ctx := context.Background()
	id := snowflake.GetEnv("DISCORD_GUILD_ID")
	log.Infof("GUILD_ID: %s", id)

	rg, err := client.Rest.GetGuild(id, false, rest.WithCtx(ctx))
	if err != nil {
		cli.Fail("GetGuild:", err)
	}

	log.Info("GetGuild",
		"guild", rg.Name,
		"url", *rg.IconURL(),
	)
}
