package  conf

import (
    "os"
	"log"
    "io/ioutil"
    "encoding/json"

    "github.com/tidwall/gjson"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type People struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func ReadConfig(group string) []People {
    var people []People

    f, err := os.Open("conf/info.json")
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
