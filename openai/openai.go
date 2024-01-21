package openai

import (
	"context"

	gpt "github.com/alarbada/chatgptwrapper"
	"github.com/sashabaranov/go-openai"
)

type client struct {
	options gpt.Options
	client  *openai.Client
}

func (this *client) Init(options gpt.Options) {
	this.options = options
}

func mapOpenaiMsgs(msgs []gpt.Message) []openai.ChatCompletionMessage {
	openaiMsgs := make([]openai.ChatCompletionMessage, len(msgs))
	for i, msg := range msgs {
		openaiMsgs[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return openaiMsgs
}

func (this *client) Complete(msgs []gpt.Message) (msg gpt.Message, err error) {
	res, err := this.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       this.options.Model,
		Messages:    mapOpenaiMsgs(msgs),
		MaxTokens:   this.options.MaxTokens,
		Temperature: this.options.Temperature,
	})
	if err != nil {
		return msg, err
	}

	assistantMsg := gpt.Message{"assistant", res.Choices[0].Message.Content}
	return assistantMsg, nil
}

func New(apikey string) gpt.Wrapper {
	return &client{
		options: gpt.Options{
			Model:       "gpt-4",
			MaxTokens:   1500,
			Temperature: 0,
		},
		client: openai.NewClient(apikey),
	}
}
