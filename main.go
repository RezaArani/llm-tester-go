// Copyright (c) 2025 John Doe
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	aillm "github.com/RezaArani/aillm/controller"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func main() {

	log.Println("Starting LLM Testing...")
	startTime := time.Now()

	models := []string{
		"stablelm2",
		"llama3.2:1b",
		"llama3.1",
		"deepseek-r1:1.5b",
		"tinyllama",
		"llama3.2",
		"aya",
		"granite3-dense:8b",
		"mistral",
		"wizardlm2",
		"qwen2",
	}
	// Run queries
	queries := []string{
		"What is SemMapas? (limiting it to a maximum of 50 words)", // limit text for testing token count
		"Sum of SemMapas annual costs?",                            // me may keep model decide the result
		// "Say Hi",
	}
	results := runTests(models, queries)

	printResultsTable(results, "FullReport")

	log.Printf("Testing completed in %v\n", time.Since(startTime))
}

// Runs tests on the given models and collects results
func runTests(models []string, queries []string) []map[string]interface{} {
	var results []map[string]interface{}

	for _, model := range models {
		log.Printf("\n\nTesting model: %s", model)

		// Initialize LLM container
		llm := initLLMContainer(model)

		// Embed data
		embeddingTitle := "rawText"
		embeddingStart := time.Now()
		re := regexp.MustCompile(`[^a-zA-Z0-9:_-]`)
		indexKey := re.ReplaceAllString(model, "_")
		embedText(llm, indexKey, embeddingTitle)
		embeddingTime := time.Since(embeddingStart)

		log.Printf("Embedding completed in: %v", embeddingTime)

		for _, query := range queries {
			log.Println("Query: " + query)
			curTestResults := executeQuery(llm, model, query)
			curTestResults["embeddingTime"] = embeddingTime
			results = append(results, curTestResults)
		}

		// Remove embeddings and test without embeddings
		llm.RemoveEmbeddingDataFromRedis(embeddingTitle, llm.WithEmbeddingPrefix(indexKey))
		log.Println("Testing after removing embeddings...")
		for _, query := range queries {
			log.Println("Query(no embeddings) : " + query)

			curTestResults := executeQuery(llm, model, query)
			curTestResults["curTestResults"] = "(no embeddings) " + query
			curTestResults["embeddingTime"] = embeddingTime
			results = append(results, curTestResults)
		}
		printResultsTable(results, indexKey)

	}

	return results
}

// Initializes the LLM container for the given model
func initLLMContainer(model string) aillm.LLMContainer {
	client := &aillm.OllamaController{
		Config: aillm.LLMConfig{
			Apiurl:  "http://127.0.0.1:11434",
			AiModel: model,
		},
	}

	llm := aillm.LLMContainer{
		Embedder:       client,
		LLMClient:      client,
		ScoreThreshold: 0.9999, // We do not want to benchmark redis, so it will always include result.
		RedisClient: aillm.RedisClient{
			Host: "localhost:6379",
		},
	}
	llm.Init()
	return llm
}

// Embeds the given text into the LLM container
func embedText(llm aillm.LLMContainer, modelPrefix, title string) {
	content := map[string]aillm.LLMEmbeddingContent{
		"en": {Text: enRawText},
	}
	llm.EmbeddText(title, content, llm.WithEmbeddingPrefix(modelPrefix))
}

// Executes a single query and returns the result
func executeQuery(llm aillm.LLMContainer, model, query string) map[string]interface{} {
	re := regexp.MustCompile(`[^a-zA-Z0-9:_-]`)
	indexKey := re.ReplaceAllString(model, "_")
	//clearing user memory for fresh information
	llm.MemoryManager.DeleteMemory("User")

	start := time.Now()
	response, err := llm.AskLLM(query, llm.WithStreamingFunc(printStream), llm.WithEmbeddingPrefix(indexKey), llm.WithSessionID("User"))
	responseTime := time.Since(start)

	generatedPrompt := ""
	for _, prmpt := range response.Prompt {
		generatedPrompt += "\nRole: " + string(prmpt.Role) + "\n"
		for _, part := range prmpt.Parts {
			generatedPrompt += strings.ReplaceAll(part.(llms.TextContent).Text, "	", "")
		}
		// generatedPrompt+=
	}

	//factcheck
	if err != nil {
		log.Printf("Error querying model %s: %v", model, err)
		return map[string]interface{}{
			"Model":        model,
			"Query":        query,
			"ResponseTime": responseTime,
			"Error":        err.Error(),
		}
	}

	// Extract response details
	var completionTokens, promptTokens, totalTokens int
	genInfo := response.Response.Choices[0].GenerationInfo
	completionTokens = genInfo["CompletionTokens"].(int)
	promptTokens = genInfo["PromptTokens"].(int)
	totalTokens = genInfo["TotalTokens"].(int)
	compareArray := []string{enRawText, response.Response.Choices[0].Content}
	similarityResults, err := CalculateSentenceSimilarity(compareArray)
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	contextualReasoningCheckReport, err := contextualReasoningCheck(enRawText, query, response.Response.Choices[0].Content)
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	fmt.Println(similarityResults)

	return map[string]interface{}{
		"Model":                    model,
		"Query":                    query,
		"ResponseTime":             responseTime,
		"CompletionTokens":         completionTokens,
		"PromptTokens":             promptTokens,
		"TotalTokens":              totalTokens,
		"ReferenceDocs":            len(response.RagDocs.([]schema.Document)),
		"ResponseBody":             response.Response.Choices[0].Content,
		"GeneratedPrompt":          generatedPrompt,
		"TimeSpans":                response.Actions,
		"Similarity":               similarityResults[0],
		"contextualReasoningCheck": contextualReasoningCheckReport,
	}
}

// Prints real-time chunks of data returned from the LLM
func printStream(ctx context.Context, chunk []byte) error {
	fmt.Print(string(chunk))
	return nil
}

// Prints results of all tests in a table format
func printResultsTable(results []map[string]interface{}, fileName string) {
	file, err := os.OpenFile(fileName+"-"+time.Now().Format("2006-01-02-15-04")+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	embeddings := make(map[string]time.Duration)
	div := "---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"
	fmt.Fprintln(file, "\nResults Table:")
	fmt.Fprintln(file, div)
	fmt.Fprint(file, fmt.Sprintf("%-20s %-40s %-15s %-10s %-16s %-16s %-16s %-20s %-30s\n", "Model", "Query", "ResponseTime", "CONTEXT", "PromptTokens", "CompletionTokens", "TotalTokens", "Cosine Similarity", "Contextual Reasoning,Correct?"))
	fmt.Fprintln(file, div)
	LLMResponses := ""

	for _, result := range results {
		embeddings[result["Model"].(string)] = result["embeddingTime"].(time.Duration)

		if result == nil {
			continue
		}
		query := result["Query"].(string)
		if len(query) > 35 {
			query = query[:35] + "..."
		}
		shortCRC := result["contextualReasoningCheck"].(string)[:strings.Index(result["contextualReasoningCheck"].(string), " ")]
		fmt.Fprintln(file, fmt.Sprintf(
			"%-20s %-40s %-15v %-10v %-16v %-16v %-18v %-20v %-30v\n",
			result["Model"],
			query,
			result["ResponseTime"],
			result["ReferenceDocs"],
			result["PromptTokens"],
			result["CompletionTokens"],
			result["TotalTokens"],
			result["Similarity"],
			shortCRC,
		))

		LLMResponses += "\nModel:" + result["Model"].(string) + "\nQuery:" + result["Query"].(string) + "\n" + getActionTable(result["TimeSpans"].([]aillm.LLMAction)) + "\n\nRAG Promot:" + result["GeneratedPrompt"].(string) + "\n\nLLM Replied:" + result["ResponseBody"].(string) + "\n\nContextual Reasoning Check(Is it correct?):" + result["contextualReasoningCheck"].(string) + "\n" + div + "\n"

	}
	fmt.Fprintln(file, "\n====================================================== Embedding Time =====================================================================")

	for modelName, EmbeddingTime := range embeddings {
		fmt.Fprintln(file, fmt.Sprintf("Model: %v\t Embedding Time:%v", modelName, EmbeddingTime))
	}
	fmt.Fprintln(file, fmt.Sprintf("Context size: %v bytes", len(enRawText)))
	fmt.Fprintln(file, "Context: "+enRawText)
	fmt.Fprintln(file, "\n====================================================== REPORT DETAILS =====================================================================")

	fmt.Fprintln(file, LLMResponses)
}

// Prints the formatted table
func getActionTable(data []aillm.LLMAction) string {
	div := "\n-------------------------------------------------------------------------------------------------------------\n"
	result := div
	result += fmt.Sprintf("%-15s", "")
	result += fmt.Sprintf("%-20s", "Vector Search")
	result += fmt.Sprintf("%-20s", "LLM Call")
	result += fmt.Sprintf("%-20s", "First Chunk")
	result += fmt.Sprintf("%-20s", "Stream Generation Time")

	result += div

	// Print row data
	result += fmt.Sprintf("%-15s", "Seconds:")

	for i := 1; i < len(data)-1; i++ {
		duration := data[i+1].TimeStamp.Sub(data[i].TimeStamp).Seconds()
		result += fmt.Sprintf("%-20.1f", duration)
	}

	result += div
	for _, action := range data {
		result += action.Action.(string) + ": " + action.TimeStamp.Format("2006-01-02 15:04:05.000") + "\n"
	}
	result += div
	return result
}

// Example raw text to be embedded
const enRawText = `Welcome to SemMapas, your strategic partner in enhancing local engagement and tourism development. Designed specifically for businesses and municipalities, SemMapas offers a powerful platform to connect with residents and visitors alike, driving growth and prosperity in your community.
Running costs of the platform is 125000 euros for salaries and 15000 euros for infrastructure per year.
Our platform goes beyond traditional mapping services, offering a comprehensive suite of features tailored to meet the diverse needs of event organizers and businesses alike. From tourism guides to event navigation, SemMapas empowers you to create immersive experiences that captivate your audience and enhance their journey.
Our project has been launched since 2023 in Portugal.`
