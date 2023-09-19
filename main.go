package main

import "log"

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer("localhost:2468", store)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
