package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func main() {
	// open json file
	jsonFile, err := os.Open("output.json")
	// if os.Open returns an error then print out it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened complex.json")
	// defer the closing of the json file
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var output Output
	err = json.Unmarshal(byteValue, &output)
	if err != nil {
		panic(err)
	}

	fmt.Println("Provider:", output.Provider)
	fmt.Println("Status:", output.Status)
	fmt.Println("Problems:", output.Problems)
	for _, result := range output.Results {

		fmt.Println("Name:", result.Name)
		fmt.Println()
		fmt.Println("Kind:", result.Kind)
		fmt.Println()
		fmt.Println("Details:", result.Details)
		fmt.Println()
		fmt.Println("ParentObject:", result.ParentObject)
		fmt.Println()
		//fmt.Println("Error:", result.Error)
		var resultErrors Errors
		err = json.Unmarshal(byteValue, &resultErrors)
		if err != nil {
			panic(err)
		}
		for _, k8sError := range result.Error {
			fmt.Println("Errors:", k8sError.Text)
			fmt.Println("KubernetesDoc:", k8sError.KubernetesDoc)
		}
		//for _, errDetail := range result.Result.Error {
		//	fmt.Println("Error Text:", errDetail.Text)
		//	for _, sensitive := range errDetail.Sensitive {
		//		fmt.Println("Unmasked:", sensitive.Unmasked, "Masked:", sensitive.Masked)
		//	}
		//}
		fmt.Println()
	}

}
