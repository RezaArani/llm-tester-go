package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

// API endpoint URL
const apiURL = "https://bge-base-en-v1-5.endpoints.kepler.ai.cloud.ovh.net/api/text2vec"

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(vec1, vec2 []float64) float64 {
	var dotProduct, normVec1, normVec2 float64
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		normVec1 += vec1[i] * vec1[i]
		normVec2 += vec2[i] * vec2[i]
	}
	return dotProduct / (math.Sqrt(normVec1) * math.Sqrt(normVec2))
}

// getEmbedding fetches the embedding vector for a given text from the API.
func getEmbedding(text, token string) ([]float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(text)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embedding []float64
	if err := json.Unmarshal(body, &embedding); err != nil {
		return nil, err
	}

	return embedding, nil
}

// CalculateSentenceSimilarity generates embeddings for sentences and computes their cosine similarities.
func CalculateSentenceSimilarity(texts []string) ([]float64, error) {
	token := os.Getenv("APITOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing API token in environment variables")
	}

	var embeddings [][]float64
	var similarityResults []float64

	for i, text := range texts {
		embedding, err := getEmbedding(text, token)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, embedding)

		if i > 0 {
			similarity := cosineSimilarity(embeddings[0], embedding)
			similarityResults = append(similarityResults,similarity )
			// similarityResults[fmt.Sprintf("Similarity with sentence n°%d", i)] = fmt.Sprintf("%.3f", similarity)
		}
	}

	return similarityResults, nil
}
