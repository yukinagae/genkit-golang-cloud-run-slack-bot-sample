package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/genkit"
	"github.com/yukinagae/genkit-golang-cloud-run-slack-bot-sample/flow"
)

func main() {
	ctx := context.Background()

	_ = flow.DefineFlow(ctx)

	if err := genkit.Init(ctx, &genkit.Options{FlowAddr: ":3400"}); err != nil {
		log.Fatal(err)
	}
}
