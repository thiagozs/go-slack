package main

import (
	"flag"
	"log"
	"os"

	"github.com/thiagozs/go-slack/options"
	"github.com/thiagozs/go-slack/pkg/utils"
	"github.com/thiagozs/go-slack/slackr"
)

var term string
var email string

func main() {

	flag.StringVar(&term, "term", "", "term to search")
	flag.StringVar(&email, "email", "", "email to search")

	flag.Parse()

	if term == "" {
		log.Fatal("term is required")
	}

	token := os.Getenv("SLACKBOT_TOKEN")

	lopts := []options.Options{
		options.CfgDebug(false),
		options.CfgToken(token),
	}

	slk, err := slackr.NewSlackClient(lopts)
	if err != nil {
		log.Fatal(err)
	}

	utils.TitleBreak("Users")

	sendMsg := false
	score := 0.0
	userFound := slackr.SlackrUser{}

	byEmail := len(email) > 0
	byMatch := len(email) == 0

	log.Printf("Filter term: %s\n", term)
	log.Printf("Filter email: %s\n", email)
	log.Printf("Filter byEmail: %t\n", byEmail)
	log.Printf("Filter byMatch: %t\n", byMatch)

	// set cache, this will load all users from slack once time
	// not mandatory, but will improve performance
	slk.SetCached(true)

	switch {
	case byEmail:
		user, err := slk.SearchByEmail(email)
		if err != nil {
			log.Fatal(err)
		}
		userFound = user
		sendMsg = true

	case byMatch:
		users, err := slk.SearchFuzzyMatch(slackr.REALNAME, term)
		if err != nil {
			log.Fatal(err)
		}
		for _, res := range users {
			if res.Score > 0.0 {
				if res.Score >= score &&
					userFound.Profile.RealName != res.User.Profile.RealName {
					log.Printf("HIGH Score, AddLastPosition Match - id:%s username:%s realname:%s email:%s\n",
						res.User.ID, res.User.Profile.FirstName,
						res.User.Profile.RealName,
						res.User.Profile.Email)
					userFound = res.User
					sendMsg = true
					score = res.Score
				}
			}
		}
	}

	if sendMsg {
		log.Printf("Sending message to id:%s realname:%s\n", userFound.ID, userFound.Profile.RealName)
		if err := slk.SendPrivateMessage(userFound.ID, "Hello from Go!"); err != nil {
			log.Fatal("Error : ", err)
		}
	}

	log.Printf("Done\n")

}
