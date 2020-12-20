package main

import (
	"flag"
)

var (
	idmap = flag.String("idmap", "message/message_id.proto", "specify message_id.proto path")
	genId = flag.Bool("genId", false, "gen idmap.json")

	idmapFile = "idmap.json"
)

func main() {
	flag.Parse()
	if *genId {
		ParseIdMap()
	}
	Init()
}
