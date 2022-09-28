package main

import (
	"flag"

	"github.com/elieser9001/AndroidRaptor/manager"
)

func main() {
	var adminId int
	var apiBotKey string

	flag.IntVar(&adminId, "uid", 0, "Your telegram user id")
	flag.StringVar(&apiBotKey, "abk", "", "API Bot Key")
	flag.Parse()

	if adminId == 0 || len(apiBotKey) <= 0 {
		flag.Usage()
		return
	}

	manager.Start(int64(adminId), apiBotKey)
}
