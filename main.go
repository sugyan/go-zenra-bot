package main

import (
	"github.com/sugyan/go-zenra"
	"github.com/sugyan/mentionbot"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"
	"unicode/utf8"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	config := &mentionbot.Config{
		UserID:            os.Getenv("USER_ID"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}
	bot := mentionbot.NewBot(config)
	bot.SetMentioner(&Zenra{
		zenrizer: zenra.NewZenrizer(),
	})
	bot.Debug(true)
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

// Zenra type implements mentionbot.Mentioner
type Zenra struct {
	zenrizer *zenra.Zenrizer
}

// Mention returns mention
func (z *Zenra) Mention(tweet *mentionbot.Tweet) (mention *string) {
	if tweet.InReplyToStatusID > 0 {
		return
	}
	if tweet.InReplyToUserID > 0 {
		return
	}
	if tweet.RetweetedStatus != nil {
		return
	}
	if tweet.User.Protected {
		return
	}

	// フォロワーが多いので適当に間引く
	if rand.Intn(100) < 50 {
		return
	}

	zenrized := z.zenrizer.Zenrize(tweet.Text)
	if tweet.Text == zenrized {
		return nil
	}
	log.Println(tweet.Text)
	screenNameLen := len(tweet.User.ScreenName)
	text := "が全裸で言った: " + zenrized
	// 140字以内に収まるように切る
	if utf8.RuneCountInString(text)+screenNameLen+2 > 140 {
		text = string([]rune(text)[0:140-len(tweet.User.ScreenName)-3]) + "…"
	}
	if !regexp.MustCompile(zenra.ZENRA).MatchString(text) {
		return nil
	}
	return &text
}
