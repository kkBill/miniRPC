package main

import (
	"encoding/gob"


)
func main() {
	gob.Register(public.ResponseQueryUser{})
}
