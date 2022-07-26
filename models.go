package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type BotSendPhoto struct {
	ChatId int    `json:"chat_id"`
	Photo  string `json:"photo"`
}

type PhotoSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
