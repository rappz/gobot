package bot

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var KING_IMAGE = []string{
	"https://media1.tenor.com/m/zzh5EGMb8KcAAAAd/yes-king.gif",
	"https://media1.tenor.com/m/wzBvSvmdhdMAAAAd/yes-king-yes.gif",
	"https://media1.tenor.com/m/1exE1H-iGGsAAAAd/martene3-yesking.gif",
	"https://media1.tenor.com/m/psMStUrhCp4AAAAd/burger-king-yes-sir.gif",
	"https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExbm1zNGNwbTFtbjFvbXlmMzByY3p3azhmNzN3cDN5d2FhMjRrcTF0eSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/L0nhTaYf038ZqDMZpY/giphy.gif",
}

func getKingImage() (*bytes.Reader, error) {
	resp, err := http.Get(KING_IMAGE[rand.Intn(len(KING_IMAGE))])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewReader(body), err
}
