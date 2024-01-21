package togetherai

import (
	gpt "github.com/alarbada/chatgptwrapper"
	"github.com/alarbada/curly"
)

type object map[string]any

type client struct {
	APIKey  string
	options gpt.Options
}

func New(apikey string) gpt.Wrapper {
	return &client{
		APIKey: apikey,
		options: gpt.Options{
			Model:       "mistralai/Mixtral-8x7B-Instruct-v0.1",
			MaxTokens:   1500,
			Temperature: 0,
		},
	}
}

func (this *client) Init(options gpt.Options) {
	this.options = options
}

func (this *client) Complete(msgs []gpt.Message) (msg gpt.Message, err error) {
	payload := map[string]any{
		"model":       this.options.Model,
		"max_tokens":  this.options.MaxTokens,
		"stop":        []string{"</s>", "[/INST]"},
		"temperature": this.options.Temperature,
		"messages":    msgs,
	}

	c := curly.New("POST", "https://api.together.xyz/v1/chat/completions").
		Header("Authorization", "Bearer "+this.APIKey).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(payload)

	statusCode := c.Do()
	if statusCode.NotOk() {
		return msg, c.Error()
	}

	var res struct {
		ID      string `json:"id"`
		Choices []struct {
			gpt.Message `json:"message"`
		} `json:"choices"`
	}
	c.Unmarshal(&res)
	return res.Choices[0].Message, c.Error()
}

func (this *client) Embeddings(input string) ([]float64, error) {
	req := curly.New("POST", "https://api.together.xyz/v1/embeddings")
	req.Header("Authorization", "Bearer "+this.APIKey)
	req.Header("Content-Type", "application/json")
	req.Body(object{
		"model": "togethercomputer/m2-bert-80M-8k-retrieval",
		"input": input,
	})

	status := req.Do()
	if status.NotOk() {
		return nil, req.Error()
	}

	var res struct {
		Object string `json:"object"`
		Data   []struct {
			Object    string    `json:"object"`
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"data"`
	}
	req.Unmarshal(&res)

	return res.Data[0].Embedding, req.Error()
}
