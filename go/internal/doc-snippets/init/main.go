// Copyright 2024 Google LLC
// SPDX-License-Identifier: Apache-2.0


// [START main]
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	// Import Genkit and the Google AI plugin
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googleai"
)

func main() {
	ctx := context.Background()

	g, err := genkit.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the Google AI plugin. When you pass nil for the
	// Config parameter, the Google AI plugin will get the API key from the
	// GOOGLE_GENAI_API_KEY environment variable, which is the recommended
	// practice.
	if err := googleai.Init(ctx, g, nil); err != nil {
		log.Fatal(err)
	}

	// Define a simple flow that prompts an LLM to generate menu suggestions.
	genkit.DefineFlow(g, "menuSuggestionFlow", func(ctx context.Context, input string) (string, error) {
		// The Google AI API provides access to several generative models. Here,
		// we specify gemini-1.5-flash.
		m := googleai.Model(g, "gemini-1.5-flash")
		if m == nil {
			return "", errors.New("menuSuggestionFlow: failed to find model")
		}

		// Construct a request and send it to the model API (Google AI).
		resp, err := genkit.Generate(ctx, g,
			ai.WithModel(m),
			ai.WithConfig(&ai.GenerationCommonConfig{Temperature: 1}),
			ai.WithTextPrompt(fmt.Sprintf(`Suggest an item for the menu of a %s themed restaurant`, input)))
		if err != nil {
			return "", err
		}

		// Handle the response from the model API. In this sample, we just
		// convert it to a string. but more complicated flows might coerce the
		// response into structured output or chain the response into another
		// LLM call.
		text := resp.Text()
		return text, nil
	})

	// Initialize Genkit and start a flow server. This call must come last,
	// after all of your plug-in configuration and flow definitions. When you
	// pass a nil configuration to Init, Genkit starts a local flow server,
	// which you can interact with using the developer UI.
	if err := g.Start(ctx, nil); err != nil {
		log.Fatal(err)
	}
}

// [END main]
