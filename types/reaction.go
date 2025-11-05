package types

type IReactionType any

type ReactionTypeEmoji struct {
	// Type string
	Emoji string `json:"emoji"`
}

type ReactionTypeCustomEmoji struct {
	// Type string
	CustomEmojiId string `json:"custom_emoji_id"`
}

type ReactionTypePaid struct {
	Type string
}

type MessageReactionUpdated struct {
	Chat        Chat
	MessageId   int  `json:"message_id"`
	User        User `json:"user,omitempty"`
	ActorChat   Chat `json:"actor_chat,omitempty"`
	Date        int
	OldReaction []IReactionType `json:"old_reaction"`
	NewReaction []IReactionType `json:"new_reaction"`
}

type ReactionCount struct {
	Type       IReactionType
	TotalCount int `json:"total_count"`
}

type MessageReactionCountUpdated struct {
	Chat      Chat
	MessageId int `json:"message_id"`
	Date      int
	Reactions []ReactionCount `json:"reactions"`
}
