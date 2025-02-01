package types

type Conversation struct {
	ChatId  int64
	channel chan string
}

func NewConversation(chatId int64) *Conversation {
	return &Conversation{
		ChatId:  chatId,
		channel: make(chan string),
	}
}

func (c *Conversation) Close() {
	close(c.channel)
}

func (c *Conversation) WaitAnswer() string {
	return <-c.channel
}

func (c *Conversation) Answer(text string) {
	c.channel <- text
}
