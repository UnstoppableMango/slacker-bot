package backup

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/unstoppablemango/slacker-bot/gen/pb/dev/unmango/discord/backup/v1alpha1"
)

func optID(id *snowflake.ID) *string {
	if id == nil {
		return nil
	}
	return new(id.String())
}

func Create(ctx context.Context, r rest.Rest, guildId snowflake.ID) (*pb.ServerBackup, error) {
	opts := []rest.RequestOpt{rest.WithCtx(ctx)}

	guild, err := r.GetGuild(guildId, false, opts...)
	if err != nil {
		return nil, err
	}

	channels, err := r.GetGuildChannels(guildId, opts...)
	if err != nil {
		return nil, err
	}

	// the `after` parameter defaults to 0, but disgo requires it
	// https://docs.discord.com/developers/resources/guild#list-guild-members
	members, err := r.GetMembers(guildId, 1000, snowflake.ID(0), opts...)
	if err != nil {
		return nil, err
	}

	webhooks, err := r.GetAllWebhooks(guildId, opts...)
	if err != nil {
		return nil, err
	}

	scheduledEvents, err := r.GetGuildScheduledEvents(guildId, false, opts...)
	if err != nil {
		return nil, err
	}

	autoModRules, err := r.GetAutoModerationRules(guildId, opts...)
	if err != nil {
		return nil, err
	}

	invites, err := r.GetGuildInvites(guildId, opts...)
	if err != nil {
		return nil, err
	}

	backup := &pb.ServerBackup_builder{
		Guild:           mapGuild(guild.Guild),
		Channels:        mapChannels(channels),
		Roles:           mapRoles(guild.Roles),
		Members:         mapMembers(members),
		Users:           mapUsers(members),
		Emojis:          mapEmojis(guild.Emojis),
		Stickers:        mapStickers(guild.Stickers),
		Webhooks:        mapWebhooks(webhooks),
		ScheduledEvents: mapScheduledEvents(scheduledEvents),
		AutoModRules:    mapAutoModRules(autoModRules),
		Invites:         mapInvites(invites),
	}

	return backup.Build(), nil
}

func mapGuild(g discord.Guild) *pb.Guild {
	features := make([]string, len(g.Features))
	for i, f := range g.Features {
		features[i] = string(f)
	}

	b := &pb.Guild_builder{
		Id:                       new(g.ID.String()),
		Name:                     new(g.Name),
		Description:              g.Description,
		IconUrl:                  g.IconURL(),
		BannerUrl:                g.BannerURL(),
		SplashUrl:                g.SplashURL(),
		DiscoverySplashUrl:       g.DiscoverySplashURL(),
		OwnerId:                  new(g.OwnerID.String()),
		AfkChannelId:             optID(g.AfkChannelID),
		AfkTimeout:               new(int32(g.AfkTimeout)),
		SystemChannelId:          optID(g.SystemChannelID),
		SystemChannelFlags:       new(int64(g.SystemChannelFlags)),
		RulesChannelId:           optID(g.RulesChannelID),
		PublicUpdatesChannelId:   optID(g.PublicUpdatesChannelID),
		SafetyAlertsChannelId:    optID(g.SafetyAlertsChannelID),
		MaxMembers:               new(int32(g.MaxMembers)),
		MaxVideoChannelUsers:     new(int32(g.MaxVideoChannelUsers)),
		WidgetEnabled:            new(g.WidgetEnabled),
		ApproximateMemberCount:   new(int32(g.ApproximateMemberCount)),
		ApproximatePresenceCount: new(int32(g.ApproximatePresenceCount)),
		PreferredLocale:          new(g.PreferredLocale),
		PremiumSubscriptionCount: new(int32(g.PremiumSubscriptionCount)),
		VanityUrlCode:            g.VanityURLCode,
		JoinedAt:                 timestamppb.New(g.JoinedAt),
		Features:                 features,

		// Discord's enumerations start at 0, but protobuf's "unspecified" convention means we need to shift them to start at 1
		VerificationLevel:           new(pb.VerificationLevel(int32(g.VerificationLevel) + 1)),
		ExplicitContentFilter:       new(pb.ExplicitContentFilter(int32(g.ExplicitContentFilter) + 1)),
		DefaultMessageNotifications: new(pb.DefaultMessageNotifications(int32(g.DefaultMessageNotifications) + 1)),
		NsfwLevel:                   new(pb.NsfwLevel(int32(g.NSFWLevel) + 1)),
		PremiumTier:                 new(pb.PremiumTier(int32(g.PremiumTier) + 1)),
	}

	if g.MaxPresences != nil {
		b.MaxPresences = new(int32(*g.MaxPresences))
	}

	if g.WidgetChannelID != 0 {
		b.WidgetChannelId = new(g.WidgetChannelID.String())
	}

	return b.Build()
}

func mapChannels(channels []discord.GuildChannel) []*pb.Channel {
	result := make([]*pb.Channel, len(channels))
	for i, ch := range channels {
		result[i] = mapChannel(ch)
	}
	return result
}

func mapChannel(ch discord.GuildChannel) *pb.Channel {
	b := &pb.Channel_builder{
		Id:                   new(ch.ID().String()),
		Name:                 new(ch.Name()),
		Type:                 new(mapChannelType(ch.Type())),
		Position:             new(int32(ch.Position())),
		ParentId:             optID(ch.ParentID()),
		PermissionOverwrites: mapPermissionOverwrites(ch.PermissionOverwrites()),
	}

	if mc, ok := ch.(discord.GuildMessageChannel); ok {
		b.Topic = mc.Topic()
		b.Nsfw = new(mc.NSFW())
		b.RateLimitPerUser = new(int32(mc.RateLimitPerUser()))
		b.DefaultAutoArchiveDuration = new(int32(mc.DefaultAutoArchiveDuration()))
	}
	if ac, ok := ch.(discord.GuildAudioChannel); ok {
		b.Bitrate = new(int32(ac.Bitrate()))
	}

	switch c := ch.(type) {
	case discord.GuildThread:
		b.ThreadMetadata = mapThreadMetadata(c.ThreadMetadata)
	case discord.GuildVoiceChannel:
		b.UserLimit = new(int32(c.UserLimit))
		vqm := pb.VideoQualityMode(int32(c.VideoQualityMode))
		b.VideoQualityMode = &vqm
	case discord.GuildStageVoiceChannel:
		vqm := pb.VideoQualityMode(int32(c.VideoQualityMode))
		b.VideoQualityMode = &vqm
	case discord.GuildForumChannel:
		b.Flags = new(int64(c.Flags))
		fl := pb.ForumLayout(int32(c.DefaultForumLayout) + 1)
		b.DefaultForumLayout = &fl
		if c.DefaultSortOrder != nil {
			so := pb.SortOrder(int32(*c.DefaultSortOrder) + 1)
			b.DefaultSortOrder = &so
		}
		b.AvailableTags = mapForumTags(c.AvailableTags)
		if c.DefaultReactionEmoji != nil {
			b.DefaultReactionEmoji = mapDefaultReaction(c.DefaultReactionEmoji)
		}
		b.DefaultThreadRateLimitPerUser = new(int32(c.DefaultThreadRateLimitPerUser))
	}

	return b.Build()
}

func mapChannelType(ct discord.ChannelType) pb.ChannelType {
	// TODO: document why this is necessary
	v := int32(ct)
	if v >= 0 && v <= 5 {
		return pb.ChannelType(v + 1)
	}
	return pb.ChannelType(v)
}

func mapPermissionOverwrites(overwrites discord.PermissionOverwrites) []*pb.PermissionOverwrite {
	if overwrites == nil {
		return nil
	}

	result := make([]*pb.PermissionOverwrite, len(overwrites))
	for i, ow := range overwrites {
		b := &pb.PermissionOverwrite_builder{
			Id:   new(ow.ID().String()),
			Type: new(pb.OverwriteType(int32(ow.Type()) + 1)),
		}

		switch o := ow.(type) {
		case discord.RolePermissionOverwrite:
			b.Allow = new(int64(o.Allow))
			b.Deny = new(int64(o.Deny))
		case discord.MemberPermissionOverwrite:
			b.Allow = new(int64(o.Allow))
			b.Deny = new(int64(o.Deny))
		}
		result[i] = b.Build()
	}

	return result
}

func mapThreadMetadata(tm discord.ThreadMetadata) *pb.ThreadMetadata {
	b := &pb.ThreadMetadata_builder{
		Archived:            new(tm.Archived),
		AutoArchiveDuration: new(int32(tm.AutoArchiveDuration)),
		ArchiveTimestamp:    timestamppb.New(tm.ArchiveTimestamp),
		Locked:              new(tm.Locked),
		Invitable:           new(tm.Invitable),
		CreateTimestamp:     timestamppb.New(tm.CreateTimestamp),
	}

	return b.Build()
}

func mapForumTags(tags []discord.ChannelTag) []*pb.ForumTag {
	result := make([]*pb.ForumTag, len(tags))
	for i, t := range tags {
		b := &pb.ForumTag_builder{
			Id:        new(t.ID.String()),
			Name:      new(t.Name),
			Moderated: new(t.Moderated),
			EmojiId:   optID(t.EmojiID),
			EmojiName: t.EmojiName,
		}

		result[i] = b.Build()
	}

	return result
}

func mapDefaultReaction(dr *discord.DefaultReactionEmoji) *pb.DefaultReaction {
	b := &pb.DefaultReaction_builder{
		EmojiId:   optID(dr.EmojiID),
		EmojiName: dr.EmojiName,
	}

	return b.Build()
}

func mapRoles(roles []discord.Role) []*pb.Role {
	result := make([]*pb.Role, len(roles))
	for i, r := range roles {
		result[i] = mapRole(r)
	}
	return result
}

func mapRole(r discord.Role) *pb.Role {
	b := &pb.Role_builder{
		Id:           new(r.ID.String()),
		Name:         new(r.Name),
		Color:        new(uint32(r.Color)),
		Hoist:        new(r.Hoist),
		IconUrl:      r.IconURL(),
		UnicodeEmoji: r.Emoji,
		Position:     new(int32(r.Position)),
		Permissions:  new(int64(r.Permissions)),
		Managed:      new(r.Managed),
		Mentionable:  new(r.Mentionable),
		Flags:        new(int64(r.Flags)),
	}

	if r.Tags != nil {
		b.Tags = mapRoleTag(r.Tags)
	}

	return b.Build()
}

func mapRoleTag(t *discord.RoleTag) *pb.RoleTag {
	b := &pb.RoleTag_builder{
		BotId:                 optID(t.BotID),
		IntegrationId:         optID(t.IntegrationID),
		SubscriptionListingId: optID(t.SubscriptionListingID),
		PremiumSubscriber:     new(t.PremiumSubscriber),
		AvailableForPurchase:  new(t.AvailableForPurchase),
		GuildConnections:      new(t.GuildConnections),
	}

	return b.Build()
}

func mapMembers(members []discord.Member) []*pb.Member {
	result := make([]*pb.Member, len(members))
	for i, m := range members {
		result[i] = mapMember(m)
	}
	return result
}

func mapMember(m discord.Member) *pb.Member {
	roleIds := make([]string, len(m.RoleIDs))
	for i, r := range m.RoleIDs {
		roleIds[i] = r.String()
	}

	b := &pb.Member_builder{
		UserId:         new(m.User.ID.String()),
		Nickname:       m.Nick,
		GuildAvatarUrl: m.AvatarURL(),
		RoleIds:        roleIds,
		Deaf:           new(m.Deaf),
		Mute:           new(m.Mute),
		Pending:        new(m.Pending),
		Flags:          new(int64(m.Flags)),
	}

	if m.JoinedAt != nil {
		b.JoinedAt = timestamppb.New(*m.JoinedAt)
	}
	if m.PremiumSince != nil {
		b.PremiumSince = timestamppb.New(*m.PremiumSince)
	}
	if m.CommunicationDisabledUntil != nil {
		b.CommunicationDisabledUntil = timestamppb.New(*m.CommunicationDisabledUntil)
	}

	return b.Build()
}

func mapUsers(members []discord.Member) []*pb.User {
	seen := make(map[snowflake.ID]bool)
	result := make([]*pb.User, 0, len(members))
	for _, m := range members {
		if seen[m.User.ID] {
			continue
		}
		seen[m.User.ID] = true
		result = append(result, mapUser(m.User))
	}

	return result
}

func mapUser(u discord.User) *pb.User {
	b := &pb.User_builder{
		Id:            new(u.ID.String()),
		Username:      new(u.Username),
		Discriminator: new(u.Discriminator),
		GlobalName:    u.GlobalName,
		AvatarUrl:     u.AvatarURL(),
		Bot:           new(u.Bot),
		System:        new(u.System),
	}

	return b.Build()
}

func mapEmojis(emojis []discord.Emoji) []*pb.Emoji {
	result := make([]*pb.Emoji, len(emojis))
	for i, e := range emojis {
		result[i] = mapEmoji(e)
	}
	return result
}

func mapEmoji(e discord.Emoji) *pb.Emoji {
	roleIds := make([]string, len(e.Roles))
	for i, r := range e.Roles {
		roleIds[i] = r.String()
	}

	b := &pb.Emoji_builder{
		Id:            new(e.ID.String()),
		Name:          new(e.Name),
		RoleIds:       roleIds,
		RequireColons: new(e.RequireColons),
		Managed:       new(e.Managed),
		Animated:      new(e.Animated),
		Available:     new(e.Available),
	}

	if e.Creator != nil {
		b.UserId = new(e.Creator.ID.String())
	}

	return b.Build()
}

func mapStickers(stickers []discord.Sticker) []*pb.Sticker {
	result := make([]*pb.Sticker, len(stickers))
	for i, s := range stickers {
		result[i] = mapSticker(s)
	}
	return result
}

func mapSticker(s discord.Sticker) *pb.Sticker {
	b := &pb.Sticker_builder{
		Id:          new(s.ID.String()),
		Name:        new(s.Name),
		Description: new(s.Description),
		Tags:        new(s.Tags),
		FormatType:  new(pb.StickerFormatType(int32(s.FormatType))),
		PackId:      optID(s.PackID),
		GuildId:     optID(s.GuildID),
		Available:   s.Available,
	}

	if s.User != nil {
		b.UserId = new(s.User.ID.String())
	}
	if s.SortValue != nil {
		b.SortValue = new(int32(*s.SortValue))
	}

	return b.Build()
}

func mapWebhooks(webhooks []discord.Webhook) []*pb.Webhook {
	result := make([]*pb.Webhook, len(webhooks))
	for i, w := range webhooks {
		result[i] = mapWebhook(w)
	}
	return result
}

func mapWebhook(w discord.Webhook) *pb.Webhook {
	b := &pb.Webhook_builder{
		Id:        new(w.ID().String()),
		Name:      new(w.Name()),
		Type:      new(pb.WebhookType(int32(w.Type()))),
		AvatarUrl: w.AvatarURL(),
	}

	if wh, ok := w.(discord.IncomingWebhook); ok {
		b.GuildId = new(wh.GuildID.String())
		b.ChannelId = new(wh.ChannelID.String())
		b.UserId = new(wh.User.ID.String())

		if wh.Token != "" {
			b.Token = new(wh.Token)
		}
		if wh.ApplicationID != nil {
			b.ApplicationId = new(wh.ApplicationID.String())
		}
	}

	return b.Build()
}

func mapScheduledEvents(events []discord.GuildScheduledEvent) []*pb.ScheduledEvent {
	result := make([]*pb.ScheduledEvent, len(events))
	for i, e := range events {
		result[i] = mapScheduledEvent(e)
	}
	return result
}

func mapScheduledEvent(e discord.GuildScheduledEvent) *pb.ScheduledEvent {
	b := &pb.ScheduledEvent_builder{
		Id:                 new(e.ID.String()),
		GuildId:            new(e.GuildID.String()),
		ChannelId:          optID(e.ChannelID),
		CreatorId:          new(e.CreatorID.String()),
		Name:               new(e.Name),
		Description:        new(e.Description),
		ScheduledStartTime: timestamppb.New(e.ScheduledStartTime),
		PrivacyLevel:       new(pb.PrivacyLevel(int32(e.PrivacyLevel))),
		Status:             new(pb.ScheduledEventStatus(int32(e.Status))),
		EntityType:         new(pb.ScheduledEventEntityType(int32(e.EntityType))),
		EntityId:           optID(e.EntityID),
		UserCount:          new(int32(e.UserCount)),
		ImageUrl:           e.Image,
	}

	if e.ScheduledEndTime != nil {
		b.ScheduledEndTime = timestamppb.New(*e.ScheduledEndTime)
	}

	if e.EntityMetaData != nil {
		mb := &pb.ScheduledEventEntityMetadata_builder{
			Location: new(e.EntityMetaData.Location),
		}
		b.EntityMetadata = mb.Build()
	}

	return b.Build()
}

func mapAutoModRules(rules []discord.AutoModerationRule) []*pb.AutoModRule {
	result := make([]*pb.AutoModRule, len(rules))
	for i, r := range rules {
		result[i] = mapAutoModRule(r)
	}
	return result
}

func mapAutoModRule(r discord.AutoModerationRule) *pb.AutoModRule {
	tmb := &pb.AutoModTriggerMetadata_builder{
		KeywordFilter:                r.TriggerMetadata.KeywordFilter,
		RegexPatterns:                r.TriggerMetadata.RegexPatterns,
		Presets:                      make([]int32, len(r.TriggerMetadata.Presets)),
		AllowList:                    r.TriggerMetadata.AllowList,
		MentionTotalLimit:            new(int32(r.TriggerMetadata.MentionTotalLimit)),
		MentionRaidProtectionEnabled: new(r.TriggerMetadata.MentionRaidProtectionEnabled),
	}
	for i, p := range r.TriggerMetadata.Presets {
		tmb.Presets[i] = int32(p)
	}

	b := &pb.AutoModRule_builder{
		Id:               new(r.ID.String()),
		GuildId:          new(r.GuildID.String()),
		Name:             new(r.Name),
		CreatorId:        new(r.CreatorID.String()),
		EventType:        new(pb.AutoModEventType(int32(r.EventType))),
		TriggerType:      new(mapAutoModTriggerType(r.TriggerType)),
		TriggerMetadata:  tmb.Build(),
		Actions:          make([]*pb.AutoModAction, len(r.Actions)),
		Enabled:          new(r.Enabled),
		ExemptRoleIds:    make([]string, len(r.ExemptRoles)),
		ExemptChannelIds: make([]string, len(r.ExemptChannels)),
	}
	for i, er := range r.ExemptRoles {
		b.ExemptRoleIds[i] = er.String()
	}
	for i, ec := range r.ExemptChannels {
		b.ExemptChannelIds[i] = ec.String()
	}
	for i, a := range r.Actions {
		b.Actions[i] = mapAutoModAction(a)
	}

	return b.Build()
}

func mapAutoModTriggerType(t discord.AutoModerationTriggerType) pb.AutoModTriggerType {
	switch t {
	case discord.AutoModerationTriggerTypeKeyword:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD
	case discord.AutoModerationTriggerTypeSpam:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_SPAM
	case discord.AutoModerationTriggerTypeKeywordPresent:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD_PRESET
	case discord.AutoModerationTriggerTypeMentionSpam:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MENTION_SPAM
	case discord.AutoModerationTriggerTypeMemberProfile:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MEMBER_PROFILE
	default:
		return pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_UNSPECIFIED
	}
}

func mapAutoModAction(a discord.AutoModerationAction) *pb.AutoModAction {
	b := &pb.AutoModAction_builder{
		Type: new(pb.AutoModActionType(int32(a.Type))),
	}

	if a.Metadata != nil {
		b.ChannelId = new(a.Metadata.ChannelID.String())
		b.DurationSeconds = new(int32(a.Metadata.DurationSeconds))
		b.CustomMessage = a.Metadata.CustomMessage
	}

	return b.Build()
}

func mapInvites(invites []discord.ExtendedInvite) []*pb.Invite {
	result := make([]*pb.Invite, len(invites))
	for i, inv := range invites {
		result[i] = mapInvite(inv)
	}
	return result
}

func mapInvite(inv discord.ExtendedInvite) *pb.Invite {
	b := &pb.Invite_builder{
		Code:      new(inv.Code),
		MaxAge:    new(int32(inv.MaxAge)),
		MaxUses:   new(int32(inv.MaxUses)),
		Temporary: new(inv.Temporary),
		Uses:      new(int32(inv.Uses)),
		CreatedAt: timestamppb.New(inv.CreatedAt),
	}

	if inv.Channel != nil {
		b.ChannelId = new(inv.Channel.ID.String())
	}
	if inv.Inviter != nil {
		b.InviterId = new(inv.Inviter.ID.String())
	}

	return b.Build()
}
