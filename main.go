package main

import (
	"fmt"
	"log"
	"model"
	"origin"
)

func main() {
	db, err := origin.Connect("recordings")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
	lists, err := db.Select("SELECT * FROM album WHERE artist = ?", "John Coltrane")
	if nil != err {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", lists)

	newalbum, newerr := db.ItemByID(2)
	if newerr != nil {
		log.Fatal(newerr)
	}
	fmt.Printf("newalbum found: %v\n", newalbum)

	albID, err := db.Add(model.Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}
