package settings

import (
	"encoding/json"
	"io/ioutil"
)

func (s *Settings) Save() {
	data, _ := json.MarshalIndent(s, "", "\t")
	_ = ioutil.WriteFile("../settings.json", data, 0644)
}
