package objects

import (
	"github.com/TeamZenithy/Araha/extensions/permissions"
	"github.com/bwmarrin/discordgo"
)

type ExtendedMember struct {
	*discordgo.Member
	session *discordgo.Session
}

func ExtendMember(member *discordgo.Member, session *discordgo.Session) *ExtendedMember {
	return &ExtendedMember{
		Member:  member,
		session: session,
	}
}

func (member *ExtendedMember) Guild() (*ExtendedGuild, error) {
	guild, err := member.session.Guild(member.GuildID)
	if err != nil {
		return nil, err
	}
	return ExtendGuild(guild, member.session), nil
}

func (member *ExtendedMember) HasAllPermissions(requestedPermissions ...int) (bool, error) {
	memberGuild, err := member.Guild()
	if err != nil {
		return false, err
	}

	if memberGuild.OwnerID == member.User.ID {
		return true, nil
	}

	if len(member.Roles) == 0 {
		var reGetMember, reGetMemberErr = member.session.GuildMember(member.GuildID, member.User.ID)
		if reGetMemberErr != nil {
			return false, reGetMemberErr
		}
		member.Roles = reGetMember.Roles
	}

	var roles, rolesErr = memberGuild.GetRolesSlice(member.Roles)
	if rolesErr != nil {
		return false, rolesErr
	}

	combinedPermissionInteger := int64(0)

	for _, role := range roles {
		combinedPermissionInteger |= role.Permissions
	}

	return permissions.IsPermittedAll(int(combinedPermissionInteger), requestedPermissions...), nil
}
