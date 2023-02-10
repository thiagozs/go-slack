package slackr

import (
	"encoding/json"

	"github.com/slack-go/slack"
	"github.com/thiagozs/go-slack/options"
	"github.com/thiagozs/go-slack/pkg/fuzzy"
)

type Kind int

const (
	EMAIL Kind = iota
	REALNAME
	FIRSTNAME
	LASTNAME
)

func (d Kind) String() string {
	return []string{"email", "realName", "firstName", "lastName"}[d]
}

type ResultFuzzy struct {
	Match   bool
	Query   string
	Score   float64
	SortKey string
	User    SlackrUser
}

type Slackr struct {
	client *slack.Client
	cfg    *options.OptionsParams
	fuzzy  *fuzzy.FzfSearcher
	terms  []string
	term   string
	cached bool
	users  []SlackrUser
}

type JSONTime int64

type SlackrUser struct {
	ID                string         `json:"id"`
	TeamID            string         `json:"team_id"`
	Name              string         `json:"name"`
	Deleted           bool           `json:"deleted"`
	Color             string         `json:"color"`
	RealName          string         `json:"real_name"`
	TZ                string         `json:"tz,omitempty"`
	TZLabel           string         `json:"tz_label"`
	TZOffset          int            `json:"tz_offset"`
	Profile           UserProfile    `json:"profile"`
	IsBot             bool           `json:"is_bot"`
	IsAdmin           bool           `json:"is_admin"`
	IsOwner           bool           `json:"is_owner"`
	IsPrimaryOwner    bool           `json:"is_primary_owner"`
	IsRestricted      bool           `json:"is_restricted"`
	IsUltraRestricted bool           `json:"is_ultra_restricted"`
	IsStranger        bool           `json:"is_stranger"`
	IsAppUser         bool           `json:"is_app_user"`
	IsInvitedUser     bool           `json:"is_invited_user"`
	Has2FA            bool           `json:"has_2fa"`
	HasFiles          bool           `json:"has_files"`
	Presence          string         `json:"presence"`
	Locale            string         `json:"locale"`
	Updated           JSONTime       `json:"updated"`
	Enterprise        EnterpriseUser `json:"enterprise_user,omitempty"`
}

type UserProfile struct {
	FirstName              string                              `json:"first_name"`
	LastName               string                              `json:"last_name"`
	RealName               string                              `json:"real_name"`
	RealNameNormalized     string                              `json:"real_name_normalized"`
	DisplayName            string                              `json:"display_name"`
	DisplayNameNormalized  string                              `json:"display_name_normalized"`
	Email                  string                              `json:"email"`
	Skype                  string                              `json:"skype"`
	Phone                  string                              `json:"phone"`
	Image24                string                              `json:"image_24"`
	Image32                string                              `json:"image_32"`
	Image48                string                              `json:"image_48"`
	Image72                string                              `json:"image_72"`
	Image192               string                              `json:"image_192"`
	Image512               string                              `json:"image_512"`
	ImageOriginal          string                              `json:"image_original"`
	Title                  string                              `json:"title"`
	BotID                  string                              `json:"bot_id,omitempty"`
	ApiAppID               string                              `json:"api_app_id,omitempty"`
	StatusText             string                              `json:"status_text,omitempty"`
	StatusEmoji            string                              `json:"status_emoji,omitempty"`
	StatusEmojiDisplayInfo []UserProfileStatusEmojiDisplayInfo `json:"status_emoji_display_info,omitempty"`
	StatusExpiration       int                                 `json:"status_expiration"`
	Team                   string                              `json:"team"`
	Fields                 UserProfileCustomFields             `json:"fields"`
}

type UserProfileStatusEmojiDisplayInfo struct {
	EmojiName    string `json:"emoji_name"`
	DisplayAlias string `json:"display_alias,omitempty"`
	DisplayURL   string `json:"display_url,omitempty"`
	Unicode      string `json:"unicode,omitempty"`
}

type UserProfileCustomFields struct {
	fields map[string]UserProfileCustomField
}

type UserProfileCustomField struct {
	Value string `json:"value"`
	Alt   string `json:"alt"`
	Label string `json:"label"`
}

type EnterpriseUser struct {
	ID             string   `json:"id"`
	EnterpriseID   string   `json:"enterprise_id"`
	EnterpriseName string   `json:"enterprise_name"`
	IsAdmin        bool     `json:"is_admin"`
	IsOwner        bool     `json:"is_owner"`
	Teams          []string `json:"teams"`
}

// UnmarshalJSON is the implementation of the json.Unmarshaler interface.
func (fields *UserProfileCustomFields) UnmarshalJSON(b []byte) error {
	// https://github.com/slack-go/slack/pull/298#discussion_r185159233
	if string(b) == "[]" {
		return nil
	}
	return json.Unmarshal(b, &fields.fields)
}

// MarshalJSON is the implementation of the json.Marshaler interface.
func (fields UserProfileCustomFields) MarshalJSON() ([]byte, error) {
	if len(fields.fields) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(fields.fields)
}

// ToMap returns a map of custom fields.
func (fields *UserProfileCustomFields) ToMap() map[string]UserProfileCustomField {
	return fields.fields
}

// Len returns the number of custom fields.
func (fields *UserProfileCustomFields) Len() int {
	return len(fields.fields)
}

// SetMap sets a map of custom fields.
func (fields *UserProfileCustomFields) SetMap(m map[string]UserProfileCustomField) {
	fields.fields = m
}

// FieldsMap returns a map of custom fields.
func (profile *UserProfile) FieldsMap() map[string]UserProfileCustomField {
	return profile.Fields.ToMap()
}

// SetFieldsMap sets a map of custom fields.
func (profile *UserProfile) SetFieldsMap(m map[string]UserProfileCustomField) {
	profile.Fields.SetMap(m)
}
