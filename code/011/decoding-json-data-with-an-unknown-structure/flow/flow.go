package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdin)

	for {
		var v map[string]interface{}
		if err := decoder.Decode(&v); err != nil {
			log.Println(err.Error())
			return
		}

		for k := range v {
			if k != "Title" {
				v[k] = nil
			}
		}

		if err := encoder.Encode(&v); err != nil {
			log.Println(err.Error())
		}
	}
}
