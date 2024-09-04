package main

import (
	"encoding/json"
	"fmt"
	"os"
	"server/handlers"
	"server/sources"
)

func main() {
	source := sources.NewSrtSource(os.Args[1])
	wordsMap, err := source.Source()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(wordsMap)
	dictapi := handlers.NewDictionaryAPIHandler()
	err = dictapi.Handle(wordsMap)
	jsonData, err := json.Marshal(wordsMap)
	if err != nil {
		fmt.Println("Error marshalling:", err)
		return
	}
	os.WriteFile("output.json", jsonData, 0644)
	// fmt.Println("Starting server at 127.0.0.1:8000...")
	// err := http.ListenAndServe("127.0.0.1:8000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
