package types

type Chat struct {
	Id              int    `json:"id"`
	Type            string `json:"type"`
	Title           string `json:"title,omitempty"`
	Username        string `json:"username,omitempty"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	IsForum         bool   `json:"is_forum,omitempty"`
	IsDirectMessage bool   `json:"is_direct_message,omitempty"`
}

// TODO: Complete the ChatFullInfo type
type ChatFullInfo struct {
	Chat
	AccentColorId        int                   `json:"accent_color_id"`
	MaxReactionCount     int                   `json:"max_reaction_count"`
	Photo                *ChatPhoto            `json:"photo,omitempty"`
	ActiveUsernames      []string              `json:"active_usernames,omitempty"`
	Birthdate            *Birthdate            `json:"birthdate,omitempty"`
	BusinessIntro        *BusinessIntro        `json:"business_intro,omitempty"`
	BusinessLocation     *BusinessLocation     `json:"business_location,omitempty"`
	BusinessOpeningHours *BusinessOpeningHours `json:"business_opening_hours,omitempty"`
	PersonalChat         *Chat                 `json:"personal_chat,omitempty"`
	ParentChat           *Chat                 `json:"parent_chat,omitempty"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}
