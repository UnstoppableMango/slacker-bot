package backup

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/unmango/go/fopt"
	"github.com/unstoppablemango/ihfs"
)

type Option func(Backup) Backup

type Backup struct {
	Id      uuid.UUID    `json:"id" yaml:"id"`
	GuildId snowflake.ID `json:"guild_id" yaml:"guild_id"`
}

func New(guildId snowflake.ID, options ...Option) Backup {
	b := Backup{
		Id:      uuid.New(),
		GuildId: guildId,
	}

	return fopt.WithAll(b, options)
}

func (b Backup) Write(ctx context.Context, r rest.Rest, fsys fs.FS) error {
	log := log.FromContext(ctx).With("guildId", b.GuildId)
	rg, err := r.GetGuild(b.GuildId, false, rest.WithCtx(ctx))
	if err != nil {
		return fmt.Errorf("get guild: %w", err)
	}

	log.Info("Got guild", "name", rg.Name)
	if err = tryPersist(ctx, rg.IconURL(), fsys); err != nil {
		return fmt.Errorf("persisting icon: %w", err)
	}
	if err = tryPersist(ctx, rg.BannerURL(), fsys); err != nil {
		return fmt.Errorf("persisting banner: %w", err)
	}
	if err = tryPersist(ctx, rg.SplashURL(), fsys); err != nil {
		return fmt.Errorf("persisting splash: %w", err)
	}

	members, err := r.GetMembers(b.GuildId, 10, snowflake.ID(0))
	if err != nil {
		return fmt.Errorf("get members: %w", err)
	}

	log.Info("Got members", "count", len(members))
	for _, m := range members {
		if err = tryPersist(ctx, m.User.AvatarURL(), fsys); err != nil {
			return fmt.Errorf("persisting avatar for user %s: %w", m.User.Username, err)
		}
		if err = tryPersist(ctx, m.User.BannerURL(), fsys); err != nil {
			return fmt.Errorf("persisting banner for user %s: %w", m.User.Username, err)
		}
		if err = tryPersist(ctx, m.User.AvatarDecorationURL(), fsys); err != nil {
			return fmt.Errorf("persisting display avatar for user %s: %w", m.User.Username, err)
		}
		if err = tryPersist(ctx, new(m.User.DefaultAvatarURL()), fsys); err != nil {
			return fmt.Errorf("persisting default avatar for user %s: %w", m.User.Username, err)
		}
		if err = tryPersist(ctx, m.User.AvatarURL(), fsys); err != nil {
			return fmt.Errorf("persisting avatar for user %s: %w", m.User.Username, err)
		}
	}

	return nil
}

func tryPersist(ctx context.Context, u *string, fsys fs.FS) error {
	log := log.FromContext(ctx).With("url", u)
	if u == nil {
		log.Info("No url, skipping")
		return nil
	}

	url, err := url.Parse(*u)
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	return persist(ctx, url, fsys)
}

func persist(ctx context.Context, u *url.URL, fsys fs.FS) error {
	log := log.FromContext(ctx).With("url", u)

	log.Info("Reading")
	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("fetching content: %w", err)
	}
	defer res.Body.Close()

	path := strings.TrimLeft(u.Path, "/")
	dir := filepath.Dir(path)

	log.Infof("Creating directory: %s", dir)
	if err = ihfs.MkdirAll(fsys, dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	log.Infof("Writing file: %s", path)
	return ihfs.WriteReader(fsys, path, res.Body, os.ModePerm)
}
