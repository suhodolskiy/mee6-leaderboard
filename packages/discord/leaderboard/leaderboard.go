package main

import (
	"context"
	"fmt"
	"leaderboard/mee6"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

var client = api.NewClient(fmt.Sprintf("Bot %s", os.Getenv("BOT_SECRET_KEY")))

func GetGuildMembers(guildID discord.GuildID) ([]discord.Member, error) {
	members, err := client.Members(guildID, 100)

	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}

	if len(members) == 0 {
		return members, nil
	}

	for {
		lm := members[len(members)-1]

		nm, err := client.MembersAfter(guildID, lm.User.ID, 100)

		if err != nil {
			return nil, fmt.Errorf("get members after (%s): %w", lm.User.ID.String(), err)
		}

		if len(nm) == 0 {
			break
		}

		members = append(members, nm...)
	}

	return members, nil
}

func GetGuildID() (discord.GuildID, error) {
	id, err := strconv.ParseUint(os.Getenv("GUILD_ID"), 10, 64)

	if err != nil {
		return discord.NullGuildID, fmt.Errorf("parse uint: %w", err)
	}

	return discord.GuildID(id), nil
}

func GetRoleID() (discord.RoleID, error) {
	id, err := strconv.ParseUint(os.Getenv("ROLE_ID"), 10, 64)

	if err != nil {
		return discord.NullRoleID, fmt.Errorf("parse uint: %w", err)
	}

	return discord.RoleID(id), nil
}

func Main(ctx context.Context) error {
	guildID, err := GetGuildID()

	if err != nil {
		return err
	}

	roleID, err := GetRoleID()

	if err != nil {
		return err
	}

	leaderboard, err := mee6.GetLeaderboard(ctx, guildID.String())

	if err != nil {
		return fmt.Errorf("get leader board: %w", err)
	}

	members, err := GetGuildMembers(guildID)

	if err != nil {
		return fmt.Errorf("get guild members: %w", err)
	}

	if len(leaderboard.Players) > 5 {
		leaderboard.Players = leaderboard.Players[:5]
	}

	var wg sync.WaitGroup

	for _, member := range members {
		var isTopPlayer bool

		for _, player := range leaderboard.Players {
			if player.ID == member.User.ID.String() {
				isTopPlayer = true

				break
			}
		}

		var hasRole bool

		for _, rid := range member.RoleIDs {
			if rid == roleID {
				hasRole = true

				break
			}
		}

		if !isTopPlayer && !hasRole || isTopPlayer && hasRole {
			continue
		}

		var roles []discord.RoleID

		if !isTopPlayer && hasRole {
			for _, rid := range member.RoleIDs {
				if rid == roleID {
					continue
				}

				roles = append(roles, roleID)
			}
		} else {
			roles = append(member.RoleIDs, roleID)
		}

		member.RoleIDs = roles

		wg.Add(1)

		go func(member discord.Member) {
			defer wg.Done()

			err := client.ModifyMember(guildID, member.User.ID, api.ModifyMemberData{
				Roles: &member.RoleIDs,
			})

			if err != nil {
				log.Printf("failed modify user (%s) roles: %s", member.User.ID.String(), err)
			}
		}(member)
	}

	wg.Wait()

	return nil
}
