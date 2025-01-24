package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"doctor-review-scraper/scraper"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Get the starting and ending doctorID to parse
	startIdPtr := flag.Int("start_id", 6871, "The starting doctor page ID to start scraping")
	endIdPtr := flag.Int("end_id", 20000, "The ending doctor page ID to end scraping")
	flag.Parse()

	// Create a new gocolly collector
	c := scraper.NewCollector(
		"www.goodoctor.com.hk",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		120*time.Second,
	)

	// Set limit to the scraper
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*goodoctor.com.hk*",
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})

	// Scrape doctor comment page according to the page ID
	scrapeDoctorCommentPage := func(pageID int) {
		var userReviews []scraper.UserReview

		collector := c.Clone()
		visitUrl := fmt.Sprintf("https://www.goodoctor.com.hk/doctor/detail/%d/", pageID)

		// If visited a page, prints it URL
		collector.OnRequest(func(r *colly.Request) {
			fmt.Printf("Visiting the page: %s\n", r.URL)
		})

		// Add error callback
		collector.OnError(func(r *colly.Response, err error) {
			log.Printf("Error visiting %s: %v", r.Request.URL, err)
		})

		// Get the basic doctor information and user info
		collector.OnHTML("body", func(h *colly.HTMLElement) {
			reviews := scraper.GetUserReviews(h)
			userReviews = append(userReviews, reviews...)
			nextPageUrl := scraper.GetNextPageUrl(h)
			if nextPageUrl != "" {
				err := h.Request.Visit(nextPageUrl)
				if err != nil {
					fmt.Println("Error visiting next page: %v\n", err)
				}
			}
			fmt.Printf("Found %d reviews\n", len(reviews))

		})

		// Add error handling for Visit
		err := collector.Visit(visitUrl)
		if err != nil {
			log.Fatalf("Failed to visit page: %v", err)
		}

		collector.Wait()

		scraper.WriteToCSV("data/data.csv", userReviews)
	}

	startId := *startIdPtr
	endId := *endIdPtr

	for id := startId; id < endId; id++ {
		scrapeDoctorCommentPage(id)
		time.Sleep(2 * time.Second)
	}

}
