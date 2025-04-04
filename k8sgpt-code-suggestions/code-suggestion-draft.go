package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//type Output struct {
//	Provider string `json:"provider"`
//	Errors   any    `json:"errors"`
//	Status   string `json:"status"`
//	Problems int    `json:"problems"`
//	Results  []struct {
//		Result struct {
//			Kind  string `json:"kind"`
//			Name  string `json:"name"`
//			Error []struct {
//				Text          string `json:"Text"`
//				KubernetesDoc string `json:"KubernetesDoc"`
//				Sensitive     []struct {
//					Unmasked string `json:"Unmasked"`
//					Masked   string `json:"Masked"`
//				} `json:"sensitive"`
//			} `json:"error"`
//			Details      string `json:"details"`
//			ParentObject string `json:"parentObject"`
//		} `json:"result"`
//	} `json:"results"`
//}

type Output struct {
	Provider string    `json:"provider"`
	Errors   any       `json:"errors"`
	Status   string    `json:"status"`
	Problems int       `json:"problems"`
	Results  []Results `json:"results"`
}
type Sensitive struct {
	Unmasked string `json:"Unmasked"`
	Masked   string `json:"Masked"`
}
type Errors struct {
	Text          string      `json:"Text"`
	KubernetesDoc string      `json:"KubernetesDoc"`
	Sensitive     []Sensitive `json:"Sensitive"`
}
type Results struct {
	Kind         string   `json:"kind"`
	Name         string   `json:"name"`
	Error        []Errors `json:"error"`
	Details      string   `json:"details"`
	ParentObject string   `json:"parentObject"`
}

type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func main() {
	// open json file
	jsonFile, err := os.Open("output.json")
	// if os.Open returns an error then print out it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened output.json")
	// defer the closing of the json file
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var output Output
	err = json.Unmarshal(byteValue, &output)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Provider:", output.Provider)
	//fmt.Println("Status:", output.Status)
	//fmt.Println("Problems:", output.Problems)
	for _, result := range output.Results {
		//fmt.Println("Name:", result.Name)
		//fmt.Println()
		//fmt.Println("Kind:", result.Kind)
		//fmt.Println()
		//fmt.Println("Details:", result.Details)
		//fmt.Println()
		//fmt.Println("ParentObject:", result.ParentObject)
		//fmt.Println()
		//fmt.Println("Error:", result.Error)
		var resourceName string
		if result.ParentObject != "" {
			resourceName = result.ParentObject
		} else {
			resourceName = result.Name
		}

		// Split the string by "/"
		resourceNameParts := strings.Split(resourceName, "/")
		// Check if we have at least two parts
		if len(resourceNameParts) < 2 {
			fmt.Println("The string does not contain a '/' to separate.")
			return
		}

		// Assign the two parts
		//namespace := resourceNameParts[0]
		resourceName = resourceNameParts[1]

		// Print the results
		//fmt.Println("First part:", namespace)
		//fmt.Println("Second part:", resourceName)

		var resultErrors Errors
		err = json.Unmarshal(byteValue, &resultErrors)
		if err != nil {
			panic(err)
		}
		var errorText string
		for _, k8sError := range result.Error {
			//fmt.Println("Errors:", k8sError.Text)
			errorText = k8sError.Text
			//fmt.Println("KubernetesDoc:", k8sError.KubernetesDoc)
		}

		// Define the starting directory and the search string.
		root := "../k8s-bugs"        // current directory; change as needed
		searchString := resourceName // the string to search for
		var yamlWitherror string

		// Walk through all files starting from the root directory.
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err // if there’s an error accessing the path, abort.
			}
			if info.IsDir() {
				return nil // skip directories
			}

			// Open the file.
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Read the entire file.
			data, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			// Check if the file contains the search string.
			var fileWithError string
			fmt.Println("Searching folder ", root, "for a k8s manifest containing string:", searchString)
			if strings.Contains(string(data), searchString) {
				fmt.Println("Found string in file:", path)
				fileWithError = path
			}

			// Open the file
			if fileWithError != "" {
				file, err = os.Open(fileWithError)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return nil
				}
				defer file.Close() // Ensure the file is closed when the function exits

				// Read the file contents
				data, err = io.ReadAll(file)
				if err != nil {
					fmt.Println("Error reading file:", err)
					return nil
				}

				// Print the file contents
				fmt.Println("This is the content of file:", string(fileWithError))
				fmt.Println(string(data))
				yamlWitherror = string(data)
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Error walking the path: %v\n", err)
		}

		// Define the API URL and payload
		url := "http://localhost:11434/api/generate"
		//fmt.Println(errorText)
		prompt := ` 

This is my kubernetes manifest:
` + yamlWitherror + `
This gives me the following error on my cluster: ` + errorText + `

Please fix the manifest for me. And only send the fixed yaml as the response.

		`
		fmt.Println("Sending the following prompt to", url, ":", prompt)
		payload := map[string]interface{}{
			"model":  "llama3:8b", // Replace with your model name
			"prompt": prompt,      // Replace with your prompt
			"stream": false,
		}

		// Convert the payload to JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		// Send the POST request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var apiResp OllamaResponse
		err = json.Unmarshal(responseData, &apiResp)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		// Now apiResp.Result contains the full output as a single string.
		finalOutput := apiResp.Response
		fmt.Println("AI Response:", finalOutput)
		fmt.Println("--------------------------------------------------------------------------------------------")
	}

}
