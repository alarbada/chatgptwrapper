package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Message struct {
	Role, Content string
}

func System(content string) Message    { return Message{"system", content} }
func User(content string) Message      { return Message{"user", content} }
func Assistant(content string) Message { return Message{"assistant", content} }

type Wrapper interface {
	Model(string) Wrapper
	MaxTokens(int) Wrapper
	Temperature(float32) Wrapper

	Complete(msgs []Message) (Message, error)
}

type openaiWrapper struct {
	req    openai.ChatCompletionRequest
	client *openai.Client
}

var defaultRequest = openai.ChatCompletionRequest{
	Model:       openai.GPT4TurboPreview,
	MaxTokens:   1500,
	Temperature: 0,
}

func newOpenaiWrapper(apikey string) *openaiWrapper {
	return &openaiWrapper{defaultRequest, openai.NewClient(apikey)}
}

func (this *openaiWrapper) Model(model string) Wrapper {
	this.req.Model = model
	return this
}

func (this *openaiWrapper) MaxTokens(maxTokens int) Wrapper {
	this.req.MaxTokens = maxTokens
	return this
}

func (this *openaiWrapper) Temperature(t float32) Wrapper {
	this.req.Temperature = t
	return this
}

func (this openaiWrapper) Complete(msgs []Message) (msg Message, err error) {
	res, err := this.client.CreateChatCompletion(context.Background(), this.req)
	if err != nil {
		return msg, err
	}

	assistantMsg := Message{"assistant", res.Choices[0].Message.Content}
	return assistantMsg, nil
}

func NewOpenAI(apikey string) Wrapper {
	return newOpenaiWrapper(apikey)
}
