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
	// in_reply_to_* 付きは無視
	if tweet.InReplyToStatusID > 0 {
		return
	}
	if tweet.InReplyToUserID > 0 {
		return
	}
	// RTも無視
	if tweet.RetweetedStatus != nil {
		return
	}
	// 非公開アカウントも無視
	if tweet.User.Protected {
		return
	}
	// 画像付きも無視
	if len(tweet.Entities.Media) > 0 {
		return
	}
	// mention付きも無視
	if len(tweet.Entities.UserMentions) > 0 {
		return
	}
	// ハッシュタグ付きも無視
	if len(tweet.Entities.Hashtags) > 0 {
		return
	}
	// シンボルも一応無視
	if len(tweet.Entities.Symbols) > 0 {
		return
	}
	// その他メディアも無視
	if len(tweet.Entities.ExtendedEntities) > 0 {
		return
	}

	// フォロワーが多いので適当に間引く
	if rand.Intn(100) < 60 {
		return
	}

	// 全裸にする
	zenrized := z.zenrizer.Zenrize(tweet.Text)
	if tweet.Text == zenrized {
		return nil
	}
	screenNameLen := len(tweet.User.ScreenName)
	text := "が全裸で言った: " + zenrized
	// 140字以内に収まるように切る
	if utf8.RuneCountInString(text)+screenNameLen+2 > 140 {
		text = string([]rune(text)[0:140-len(tweet.User.ScreenName)-3]) + "…"
	}

	// 最終チェック
	if !regexp.MustCompile(zenra.ZENRA).MatchString(text) {
		return nil
	}
	return &text
}
