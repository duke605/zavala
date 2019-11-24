package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/duke605/zavala/app"
	"github.com/duke605/zavala/destiny2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("CONFIG_PATH", ".env")

	viper.SetConfigFile(viper.GetString("CONFIG_PATH"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Getting ngrok address
	r, err := http.Get("http://localhost:4040/api/tunnels/command_line")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var resp struct {
		URL string `json:"public_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		panic(err)
	}

	viper.Set("NGROK_URL", resp.URL)
}

func main() {

	// Setting up discord bot
	c, err := discordgo.New(fmt.Sprintf("Bot %s", viper.GetString("DISCORD_BOT_TOKEN")))
	if err != nil {
		panic(err)
	}

	// Setting up database
	repo := app.NewRepo("mysql", fmt.Sprintf("%s@tcp(%s:%d)/%s?parseTime=true",
		viper.GetString("DB_USER"),
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_DATABASE"),
	))

	// Setting up destiny 2 client
	d2Client := destiny2.NewClient(viper.GetString("BUNGIE_API_KEY"))
	d2Client.SetOAuthCredentials(
		viper.GetString("BUNGIE_CLIENT_ID"),
		viper.GetString("BUNGIE_CLIENT_SECRET"),
	)

	app := app.New(d2Client, c, repo)
	app.RunUntilInterupt()
}
