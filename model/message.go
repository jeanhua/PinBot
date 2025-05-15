package model

type Message struct {
	Time     int    `json:"time"`
	PostType string `json:"post_type"`
	SelfId   int    `json:"self_id"`
}

type FriendMessage struct {
	Time        int           `json:"time"`
	PostType    string        `json:"post_type"`
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageId   int           `json:"message_id"`
	UserId      int           `json:"user_id"`
	Message     []OB11Segment `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      FriendSender  `json:"sender"`
	SelfId      int           `json:"self_id"`
}

type GroupMessage struct {
	Time        int           `json:"time"`
	PostType    string        `json:"post_type"`
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageId   int           `json:"message_id"`
	UserId      int           `json:"user_id"`
	GroupId     int           `json:"group_id"`
	Message     []OB11Segment `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      GroupSender   `json:"sender"`
	SelfId      int           `json:"self_id"`
}

type OB11Segment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type FriendSender struct {
	UserId   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
}

type GroupSender struct {
	UserId   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}
