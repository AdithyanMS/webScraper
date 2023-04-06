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

func startPage() {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"),
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
		r.Headers.Set("Accept", "*/*")
		// r.Headers.Set(":authority", "www.animeout.xyz")
		// r.Headers.Set(":method", "POST")
		// r.Headers.Set(":path", "/tokyo-revengers-1080p-300mb720p-150mbepisode-1/")
		// r.Headers.Set(":scheme", "https")
		r.Headers.Set("accept-encoding", "gzip, deflate, br")
		r.Headers.Set("accept-language", "en-GB,en;q=0.6")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("content-type", "application/x-www-form-urlencoded")
		r.Headers.Set("origin", "https://www.animeout.xyz")
		// r.Headers.Set("sec-ch-ua", "Brave";v="111", "Not(A:Brand";v="8", "Chromium";v="111")
		r.Headers.Set("sec-ch-ua-platform", "Linux")
		r.Headers.Set("upgrade-insecure-requests", "1")
		r.Headers.Set("cookie", "PHPSESSID=egthsq71sj4rmi80hue3omdk5p; spu_closing_118900=true; cf_chl_2=6686f7cdf590d43; cf_clearance=xPCX71b70VvSnMGxxshtlFplg9XMgvuYazIJj1.MeC0-1680768349-0-160")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err.Error())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		c.DetectCharset = true
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
	c.Visit("https://www.animeout.xyz/tokyo-revengers-1080p-300mb720p-150mbepisode-1/")
	return
}

func main() {
	// fmt.Println(getAnimeListing())
	startPage()
}
