package gpt

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func System(content string) Message    { return Message{"system", content} }
func User(content string) Message      { return Message{"user", content} }
func Assistant(content string) Message { return Message{"assistant", content} }

type Options struct {
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

type Wrapper interface {
	Init(Options)
	Complete([]Message) (Message, error)
}
