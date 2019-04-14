package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type People struct {
	Host      string `json:"host"`
	User      string `json:"user"`
	Password  string `json:"password"`
	PublicKey string `json:"publickey"`
	Port      int    `json:"port"`
}

func ReadConfig(pwdfile, group string) []People {
	var people []People

	f, err := os.Open(pwdfile)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonByte, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	info := gjson.Get(string(jsonByte), group)
	err = json.Unmarshal([]byte(info.String()), &people)
	if err != nil {
		log.Fatal(err)
	}
	return people
}
