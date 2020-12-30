package permissions

// https://discord.com/developers/docs/topics/permissions
//goland:noinspection GoSnakeCaseUsage,GoUnnecessarilyExportedIdentifiers
const (
	NONE                  = 0x00000000 // No Permission Required
	CREATE_INSTANT_INVITE = 0x00000001 // Allows creation of instant invites
	KICK_MEMBERS          = 0x00000002 // Allows kicking members
	BAN_MEMBERS           = 0x00000004 // Allows banning members
	ADMINISTRATOR         = 0x00000008 // Allows all permissions and bypasses channel permission overwrites
	MANAGE_CHANNELS       = 0x00000010 // Allows management and editing of channels
	MANAGE_GUILD          = 0x00000020 // Allows management and editing of the guild
	ADD_REACTIONS         = 0x00000040 // Allows for the addition of reactions to messages
	VIEW_AUDIT_LOG        = 0x00000080 // Allows for viewing of audit logs
	PRIORITY_SPEAKER      = 0x00000100 // Allows for using priority speaker in a voice channel
	STREAM                = 0x00000200 // Allows the user to go live
	VIEW_CHANNEL          = 0x00000400 // Allows guild members to view a channel, which includes reading messages in text channels
	SEND_MESSAGES         = 0x00000800 // Allows for sending messages in a channel
	SEND_TTS_MESSAGES     = 0x00001000 // Allows for sending of /tts messages
	MANAGE_MESSAGES       = 0x00002000 // Allows for deletion of other users messages
	EMBED_LINKS           = 0x00004000 // Links sent by users with this permission will be auto-embedded
	ATTACH_FILES          = 0x00008000 // Allows for uploading images and files
	READ_MESSAGE_HISTORY  = 0x00010000 // Allows for reading of message history
	MENTION_EVERYONE      = 0x00020000 // Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel
	USE_EXTERNAL_EMOJIS   = 0x00040000 // Allows the usage of custom emojis from other servers
	VIEW_GUILD_INSIGHTS   = 0x00080000 // Allows for viewing guild insights
	CONNECT               = 0x00100000 // Allows for joining of a voice channel
	SPEAK                 = 0x00200000 // Allows for speaking in a voice channel
	MUTE_MEMBERS          = 0x00400000 // Allows for muting members in a voice channel
	DEAFEN_MEMBERS        = 0x00800000 // Allows for deafening of members in a voice channel
	MOVE_MEMBERS          = 0x01000000 // Allows for moving of members between voice channels
	USE_VAD               = 0x02000000 // Allows for using voice-activity-detection in a voice channel
	CHANGE_NICKNAME       = 0x04000000 // Allows for modification of own nickname
	MANAGE_NICKNAMES      = 0x08000000 // Allows for modification of other users nicknames
	MANAGE_ROLES          = 0x10000000 // Allows management and editing of roles
	MANAGE_WEBHOOKS       = 0x20000000 // Allows management and editing of webhooks
	MANAGE_EMOJIS         = 0x40000000 // Allows management and editing of emojis
)

var (
	NameWithValue = map[string]int{
		"NONE":                  NONE,
		"CREATE_INSTANT_INVITE": CREATE_INSTANT_INVITE,
		"KICK_MEMBERS":          KICK_MEMBERS,
		"BAN_MEMBERS":           BAN_MEMBERS,
		"ADMINISTRATOR":         ADMINISTRATOR,
		"MANAGE_CHANNELS":       MANAGE_CHANNELS,
		"MANAGE_GUILD":          MANAGE_GUILD,
		"ADD_REACTIONS":         ADD_REACTIONS,
		"VIEW_AUDIT_LOG":        VIEW_AUDIT_LOG,
		"PRIORITY_SPEAKER":      PRIORITY_SPEAKER,
		"STREAM":                STREAM,
		"VIEW_CHANNEL":          VIEW_CHANNEL,
		"SEND_MESSAGES":         SEND_MESSAGES,
		"SEND_TTS_MESSAGES":     SEND_TTS_MESSAGES,
		"MANAGE_MESSAGES":       MANAGE_MESSAGES,
		"EMBED_LINKS":           EMBED_LINKS,
		"ATTACH_FILES":          ATTACH_FILES,
		"READ_MESSAGE_HISTORY":  READ_MESSAGE_HISTORY,
		"MENTION_EVERYONE":      MENTION_EVERYONE,
		"USE_EXTERNAL_EMOJIS":   USE_EXTERNAL_EMOJIS,
		"VIEW_GUILD_INSIGHTS":   VIEW_GUILD_INSIGHTS,
		"CONNECT":               CONNECT,
		"SPEAK":                 SPEAK,
		"MUTE_MEMBERS":          MUTE_MEMBERS,
		"DEAFEN_MEMBERS":        DEAFEN_MEMBERS,
		"MOVE_MEMBERS":          MOVE_MEMBERS,
		"USE_VAD":               USE_VAD,
		"CHANGE_NICKNAME":       CHANGE_NICKNAME,
		"MANAGE_NICKNAMES":      MANAGE_NICKNAMES,
		"MANAGE_ROLES":          MANAGE_ROLES,
		"MANAGE_WEBHOOKS":       MANAGE_WEBHOOKS,
		"MANAGE_EMOJIS":         MANAGE_EMOJIS,
	}

	ValueWithName = map[int]string{
		NONE:                  "NONE",
		CREATE_INSTANT_INVITE: "CREATE_INSTANT_INVITE",
		KICK_MEMBERS:          "KICK_MEMBERS",
		BAN_MEMBERS:           "BAN_MEMBERS",
		ADMINISTRATOR:         "ADMINISTRATOR",
		MANAGE_CHANNELS:       "MANAGE_CHANNELS",
		MANAGE_GUILD:          "MANAGE_GUILD",
		ADD_REACTIONS:         "ADD_REACTIONS",
		VIEW_AUDIT_LOG:        "VIEW_AUDIT_LOG",
		PRIORITY_SPEAKER:      "PRIORITY_SPEAKER",
		STREAM:                "STREAM",
		VIEW_CHANNEL:          "VIEW_CHANNEL",
		SEND_MESSAGES:         "SEND_MESSAGES",
		SEND_TTS_MESSAGES:     "SEND_TTS_MESSAGES",
		MANAGE_MESSAGES:       "MANAGE_MESSAGES",
		EMBED_LINKS:           "EMBED_LINKS",
		ATTACH_FILES:          "ATTACH_FILES",
		READ_MESSAGE_HISTORY:  "READ_MESSAGE_HISTORY",
		MENTION_EVERYONE:      "MENTION_EVERYONE",
		USE_EXTERNAL_EMOJIS:   "USE_EXTERNAL_EMOJIS",
		VIEW_GUILD_INSIGHTS:   "VIEW_GUILD_INSIGHTS",
		CONNECT:               "CONNECT",
		SPEAK:                 "SPEAK",
		MUTE_MEMBERS:          "MUTE_MEMBERS",
		DEAFEN_MEMBERS:        "DEAFEN_MEMBERS",
		MOVE_MEMBERS:          "MOVE_MEMBERS",
		USE_VAD:               "USE_VAD",
		CHANGE_NICKNAME:       "CHANGE_NICKNAME",
		MANAGE_NICKNAMES:      "MANAGE_NICKNAMES",
		MANAGE_ROLES:          "MANAGE_ROLES",
		MANAGE_WEBHOOKS:       "MANAGE_WEBHOOKS",
		MANAGE_EMOJIS:         "MANAGE_EMOJIS",
	}
)
