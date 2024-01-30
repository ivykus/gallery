package main

import (
	"fmt"

	"github.com/ivykus/gallery/models"
)

// Host	sandbox-smtp.mailcatch.app
// Port	25, 1025, 2525
// Username	c8648ea7b08f
// Password	711d678a04e4

const (
	host     = "sandbox-smtp.mailcatch.app"
	port     = 1025
	username = "c8648ea7b08f"
	password = "711d678a04e4"
)

func main() {
	gs := models.GalleryService{}
	images, err := gs.Images(5)
	if err != nil {
		panic(err)
	}
	fmt.Println(images)

}
