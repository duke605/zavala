package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/duke605/zavala/destiny2"
)

type App struct {
	d2Client *destiny2.Client
	bot      *discordgo.Session
	repo     *Repo
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates and returns new instance of App
func New(d2Client *destiny2.Client, bot *discordgo.Session, repo *Repo) *App {
	ctx, cancel := context.WithCancel(context.Background())
	a := &App{
		d2Client: d2Client,
		repo:     repo,
		bot:      bot,
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
	}

	// Starting go routines when bot is ready
	a.bot.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.Ready) {
		fmt.Println("Bot ready!")
		go a.scheduleEvery(time.Minute, a.DeleteRemovedGuilds, a.SetNicknames)
	})

	return a
}

// RunUntilInterupt connects to discord and will run until the program is interupted
func (a *App) RunUntilInterupt() {
	a.bot.Open()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	a.cancel()
	a.bot.Close()
	a.wg.Wait()
}

// HandleMessage handles message create events from discord
func (a *App) HandleMessage(sess *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignoring self, very important
	if m.Member.User.ID == sess.State.User.ID {
		return
	}
}

func (a *App) getGuildConfig(g *discordgo.Guild) (Guild, error) {
	guildID, _ := strconv.ParseUint(g.ID, 10, 64)

	// Getting guild from database for setting and group id
	dbGuild, err := a.repo.GetGuildByID(a.ctx, guildID)
	if err != nil {
		return dbGuild, err
	}

	return dbGuild, nil
}
