package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/midbel/uuid"
)

func main() {
	typ := flag.Int("t", 5, "version")
	ns := flag.String("n", "", "namespace")
	flag.Parse()

	uid := uuid.Nil
	switch *typ {
	case 1:
	case 4:
		uid = uuid.UUID4()
	case 3, 5:
		uid = uuid3or5(*ns, flag.Arg(0), *typ)
	default:
	}
	fmt.Println(uid.String())
}

func uuid3or5(ns, str string, version int) uuid.UUID {
	var uid uuid.UUID
	switch strings.ToLower(ns) {
	default:
		return uuid.Nil
	case "dns":
		uid = uuid.DNS
	case "oid":
		uid = uuid.OID
	case "url":
		uid = uuid.URL
	case "dn":
		uid = uuid.DN
	}
	switch version {
	default:
		return uuid.Nil
	case 3:
		return uuid.UUID3([]byte(str), uid)
	case 5:
		return uuid.UUID5([]byte(str), uid)
	}
}
