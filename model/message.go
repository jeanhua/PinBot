package model

type Message struct {
	Time     uint   `json:"time"`
	PostType string `json:"post_type"`
	SelfId   uint   `json:"self_id"`
}

type FriendMessage struct {
	Time        uint          `json:"time"`
	PostType    string        `json:"post_type"`
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageId   uint          `json:"message_id"`
	UserId      uint          `json:"user_id"`
	Message     []OB11Segment `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      FriendSender  `json:"sender"`
	SelfId      uint          `json:"self_id"`
}

type GroupMessage struct {
	Time        uint          `json:"time"`
	PostType    string        `json:"post_type"`
	MessageType string        `json:"message_type"`
	SubType     string        `json:"sub_type"`
	MessageId   uint          `json:"message_id"`
	UserId      uint          `json:"user_id"`
	GroupId     uint          `json:"group_id"`
	Message     []OB11Segment `json:"message"`
	RawMessage  string        `json:"raw_message"`
	Sender      GroupSender   `json:"sender"`
	SelfId      uint          `json:"self_id"`
}

type OB11Segment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type FriendSender struct {
	UserId   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
}

type GroupSender struct {
	UserId   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

type Response struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		MessageId int `json:"message_id"`
	} `json:"data"`
	Message string `json:"message"`
	Wording string `json:"wording"`
}
