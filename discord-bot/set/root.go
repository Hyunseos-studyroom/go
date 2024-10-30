package set

import (
	"discord-bot/handler"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
)

func Setup(db *mongo.Client) {
	message := handler.Message{DB: db}
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(message.MessageInfoMsg)
	dg.AddHandler(message.SendingEmbed)
	dg.AddHandler(message.ThisMonthSalary)

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.UpdateListeningStatus("명령어 리스트는 $사용법")
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
