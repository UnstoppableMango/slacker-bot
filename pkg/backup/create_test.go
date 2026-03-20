package backup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"go.uber.org/mock/gomock"

	pb "github.com/unstoppablemango/slacker-bot/gen/dev/unmango/discord/backup/v1alpha1"
	"github.com/unstoppablemango/slacker-bot/pkg/mocks"
)

var (
	testGuildID = snowflake.ID(123456789)
	testUserID  = snowflake.ID(987654321)
	testRoleID  = snowflake.ID(111111111)
	testChanID  = snowflake.ID(222222222)

	errTest = errors.New("test error")
)

func testRestGuild() *discord.RestGuild {
	return &discord.RestGuild{
		Guild: discord.Guild{
			ID:      testGuildID,
			Name:    "Test Guild",
			OwnerID: testUserID,
		},
	}
}

func testMember() discord.Member {
	return discord.Member{
		User: discord.User{
			ID:            testUserID,
			Username:      "testuser",
			Discriminator: "0",
		},
	}
}

func setupCreate(t *testing.T) (*mocks.MockRest, snowflake.ID) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mocks.NewMockRest(ctrl), testGuildID
}

func expectHappyPath(m *mocks.MockRest, guildID snowflake.ID) {
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, nil)
}

func TestCreate_HappyPath(t *testing.T) {
	m, guildID := setupCreate(t)
	expectHappyPath(m, guildID)

	backup, err := Create(context.Background(), m, guildID, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backup == nil {
		t.Fatal("expected non-nil backup")
	}
	if got := backup.GetGuild().GetName(); got != "Test Guild" {
		t.Errorf("guild name = %q, want %q", got, "Test Guild")
	}
}

func TestCreate_GetGuildError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetGuildChannelsError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetMembersError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetAllWebhooksError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetGuildScheduledEventsError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetAutoModerationRulesError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_GetGuildInvitesError(t *testing.T) {
	m, guildID := setupCreate(t)
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, errTest)

	_, err := Create(context.Background(), m, guildID, nil)

	if !errors.Is(err, errTest) {
		t.Errorf("err = %v, want %v", err, errTest)
	}
}

func TestCreate_MapsMembers(t *testing.T) {
	m, guildID := setupCreate(t)
	members := []discord.Member{testMember()}
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(members, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, nil)

	backup, err := Create(context.Background(), m, guildID, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(backup.GetMembers()); got != 1 {
		t.Errorf("members count = %d, want 1", got)
	}
	if got := len(backup.GetUsers()); got != 1 {
		t.Errorf("users count = %d, want 1", got)
	}
	if got := backup.GetUsers()[0].GetUsername(); got != "testuser" {
		t.Errorf("username = %q, want %q", got, "testuser")
	}
}

func TestMapUsers_Deduplication(t *testing.T) {
	user := discord.User{ID: testUserID, Username: "testuser", Discriminator: "0"}
	members := []discord.Member{
		{User: user},
		{User: user},
		{User: user},
	}

	users := mapUsers(members)

	if got := len(users); got != 1 {
		t.Errorf("len(users) = %d, want 1", got)
	}
	if got := users[0].GetUsername(); got != "testuser" {
		t.Errorf("username = %q, want %q", got, "testuser")
	}
}

func TestMapUsers_MultipleDistinct(t *testing.T) {
	members := []discord.Member{
		{User: discord.User{ID: snowflake.ID(1), Username: "alice", Discriminator: "0"}},
		{User: discord.User{ID: snowflake.ID(2), Username: "bob", Discriminator: "0"}},
	}

	users := mapUsers(members)

	if got := len(users); got != 2 {
		t.Errorf("len(users) = %d, want 2", got)
	}
}

func TestMapChannelType_OffsetApplied(t *testing.T) {
	tests := []struct {
		input discord.ChannelType
		want  pb.ChannelType
	}{
		{0, pb.ChannelType(1)},
		{1, pb.ChannelType(2)},
		{5, pb.ChannelType(6)},
	}

	for _, tt := range tests {
		got := mapChannelType(tt.input)
		if got != tt.want {
			t.Errorf("mapChannelType(%d) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMapChannelType_NoOffsetAbove5(t *testing.T) {
	// Types > 5 should not have the offset applied
	tests := []struct {
		input discord.ChannelType
		want  pb.ChannelType
	}{
		{6, pb.ChannelType(6)},
		{10, pb.ChannelType(10)},
		{13, pb.ChannelType(13)},
	}

	for _, tt := range tests {
		got := mapChannelType(tt.input)
		if got != tt.want {
			t.Errorf("mapChannelType(%d) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMapAutoModTriggerType(t *testing.T) {
	tests := []struct {
		input discord.AutoModerationTriggerType
		want  pb.AutoModTriggerType
	}{
		{discord.AutoModerationTriggerTypeKeyword, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD},
		{discord.AutoModerationTriggerTypeSpam, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_SPAM},
		{discord.AutoModerationTriggerTypeKeywordPresent, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD_PRESET},
		{discord.AutoModerationTriggerTypeMentionSpam, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MENTION_SPAM},
		{discord.AutoModerationTriggerTypeMemberProfile, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MEMBER_PROFILE},
		{discord.AutoModerationTriggerType(999), pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_UNSPECIFIED},
	}

	for _, tt := range tests {
		got := mapAutoModTriggerType(tt.input)
		if got != tt.want {
			t.Errorf("mapAutoModTriggerType(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMapScheduledEvent_WithEndTime(t *testing.T) {
	endTime := time.Now().Add(time.Hour)
	event := discord.GuildScheduledEvent{
		ID:                 snowflake.ID(1),
		GuildID:            testGuildID,
		CreatorID:          testUserID,
		Name:               "Test Event",
		ScheduledStartTime: time.Now(),
		ScheduledEndTime:   &endTime,
	}

	result := mapScheduledEvent(event)

	if result.GetScheduledEndTime() == nil {
		t.Error("expected ScheduledEndTime to be set")
	}
}

func TestMapScheduledEvent_WithEntityMetadata(t *testing.T) {
	location := "Test Location"
	event := discord.GuildScheduledEvent{
		ID:                 snowflake.ID(1),
		GuildID:            testGuildID,
		CreatorID:          testUserID,
		Name:               "Test Event",
		ScheduledStartTime: time.Now(),
		EntityMetaData:     &discord.EntityMetaData{Location: location},
	}

	result := mapScheduledEvent(event)

	if result.GetEntityMetadata().GetLocation() != location {
		t.Errorf("location = %q, want %q", result.GetEntityMetadata().GetLocation(), location)
	}
}

func TestMapGuild_EnumOffset(t *testing.T) {
	guild := discord.Guild{
		ID:                          testGuildID,
		Name:                        "Test",
		OwnerID:                     testUserID,
		VerificationLevel:           discord.VerificationLevelLow, // 1
		ExplicitContentFilter:       discord.ExplicitContentFilterLevelMembersWithoutRoles, // 1
		DefaultMessageNotifications: discord.MessageNotificationsLevelAllMessages,          // 0
		NSFWLevel:                   discord.NSFWLevelDefault,                              // 0
		PremiumTier:                 discord.PremiumTierNone,                               // 0
	}

	result := mapGuild(guild)

	// Each enum value should be shifted +1
	if got := result.GetVerificationLevel(); got != pb.VerificationLevel(int32(discord.VerificationLevelLow)+1) {
		t.Errorf("VerificationLevel = %v, want %v", got, pb.VerificationLevel(int32(discord.VerificationLevelLow)+1))
	}
}

func TestMapRole_WithTags(t *testing.T) {
	botID := snowflake.ID(42)
	role := discord.Role{
		ID:   testRoleID,
		Name: "Test Role",
		Tags: &discord.RoleTag{BotID: &botID},
	}

	result := mapRole(role)

	if result.GetTags() == nil {
		t.Error("expected Tags to be set")
	}
	if got := result.GetTags().GetBotId(); got != botID.String() {
		t.Errorf("BotId = %q, want %q", got, botID.String())
	}
}

func TestMapMember_WithOptionalTimes(t *testing.T) {
	now := time.Now()
	member := discord.Member{
		User:         discord.User{ID: testUserID, Username: "u", Discriminator: "0"},
		JoinedAt:     &now,
		PremiumSince: &now,
		CommunicationDisabledUntil: &now,
	}

	result := mapMember(member)

	if result.GetJoinedAt() == nil {
		t.Error("expected JoinedAt to be set")
	}
	if result.GetPremiumSince() == nil {
		t.Error("expected PremiumSince to be set")
	}
	if result.GetCommunicationDisabledUntil() == nil {
		t.Error("expected CommunicationDisabledUntil to be set")
	}
}

func TestMapInvite_WithChannelAndInviter(t *testing.T) {
	channel := &discord.InviteChannel{ID: testChanID}
	inviter := &discord.User{ID: testUserID}
	inv := discord.ExtendedInvite{
		Invite: discord.Invite{
			Code:    "abc123",
			Channel: channel,
			Inviter: inviter,
		},
		CreatedAt: time.Now(),
	}

	result := mapInvite(inv)

	if got := result.GetCode(); got != "abc123" {
		t.Errorf("Code = %q, want %q", got, "abc123")
	}
	if got := result.GetChannelId(); got != testChanID.String() {
		t.Errorf("ChannelId = %q, want %q", got, testChanID.String())
	}
	if got := result.GetInviterId(); got != testUserID.String() {
		t.Errorf("InviterId = %q, want %q", got, testUserID.String())
	}
}
