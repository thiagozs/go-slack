package slackr

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/slack-go/slack"
	"github.com/thiagozs/go-slack/options"
	"github.com/thiagozs/go-slack/pkg/fuzzy"
)

func NewSlackClient(opts []options.Options) (*Slackr, error) {
	c := &options.OptionsParams{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	api := slack.New(
		c.Token,
		slack.OptionDebug(c.Debug),
		slack.OptionLog(log.New(os.Stdout, "log: ", log.Lshortfile|log.LstdFlags)),
	)

	return &Slackr{
		client: api,
		cfg:    c,
	}, nil

}

func (s *Slackr) SearchFuzzyLoadTerms(terms []string) {
	s.terms = terms
}

func (s *Slackr) GetUsers() ([]SlackrUser, error) {
	uopts := []slack.GetUsersOption{
		slack.GetUsersOptionPresence(true),
		slack.GetUsersOptionLimit(1000),
	}

	users, err := s.client.GetUsers(uopts...)
	if err != nil {
		return nil, err
	}

	sus := []SlackrUser{}
	for _, user := range users {
		bts, err := json.Marshal(user)
		if err != nil {
			fmt.Println("json Marshal, ", err)
			return nil, err
		}
		su := SlackrUser{}
		if err := json.Unmarshal(bts, &su); err != nil {
			fmt.Println("json Unmarshal, ", err)

			return nil, err
		}
		sus = append(sus, su)
	}

	return sus, nil
}

func (s *Slackr) GetUserGroups() ([]slack.UserGroup, error) {
	opts := []slack.GetUserGroupsOption{
		slack.GetUserGroupsOptionIncludeUsers(true),
	}

	groups, err := s.client.GetUserGroups(opts...)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *Slackr) GetClient() *slack.Client {
	return s.client
}

func (s *Slackr) GetConfig() *options.OptionsParams {
	return s.cfg
}

func (s *Slackr) GetToken() string {
	return s.cfg.Token
}

func (s *Slackr) getUserByEmail(email string) (SlackrUser, error) {

	var usersSlack []SlackrUser

	if s.cached {
		users, err := s.GetUsersFromCached()
		if err != nil {
			return SlackrUser{}, err
		}
		usersSlack = users
	} else {
		users, err := s.GetUsers()
		if err != nil {
			return SlackrUser{}, err
		}
		usersSlack = users
	}

	for _, user := range usersSlack {
		if user.Profile.Email == email {
			return user, nil
		}
	}

	return SlackrUser{}, nil
}

func (s *Slackr) SendPrivateMessage(user string, text string) error {

	// response.Channel, response.NoOp, response.AlreadyOpen, response.Err()
	ch, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user},
	})
	if err != nil {
		return err
	}

	opts := slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text},
		},
	)

	// response.Channel, response.Timestamp, err
	_, _, err = s.client.PostMessage(ch.ID, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Slackr) SearchFuzzyMatch(kind Kind, term string) ([]ResultFuzzy, error) {
	s.term = term
	s.terms = []string{}

	var usersSlack []SlackrUser

	if s.cached {
		users, err := s.GetUsersFromCached()
		if err != nil {
			return []ResultFuzzy{}, err
		}
		usersSlack = users
	} else {
		users, err := s.GetUsers()
		if err != nil {
			return []ResultFuzzy{}, err
		}
		usersSlack = users
	}

	for _, user := range usersSlack {
		switch kind {
		case EMAIL:
			s.terms = append(s.terms, user.Profile.Email)
		case REALNAME:
			s.terms = append(s.terms, user.Profile.RealName)
		case FIRSTNAME:
			s.terms = append(s.terms, user.Profile.FirstName)
		case LASTNAME:
			s.terms = append(s.terms, user.Profile.LastName)
		}
	}

	s.fuzzy = fuzzy.NewFzfSearcher(s.terms)

	rf := s.fuzzy.Search(s.term)
	rr := []ResultFuzzy{}

	for _, r := range rf {
		for _, u := range usersSlack {

			user := SlackrUser{}

			switch kind {
			case EMAIL:
				if u.Profile.Email == r.SortKey {
					user = u
				}
			case REALNAME:
				if u.Profile.RealName == r.SortKey {
					user = u
				}
			case FIRSTNAME:
				if u.Profile.FirstName == r.SortKey {
					user = u
				}
			case LASTNAME:
				if u.Profile.LastName == r.SortKey {
					user = u
				}
			}

			if !reflect.DeepEqual(user, SlackrUser{}) {
				rr = append(rr, ResultFuzzy{
					Match:   r.Match,
					Query:   r.Query,
					Score:   r.Score,
					SortKey: r.SortKey,
					User:    user,
				})
			}
		}
	}

	return rr, nil
}

func (s *Slackr) SearchByEmail(email string) (SlackrUser, error) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return SlackrUser{}, err
	}

	bts, err := json.Marshal(user)
	if err != nil {
		return SlackrUser{}, err
	}

	u := SlackrUser{}
	if err := json.Unmarshal(bts, &u); err != nil {
		return SlackrUser{}, err
	}

	return u, nil
}

func (s *Slackr) SendMessageChannel(channel string, text string) error {
	opts := slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text},
		},
	)

	// response.Channel, response.Timestamp, err
	_, _, err := s.client.PostMessage(channel, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Slackr) SetCached(cached bool) {
	s.cached = cached
}

func (s *Slackr) GetCached() bool {
	return s.cached
}

func (s *Slackr) GetUsersSlack() []SlackrUser {
	return s.users
}

func (s *Slackr) SetUsersSlack(users []SlackrUser) {
	s.users = users
}

func (s *Slackr) GetUsersFromCached() ([]SlackrUser, error) {
	if s.cached && len(s.users) == 0 {

		users, err := s.GetUsers()
		if err != nil {
			return []SlackrUser{}, err
		}

		s.users = users
	}

	return s.users, nil
}
