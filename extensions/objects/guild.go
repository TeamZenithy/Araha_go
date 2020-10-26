package objects

import (
	"github.com/TeamZenithy/Araha/extensions"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ExtendedGuild struct {
	*discordgo.Guild
	session *discordgo.Session
}

func ExtendGuild(guild *discordgo.Guild, session *discordgo.Session) *ExtendedGuild {
	return &ExtendedGuild{
		Guild:   guild,
		session: session,
	}
}

func (guild *ExtendedGuild) GetRole(roleID string) (*discordgo.Role, error) {
	for _, role := range guild.Roles {
		if role.ID == roleID {
			return role, nil
		}
	}
	return nil, errors.New(fmt.Sprint(extensions.RoleNotFoundError,
		"Role ", roleID, " not found in guild ", guild.ID))
}

// Better than GetRole as this finds all the Roles in 1 pass
func (guild *ExtendedGuild) GetRoles(roleIDs map[string]struct{}) (roles []*discordgo.Role, err error) {
	if len(guild.Roles) == 0 {
		var guildRolesErr error
		guild.Roles, guildRolesErr = guild.session.GuildRoles(guild.ID)
		if guildRolesErr != nil {
			return nil, guildRolesErr
		}
	}
	for _, role := range guild.Roles {
		var _, contains = roleIDs[role.ID]
		if contains {
			roles = append(roles, role)
		}
	}
	return
}

// GetRolesSlice GetRole but your slice is converted into a map
func (guild *ExtendedGuild) GetRolesSlice(roleIDs []string) (roles []*discordgo.Role, err error) {
	var roleIDMap = make(map[string]struct{})
	for _, roleID := range roleIDs {
		roleIDMap[roleID] = struct{}{}
	}
	return guild.GetRoles(roleIDMap)
}

func (guild *ExtendedGuild) GetMembers() (members []*discordgo.Member, err error) {
	members = guild.Members
	if len(members) == 0 {
		members, err = guild.session.GuildMembers(guild.ID, "0", 1000)
		if err != nil {
			return nil, err
		}
	}
	return members, nil
}