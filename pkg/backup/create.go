package backup

import (
	"context"
	"io/fs"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/unstoppablemango/slacker-bot/gen/dev/unmango/discord/backup/v1alpha1"
)

func ptr[T any](v T) *T {
	return &v
}

func optID(id *snowflake.ID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

func Create(ctx context.Context, r rest.Rest, guildId snowflake.ID, fsys fs.FS) (*pb.ServerBackup, error) {
	opts := []rest.RequestOpt{rest.WithCtx(ctx)}

	g, err := r.GetGuild(guildId, true, opts...)
	if err != nil {
		return nil, err
	}

	channels, err := r.GetGuildChannels(guildId, opts...)
	if err != nil {
		return nil, err
	}

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

	return (&pb.ServerBackup_builder{
		Guild:           mapGuild(g.Guild),
		Channels:        mapChannels(channels),
		Roles:           mapRoles(g.Roles),
		Members:         mapMembers(members),
		Users:           mapUsers(members),
		Emojis:          mapEmojis(g.Emojis),
		Stickers:        mapStickers(g.Stickers),
		Webhooks:        mapWebhooks(webhooks),
		ScheduledEvents: mapScheduledEvents(scheduledEvents),
		AutoModRules:    mapAutoModRules(autoModRules),
		Invites:         mapInvites(invites),
	}).Build(), nil
}

func mapGuild(g discord.Guild) *pb.Guild {
	id := g.ID.String()
	ownerID := g.OwnerID.String()
	preferredLocale := g.PreferredLocale
	premiumCount := int32(g.PremiumSubscriptionCount)
	maxMembers := int32(g.MaxMembers)
	maxVideoUsers := int32(g.MaxVideoChannelUsers)
	approxMembers := int32(g.ApproximateMemberCount)
	approxPresence := int32(g.ApproximatePresenceCount)

	features := make([]string, len(g.Features))
	for i, f := range g.Features {
		features[i] = string(f)
	}

	b := &pb.Guild_builder{
		Id:                          &id,
		Name:                        &g.Name,
		Description:                 g.Description,
		IconUrl:                     g.IconURL(),
		BannerUrl:                   g.BannerURL(),
		SplashUrl:                   g.SplashURL(),
		DiscoverySplashUrl:          g.DiscoverySplashURL(),
		OwnerId:                     &ownerID,
		AfkChannelId:                optID(g.AfkChannelID),
		AfkTimeout:                  ptr(int32(g.AfkTimeout)),
		SystemChannelId:             optID(g.SystemChannelID),
		SystemChannelFlags:          ptr(int64(g.SystemChannelFlags)),
		RulesChannelId:              optID(g.RulesChannelID),
		PublicUpdatesChannelId:      optID(g.PublicUpdatesChannelID),
		SafetyAlertsChannelId:       optID(g.SafetyAlertsChannelID),
		MaxMembers:                  &maxMembers,
		MaxVideoChannelUsers:        &maxVideoUsers,
		WidgetEnabled:               ptr(g.WidgetEnabled),
		ApproximateMemberCount:      &approxMembers,
		ApproximatePresenceCount:    &approxPresence,
		PreferredLocale:             &preferredLocale,
		PremiumSubscriptionCount:    &premiumCount,
		VanityUrlCode:               g.VanityURLCode,
		JoinedAt:                    timestamppb.New(g.JoinedAt),
		Features:                    features,
		VerificationLevel:           ptr(pb.VerificationLevel(int32(g.VerificationLevel) + 1)),
		ExplicitContentFilter:       ptr(pb.ExplicitContentFilter(int32(g.ExplicitContentFilter) + 1)),
		DefaultMessageNotifications: ptr(pb.DefaultMessageNotifications(int32(g.DefaultMessageNotifications) + 1)),
		NsfwLevel:                   ptr(pb.NsfwLevel(int32(g.NSFWLevel) + 1)),
		PremiumTier:                 ptr(pb.PremiumTier(int32(g.PremiumTier) + 1)),
	}

	if g.MaxPresences != nil {
		b.MaxPresences = ptr(int32(*g.MaxPresences))
	}

	if g.WidgetChannelID != 0 {
		wid := g.WidgetChannelID.String()
		b.WidgetChannelId = &wid
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
	id := ch.ID().String()
	name := ch.Name()
	ct := mapChannelType(ch.Type())
	pos := int32(ch.Position())

	b := &pb.Channel_builder{
		Id:                   &id,
		Name:                 &name,
		Type:                 &ct,
		Position:             &pos,
		ParentId:             optID(ch.ParentID()),
		PermissionOverwrites: mapPermissionOverwrites(ch.PermissionOverwrites()),
	}

	if mc, ok := ch.(discord.GuildMessageChannel); ok {
		b.Topic = mc.Topic()
		b.Nsfw = ptr(mc.NSFW())
		b.RateLimitPerUser = ptr(int32(mc.RateLimitPerUser()))
		b.DefaultAutoArchiveDuration = ptr(int32(mc.DefaultAutoArchiveDuration()))
	}

	if ac, ok := ch.(discord.GuildAudioChannel); ok {
		b.Bitrate = ptr(int32(ac.Bitrate()))
	}

	switch c := ch.(type) {
	case discord.GuildThread:
		b.ThreadMetadata = mapThreadMetadata(c.ThreadMetadata)
	case discord.GuildVoiceChannel:
		b.UserLimit = ptr(int32(c.UserLimit))
		vqm := pb.VideoQualityMode(int32(c.VideoQualityMode))
		b.VideoQualityMode = &vqm
	case discord.GuildStageVoiceChannel:
		vqm := pb.VideoQualityMode(int32(c.VideoQualityMode))
		b.VideoQualityMode = &vqm
	case discord.GuildForumChannel:
		b.Flags = ptr(int64(c.Flags))
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
		b.DefaultThreadRateLimitPerUser = ptr(int32(c.DefaultThreadRateLimitPerUser))
	}

	return b.Build()
}

func mapChannelType(ct discord.ChannelType) pb.ChannelType {
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
		id := ow.ID().String()
		owType := pb.OverwriteType(int32(ow.Type()) + 1)
		b := &pb.PermissionOverwrite_builder{
			Id:   &id,
			Type: &owType,
		}
		switch o := ow.(type) {
		case discord.RolePermissionOverwrite:
			b.Allow = ptr(int64(o.Allow))
			b.Deny = ptr(int64(o.Deny))
		case discord.MemberPermissionOverwrite:
			b.Allow = ptr(int64(o.Allow))
			b.Deny = ptr(int64(o.Deny))
		}
		result[i] = b.Build()
	}
	return result
}

func mapThreadMetadata(tm discord.ThreadMetadata) *pb.ThreadMetadata {
	return (&pb.ThreadMetadata_builder{
		Archived:            ptr(tm.Archived),
		AutoArchiveDuration: ptr(int32(tm.AutoArchiveDuration)),
		ArchiveTimestamp:    timestamppb.New(tm.ArchiveTimestamp),
		Locked:              ptr(tm.Locked),
		Invitable:           ptr(tm.Invitable),
		CreateTimestamp:     timestamppb.New(tm.CreateTimestamp),
	}).Build()
}

func mapForumTags(tags []discord.ChannelTag) []*pb.ForumTag {
	result := make([]*pb.ForumTag, len(tags))
	for i, t := range tags {
		id := t.ID.String()
		result[i] = (&pb.ForumTag_builder{
			Id:        &id,
			Name:      &t.Name,
			Moderated: ptr(t.Moderated),
			EmojiId:   optID(t.EmojiID),
			EmojiName: t.EmojiName,
		}).Build()
	}
	return result
}

func mapDefaultReaction(dr *discord.DefaultReactionEmoji) *pb.DefaultReaction {
	return (&pb.DefaultReaction_builder{
		EmojiId:   optID(dr.EmojiID),
		EmojiName: dr.EmojiName,
	}).Build()
}

func mapRoles(roles []discord.Role) []*pb.Role {
	result := make([]*pb.Role, len(roles))
	for i, r := range roles {
		result[i] = mapRole(r)
	}
	return result
}

func mapRole(r discord.Role) *pb.Role {
	id := r.ID.String()
	name := r.Name
	color := uint32(r.Color)
	pos := int32(r.Position)
	perms := int64(r.Permissions)
	flags := int64(r.Flags)

	b := &pb.Role_builder{
		Id:           &id,
		Name:         &name,
		Color:        &color,
		Hoist:        ptr(r.Hoist),
		IconUrl:      r.IconURL(),
		UnicodeEmoji: r.Emoji,
		Position:     &pos,
		Permissions:  &perms,
		Managed:      ptr(r.Managed),
		Mentionable:  ptr(r.Mentionable),
		Flags:        &flags,
	}

	if r.Tags != nil {
		b.Tags = mapRoleTag(r.Tags)
	}

	return b.Build()
}

func mapRoleTag(t *discord.RoleTag) *pb.RoleTag {
	return (&pb.RoleTag_builder{
		BotId:                 optID(t.BotID),
		IntegrationId:         optID(t.IntegrationID),
		SubscriptionListingId: optID(t.SubscriptionListingID),
		PremiumSubscriber:     ptr(t.PremiumSubscriber),
		AvailableForPurchase:  ptr(t.AvailableForPurchase),
		GuildConnections:      ptr(t.GuildConnections),
	}).Build()
}

func mapMembers(members []discord.Member) []*pb.Member {
	result := make([]*pb.Member, len(members))
	for i, m := range members {
		result[i] = mapMember(m)
	}
	return result
}

func mapMember(m discord.Member) *pb.Member {
	userId := m.User.ID.String()
	roleIds := make([]string, len(m.RoleIDs))
	for i, r := range m.RoleIDs {
		roleIds[i] = r.String()
	}

	b := &pb.Member_builder{
		UserId:         &userId,
		Nickname:       m.Nick,
		GuildAvatarUrl: m.AvatarURL(),
		RoleIds:        roleIds,
		Deaf:           ptr(m.Deaf),
		Mute:           ptr(m.Mute),
		Pending:        ptr(m.Pending),
		Flags:          ptr(int64(m.Flags)),
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
	id := u.ID.String()
	username := u.Username
	discriminator := u.Discriminator
	return (&pb.User_builder{
		Id:            &id,
		Username:      &username,
		Discriminator: &discriminator,
		GlobalName:    u.GlobalName,
		AvatarUrl:     u.AvatarURL(),
		Bot:           ptr(u.Bot),
		System:        ptr(u.System),
	}).Build()
}

func mapEmojis(emojis []discord.Emoji) []*pb.Emoji {
	result := make([]*pb.Emoji, len(emojis))
	for i, e := range emojis {
		result[i] = mapEmoji(e)
	}
	return result
}

func mapEmoji(e discord.Emoji) *pb.Emoji {
	id := e.ID.String()
	name := e.Name
	roleIds := make([]string, len(e.Roles))
	for i, r := range e.Roles {
		roleIds[i] = r.String()
	}

	b := &pb.Emoji_builder{
		Id:            &id,
		Name:          &name,
		RoleIds:       roleIds,
		RequireColons: ptr(e.RequireColons),
		Managed:       ptr(e.Managed),
		Animated:      ptr(e.Animated),
		Available:     ptr(e.Available),
	}

	if e.Creator != nil {
		userId := e.Creator.ID.String()
		b.UserId = &userId
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
	id := s.ID.String()
	name := s.Name
	description := s.Description
	tags := s.Tags
	formatType := pb.StickerFormatType(int32(s.FormatType))

	b := &pb.Sticker_builder{
		Id:          &id,
		Name:        &name,
		Description: &description,
		Tags:        &tags,
		FormatType:  &formatType,
		PackId:      optID(s.PackID),
		GuildId:     optID(s.GuildID),
		Available:   s.Available,
	}

	if s.User != nil {
		userId := s.User.ID.String()
		b.UserId = &userId
	}

	if s.SortValue != nil {
		b.SortValue = ptr(int32(*s.SortValue))
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
	id := w.ID().String()
	name := w.Name()
	wType := pb.WebhookType(int32(w.Type()))

	b := &pb.Webhook_builder{
		Id:        &id,
		Name:      &name,
		Type:      &wType,
		AvatarUrl: w.AvatarURL(),
	}

	if wh, ok := w.(discord.IncomingWebhook); ok {
		guildId := wh.GuildID.String()
		channelId := wh.ChannelID.String()
		userId := wh.User.ID.String()
		b.GuildId = &guildId
		b.ChannelId = &channelId
		b.UserId = &userId
		if wh.Token != "" {
			b.Token = &wh.Token
		}
		if wh.ApplicationID != nil {
			appId := wh.ApplicationID.String()
			b.ApplicationId = &appId
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
	id := e.ID.String()
	guildId := e.GuildID.String()
	creatorId := e.CreatorID.String()
	name := e.Name
	description := e.Description
	privacyLevel := pb.PrivacyLevel(int32(e.PrivacyLevel))
	status := pb.ScheduledEventStatus(int32(e.Status))
	entityType := pb.ScheduledEventEntityType(int32(e.EntityType))
	userCount := int32(e.UserCount)

	b := &pb.ScheduledEvent_builder{
		Id:                 &id,
		GuildId:            &guildId,
		ChannelId:          optID(e.ChannelID),
		CreatorId:          &creatorId,
		Name:               &name,
		Description:        &description,
		ScheduledStartTime: timestamppb.New(e.ScheduledStartTime),
		PrivacyLevel:       &privacyLevel,
		Status:             &status,
		EntityType:         &entityType,
		EntityId:           optID(e.EntityID),
		UserCount:          &userCount,
		ImageUrl:           e.Image,
	}

	if e.ScheduledEndTime != nil {
		b.ScheduledEndTime = timestamppb.New(*e.ScheduledEndTime)
	}

	if e.EntityMetaData != nil {
		b.EntityMetadata = (&pb.ScheduledEventEntityMetadata_builder{
			Location: ptr(e.EntityMetaData.Location),
		}).Build()
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
	id := r.ID.String()
	guildId := r.GuildID.String()
	name := r.Name
	creatorId := r.CreatorID.String()
	eventType := pb.AutoModEventType(int32(r.EventType))
	triggerType := mapAutoModTriggerType(r.TriggerType)

	exemptRoles := make([]string, len(r.ExemptRoles))
	for i, er := range r.ExemptRoles {
		exemptRoles[i] = er.String()
	}
	exemptChannels := make([]string, len(r.ExemptChannels))
	for i, ec := range r.ExemptChannels {
		exemptChannels[i] = ec.String()
	}
	actions := make([]*pb.AutoModAction, len(r.Actions))
	for i, a := range r.Actions {
		actions[i] = mapAutoModAction(a)
	}
	presets := make([]int32, len(r.TriggerMetadata.Presets))
	for i, p := range r.TriggerMetadata.Presets {
		presets[i] = int32(p)
	}

	return (&pb.AutoModRule_builder{
		Id:          &id,
		GuildId:     &guildId,
		Name:        &name,
		CreatorId:   &creatorId,
		EventType:   &eventType,
		TriggerType: &triggerType,
		TriggerMetadata: (&pb.AutoModTriggerMetadata_builder{
			KeywordFilter:                r.TriggerMetadata.KeywordFilter,
			RegexPatterns:                r.TriggerMetadata.RegexPatterns,
			Presets:                      presets,
			AllowList:                    r.TriggerMetadata.AllowList,
			MentionTotalLimit:            ptr(int32(r.TriggerMetadata.MentionTotalLimit)),
			MentionRaidProtectionEnabled: ptr(r.TriggerMetadata.MentionRaidProtectionEnabled),
		}).Build(),
		Actions:          actions,
		Enabled:          ptr(r.Enabled),
		ExemptRoleIds:    exemptRoles,
		ExemptChannelIds: exemptChannels,
	}).Build()
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
	actionType := pb.AutoModActionType(int32(a.Type))
	b := &pb.AutoModAction_builder{
		Type: &actionType,
	}
	if a.Metadata != nil {
		channelId := a.Metadata.ChannelID.String()
		b.ChannelId = &channelId
		b.DurationSeconds = ptr(int32(a.Metadata.DurationSeconds))
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
	code := inv.Code
	maxAge := int32(inv.MaxAge)
	maxUses := int32(inv.MaxUses)
	uses := int32(inv.Uses)

	b := &pb.Invite_builder{
		Code:      &code,
		MaxAge:    &maxAge,
		MaxUses:   &maxUses,
		Temporary: ptr(inv.Temporary),
		Uses:      &uses,
		CreatedAt: timestamppb.New(inv.CreatedAt),
	}

	if inv.Channel != nil {
		channelId := inv.Channel.ID.String()
		b.ChannelId = &channelId
	}

	if inv.Inviter != nil {
		inviterId := inv.Inviter.ID.String()
		b.InviterId = &inviterId
	}

	return b.Build()
}
