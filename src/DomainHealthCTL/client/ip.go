package client

import (
	"github.com/ipipdotnet/ipdb-go"
	"log"
)

func IpInfo(address string) (map[string]string, error) {
	db, err := ipdb.NewCity("E:\\code\\xxscloud\\GoCtl\\src\\DomainHealthCTL\\assets\\ipipfree.ipdb")
	if err != nil {
		log.Fatal(err)
	}
	return db.FindMap(address, "CN")
}
