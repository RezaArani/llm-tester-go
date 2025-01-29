package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func contextualReasoningCheck(context, query, answer string) (string, error) {
	url := "https://llama-3-1-70b-instruct.endpoints.kepler.ai.cloud.ovh.net/api/openai_compat/v1/chat/completions"

	prompt :=  `Is this response "`+answer+`" is a correct answer to "`+query+`" about this context "`+context+`"? Say Yes or No  and why,very very shortly.`
	payload := map[string]interface{}{
		"max_tokens": 512,
		"messages": []map[string]interface{}{
			{
				"content": prompt,
				"name":    "User",
				"role":    "user",
			},
		},
		"model":       "Meta-Llama-3_1-70B-Instruct",
		"temperature": 0,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return "",err
	}

	accessToken :=os.Getenv("APITOKEN")
	if accessToken == "" {
		fmt.Println("Error: APITOKEN is not set")
		return "",err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "",err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "",err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return "",err
		}

		var responseData map[string]interface{}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return "",err
		}

		choices, ok := responseData["choices"].([]interface{})
		if !ok {
			fmt.Println("Error: choices field is not an array")
			return "",err
		}

		for _, choice := range choices {
			choiceMap, ok := choice.(map[string]interface{})
			if !ok {
				fmt.Println("Error: choice is not a map")
				return "",err
			}
			message, ok := choiceMap["message"].(map[string]interface{})
			if !ok {
				fmt.Println("Error: message field is not a map")
				return "",err
			}
			content, ok := message["content"].(string)
			if !ok {
				fmt.Println("Error: content is not a string")
				return "",err
			}

			return content,nil
		}
	} else {
		return "",err
	}
	return "",nil
}