package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gocolly/colly"
)

func getAnimeListing() (animes []string) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
		r.Headers.Set("Accept", "*/*")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err.Error())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("li", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, elem *colly.HTMLElement) {
			animes = append(animes, elem.Text)
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
	c.Visit("https://www.animeout.xyz/download-anime/")
	return
}

func downloadUrl(url string) (err error) {
	client := grab.NewClient()
	req, _ := grab.NewRequest("/home/adithyan/Videos", url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err = resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
	return
}

func main() {
	fmt.Println(getAnimeListing())
}
