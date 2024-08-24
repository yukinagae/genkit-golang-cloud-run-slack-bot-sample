package flow

import (
	"context"
	"fmt"
	"log"

	_ "embed"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/googleai"
	"github.com/invopop/jsonschema"
)

//go:embed answer.prompt
var promptTemplate string

type promptInput struct {
	Question string `json:"question"`
}

func DefineFlow(ctx context.Context) *genkit.Flow[string, string, struct{}] {
	// Initialize the Google AI plugin
	if err := googleai.Init(ctx, nil); err != nil {
		log.Fatalf("Failed to initialize Google AI plugin: %v", err)
	}

	model := googleai.Model("gemini-1.5-flash")

	answerPrompt, err := dotprompt.Define("answerPrompt",
		promptTemplate,
		dotprompt.Config{
			Model: model,
			Tools: []ai.Tool{},
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 1,
			},
			InputSchema:  jsonschema.Reflect(promptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	if err != nil {
		log.Fatalf("Failed to initialize prompt: %v", err)
	}

	return genkit.DefineFlow("answerFlow", func(ctx context.Context, input string) (string, error) {
		resp, err := answerPrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: &promptInput{
					Question: input,
				},
			},
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate answer: %w", err)
		}
		return resp.Text(), nil
	})
}
