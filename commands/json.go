package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getJSON(request string, i interface{}) error {
	resp, err := http.Get(request)

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Issue reading json")
		return err
	}
	defer resp.Body.Close()
	return json.Unmarshal(contents, &i)
}
