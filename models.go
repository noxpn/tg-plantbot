package main

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Id       int      `json:"message_id"`
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Photo    []Photo  `json:"photo"`
	Document Document `json:"document"`
}

type Chat struct {
	Id   int    `json:"id"`
	Name string `json:"username"`
}

type Photo struct {
	Id        string `json:"file_id"`
	Unique_id string `json:"file_unique_id"`
	Size      int    `json:"file_size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Document struct {
	Id        string    `json:"file_id"`
	Unique_id string    `json:"file_unique_id"`
	Name      string    `json:"file_name"`
	Size      int       `json:"file_size"`
	Type      string    `json:"mime_type"`
	Thumb     PhotoSize `json:"thumb"`
}

func (d *Document) isEmpty() bool {
	return len(d.Id) == 0
}

type PhotoSize struct {
	Id             string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}
type RestResponse struct {
	Result []Update `json:"result"`
}

type BotRespond struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

