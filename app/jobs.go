package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/duke605/zavala/destiny2"
)

// scheduleEvery scuedules a function to be called everytime d elapses.
// Multiple functions can be provided and they well be called sequentially.
func (a *App) scheduleEvery(d time.Duration, fns ...func(ctx context.Context) error) {
	ticker := time.NewTicker(d)
	done := a.ctx.Done()
	a.wg.Add(1)
	defer a.wg.Done()
	defer ticker.Stop()

	for {

		// Calling the functions sequantially
		for _, fn := range fns {
			if err := fn(a.ctx); err != nil {

				// If the error was a cancelled error we break from the loop to stop all work
				if errors.Is(err, context.Canceled) {
					break
				}

				fmt.Printf("Errored running scheduled job: %s\n", err.Error())
			}
		}

		select {
		case <-done:
			return
		case <-ticker.C:
		}
	}
}

// DeleteRemovedGuilds syncs discord guild information with the database
func (a *App) DeleteRemovedGuilds(ctx context.Context) error {
	guilds := a.bot.State.Guilds

	// Syncing guild information
	sGuildIDs := make([]uint64, len(guilds))
	for i, guild := range guilds {
		guildID, _ := strconv.ParseUint(guild.ID, 10, 64)
		sGuildIDs[i] = guildID
	}

	// Deleting guilds that no longer have the bot
	return a.repo.SyncGuilds(ctx, sGuildIDs)
}

// SetNicknames sets the nickname of every member, if they have their Destiny 2
// account registered with the bot, of every guild the bot is connected to
func (a *App) SetNicknames(ctx context.Context) error {
	guilds := a.bot.State.Guilds
	wg := &sync.WaitGroup{}
	config := a.d2Client.GetOAuthConfig()

	// Sets the nickname of every member (if applicable) for the provided guild
	setNicksForGuild := func(guild *discordgo.Guild) {
		defer wg.Done()

		for _, member := range guild.Members {
			memID, _ := strconv.ParseUint(member.User.ID, 10, 64)

			// Getting the user record for the current member if they have one
			user, err := a.repo.GetUserByID(ctx, memID)
			if err != nil {

				// Member does not have a Destiny 2 account registered
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}

				fmt.Printf("Error getting user from database: %s\n", err.Error())
				continue
			}

			// Creasting token source instead of client so we can see if the token has changed
			t := user.Token()
			ts := config.TokenSource(ctx, t)
			tNew, err := ts.Token()
			if err != nil {
				fmt.Printf("Failed to refresh token for user '%d': %s\n", user.ID, err.Error())
				continue
			}

			// Getting destiny user data
			dUser, err := a.d2Client.UserService.GetMembershipDataForCurrentUser(destiny2.OptionOAuthToken(tNew))
			fmt.Println(dUser)
			if err != nil {
				fmt.Printf("Error getting user data: %s\n", err.Error())
				continue
			}

			// Finding active destiny 2 account
			for _, d2Acc := range dUser.DestinyMemberships {
				if d2Acc.CrossSaveOverride == d2Acc.MembershipType {
					if err := a.bot.GuildMemberNickname(guild.ID, member.User.ID, d2Acc.DisplayName); err != nil {
						fmt.Printf("Error changing nickname: %s\n", err.Error())
					}

					break
				}
			}
		}
	}

	// Starting a go routine for every guild
	for _, guild := range guilds {
		wg.Add(1)
		go setNicksForGuild(guild)
	}

	wg.Wait()
	return nil
}
