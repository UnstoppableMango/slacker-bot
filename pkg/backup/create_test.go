package backup

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/unstoppablemango/slacker-bot/gen/mocks"
	pb "github.com/unstoppablemango/slacker-bot/gen/pb/dev/unmango/discord/backup/v1alpha1"
	"go.uber.org/mock/gomock"
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

func expectHappyPath(m *mocks.MockRest, guildID snowflake.ID) {
	m.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
	m.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
	m.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, nil)
}

var _ = Describe("Create", func() {
	var (
		ctrl    *gomock.Controller
		rest    *mocks.MockRest
		guildID snowflake.ID
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rest = mocks.NewMockRest(ctrl)
		guildID = testGuildID
	})

	Context("when all API calls succeed", func() {
		BeforeEach(func() {
			expectHappyPath(rest, guildID)
		})

		It("returns a non-nil backup", func() {
			backup, err := Create(context.Background(), rest, guildID)
			Expect(err).NotTo(HaveOccurred())
			Expect(backup).NotTo(BeNil())
		})

		It("maps the guild name", func() {
			backup, err := Create(context.Background(), rest, guildID)
			Expect(err).NotTo(HaveOccurred())
			Expect(backup.GetGuild().GetName()).To(Equal("Test Guild"))
		})
	})

	Context("when members are returned", func() {
		BeforeEach(func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return([]discord.Member{testMember()}, nil)
			rest.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, nil)
		})

		It("maps members", func() {
			backup, err := Create(context.Background(), rest, guildID)
			Expect(err).NotTo(HaveOccurred())
			Expect(backup.GetMembers()).To(HaveLen(1))
		})

		It("maps users from members", func() {
			backup, err := Create(context.Background(), rest, guildID)
			Expect(err).NotTo(HaveOccurred())
			Expect(backup.GetUsers()).To(HaveLen(1))
			Expect(backup.GetUsers()[0].GetUsername()).To(Equal("testuser"))
		})
	})

	DescribeTable("propagates API errors",
		func(setup func()) {
			setup()
			_, err := Create(context.Background(), rest, guildID)
			Expect(errors.Is(err, errTest)).To(BeTrue())
		},
		Entry("GetGuild", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetGuildChannels", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetMembers", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetAllWebhooks", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetGuildScheduledEvents", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetAutoModerationRules", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, errTest)
		}),
		Entry("GetGuildInvites", func() {
			rest.EXPECT().GetGuild(guildID, false, gomock.Any()).Return(testRestGuild(), nil)
			rest.EXPECT().GetGuildChannels(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetMembers(guildID, 1000, snowflake.ID(0), gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAllWebhooks(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildScheduledEvents(guildID, false, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetAutoModerationRules(guildID, gomock.Any()).Return(nil, nil)
			rest.EXPECT().GetGuildInvites(guildID, gomock.Any()).Return(nil, errTest)
		}),
	)
})

var _ = Describe("mapUsers", func() {
	It("deduplicates members with the same user ID", func() {
		user := discord.User{ID: testUserID, Username: "testuser", Discriminator: "0"}
		members := []discord.Member{{User: user}, {User: user}, {User: user}}

		users := mapUsers(members)

		Expect(users).To(HaveLen(1))
		Expect(users[0].GetUsername()).To(Equal("testuser"))
	})

	It("includes multiple distinct users", func() {
		members := []discord.Member{
			{User: discord.User{ID: snowflake.ID(1), Username: "alice", Discriminator: "0"}},
			{User: discord.User{ID: snowflake.ID(2), Username: "bob", Discriminator: "0"}},
		}

		Expect(mapUsers(members)).To(HaveLen(2))
	})
})

var _ = Describe("mapChannelType", func() {
	DescribeTable("applies +1 offset for types 0-5",
		func(input discord.ChannelType, want pb.ChannelType) {
			Expect(mapChannelType(input)).To(Equal(want))
		},
		Entry("type 0", discord.ChannelType(0), pb.ChannelType(1)),
		Entry("type 1", discord.ChannelType(1), pb.ChannelType(2)),
		Entry("type 5", discord.ChannelType(5), pb.ChannelType(6)),
	)

	DescribeTable("does not apply offset for types above 5",
		func(input discord.ChannelType, want pb.ChannelType) {
			Expect(mapChannelType(input)).To(Equal(want))
		},
		Entry("type 6", discord.ChannelType(6), pb.ChannelType(6)),
		Entry("type 10", discord.ChannelType(10), pb.ChannelType(10)),
		Entry("type 13", discord.ChannelType(13), pb.ChannelType(13)),
	)
})

var _ = Describe("mapAutoModTriggerType", func() {
	DescribeTable("maps all trigger types",
		func(input discord.AutoModerationTriggerType, want pb.AutoModTriggerType) {
			Expect(mapAutoModTriggerType(input)).To(Equal(want))
		},
		Entry("keyword", discord.AutoModerationTriggerTypeKeyword, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD),
		Entry("spam", discord.AutoModerationTriggerTypeSpam, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_SPAM),
		Entry("keyword present", discord.AutoModerationTriggerTypeKeywordPresent, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_KEYWORD_PRESET),
		Entry("mention spam", discord.AutoModerationTriggerTypeMentionSpam, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MENTION_SPAM),
		Entry("member profile", discord.AutoModerationTriggerTypeMemberProfile, pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_MEMBER_PROFILE),
		Entry("unknown", discord.AutoModerationTriggerType(999), pb.AutoModTriggerType_AUTO_MOD_TRIGGER_TYPE_UNSPECIFIED),
	)
})

var _ = Describe("mapScheduledEvent", func() {
	It("maps optional ScheduledEndTime when present", func() {
		endTime := time.Now().Add(time.Hour)
		event := discord.GuildScheduledEvent{
			ID:                 snowflake.ID(1),
			GuildID:            testGuildID,
			CreatorID:          testUserID,
			Name:               "Test Event",
			ScheduledStartTime: time.Now(),
			ScheduledEndTime:   &endTime,
		}

		Expect(mapScheduledEvent(event).GetScheduledEndTime()).NotTo(BeNil())
	})

	It("maps entity metadata location when present", func() {
		location := "Test Location"
		event := discord.GuildScheduledEvent{
			ID:                 snowflake.ID(1),
			GuildID:            testGuildID,
			CreatorID:          testUserID,
			Name:               "Test Event",
			ScheduledStartTime: time.Now(),
			EntityMetaData:     &discord.EntityMetaData{Location: location},
		}

		Expect(mapScheduledEvent(event).GetEntityMetadata().GetLocation()).To(Equal(location))
	})
})

var _ = Describe("mapGuild", func() {
	It("shifts enum values by +1", func() {
		guild := discord.Guild{
			ID:                          testGuildID,
			Name:                        "Test",
			OwnerID:                     testUserID,
			VerificationLevel:           discord.VerificationLevelLow,
			ExplicitContentFilter:       discord.ExplicitContentFilterLevelMembersWithoutRoles,
			DefaultMessageNotifications: discord.MessageNotificationsLevelAllMessages,
			NSFWLevel:                   discord.NSFWLevelDefault,
			PremiumTier:                 discord.PremiumTierNone,
		}

		result := mapGuild(guild)

		Expect(result.GetVerificationLevel()).To(Equal(pb.VerificationLevel(int32(discord.VerificationLevelLow) + 1)))
	})
})

var _ = Describe("mapRole", func() {
	It("maps optional tags when present", func() {
		botID := snowflake.ID(42)
		role := discord.Role{
			ID:   testRoleID,
			Name: "Test Role",
			Tags: &discord.RoleTag{BotID: &botID},
		}

		result := mapRole(role)

		Expect(result.GetTags()).NotTo(BeNil())
		Expect(result.GetTags().GetBotId()).To(Equal(botID.String()))
	})
})

var _ = Describe("mapMember", func() {
	It("maps optional time fields when present", func() {
		now := time.Now()
		member := discord.Member{
			User:                       discord.User{ID: testUserID, Username: "u", Discriminator: "0"},
			JoinedAt:                   &now,
			PremiumSince:               &now,
			CommunicationDisabledUntil: &now,
		}

		result := mapMember(member)

		Expect(result.GetJoinedAt()).NotTo(BeNil())
		Expect(result.GetPremiumSince()).NotTo(BeNil())
		Expect(result.GetCommunicationDisabledUntil()).NotTo(BeNil())
	})
})

var _ = Describe("mapInvite", func() {
	It("maps optional channel and inviter when present", func() {
		inv := discord.ExtendedInvite{
			Invite: discord.Invite{
				Code:    "abc123",
				Channel: &discord.InviteChannel{ID: testChanID},
				Inviter: &discord.User{ID: testUserID},
			},
			CreatedAt: time.Now(),
		}

		result := mapInvite(inv)

		Expect(result.GetCode()).To(Equal("abc123"))
		Expect(result.GetChannelId()).To(Equal(testChanID.String()))
		Expect(result.GetInviterId()).To(Equal(testUserID.String()))
	})
})
