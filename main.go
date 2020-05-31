package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sclevine/agouti"
)

func main() {
	var url string
	fmt.Print("Enter site url: ")
	fmt.Scanf("%s\n", &url)
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox", "--log-level=3"}),
	)

	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	page, err := driver.NewPage()
	if err != nil {
		log.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate(url); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	vid, err := page.Find(`video`).Attribute("src")
	log.Println(vid == "")
	if err != nil {
		log.Fatal("Failed search source tag")
	}
	if vid == "" {
		log.Println("Failed search video tag. Trying searcg source tag")
		vid, err := page.Find(`source`).Attribute("src")
		if err != nil {
			log.Fatal("Failed search source tag")
		}
		log.Println(vid)
		download(vid)
	} else {
		log.Println(vid)
		download(vid)
	}

	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}

func download(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var filename = getFileName(url)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, res.Body)

		fmt.Println("Saved file ", filename)
	} else {
		fmt.Println(filename, " already exist")
	}
}

func getFileName(url string) string {
	hasher := md5.New()
	io.WriteString(hasher, url)

	var hashStr = hex.EncodeToString(hasher.Sum(nil)) + ".mp4"

	fmt.Printf("%#v", hashStr)
	return hashStr
}
