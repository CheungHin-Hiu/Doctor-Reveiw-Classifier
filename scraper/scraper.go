package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	c := newCollector(
		"www.seedoctor.com.hk",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		120*time.Second,
	)

	scrapeDoctorPage := func(doctorID int) {
		var doctorInfo DoctorInfo
		var userReviews []UserReview
		var executed = false

		collector := c.Clone()

		collector.OnRequest(func(r *colly.Request) {
			fmt.Printf("Visiting doctor ID %d: %s\n", doctorID, r.URL)
		})

		collector.OnHTML("#dr_info", func(h *colly.HTMLElement) {
			if doctorInfo == (DoctorInfo{}) {
				doctorInfo = getDoctorInfo(h)
			}
		})

		collector.OnHTML("body", func(h *colly.HTMLElement) {
			retrievedReviews := getUserReviews(h)
			userReviews = append(userReviews, retrievedReviews...)
			nextPageURL := h.DOM.Find("input[value='下一頁']").Parent().AttrOr("href", "")
			if nextPageURL != "" {
				absoluteURL := h.Request.AbsoluteURL(nextPageURL)
				fmt.Printf("Found next page URL: %s\n", absoluteURL)
				err := h.Request.Visit(absoluteURL)
				if err != nil {
					fmt.Printf("Error visiting next page: %v\n", err)
				}
			}
		})

		collector.OnScraped(func(r *colly.Response) {
			// Only write to CSV when scraping is complete (no more pages)
			if !executed {
				review := DoctorReview{doctorInfo, userReviews}
				fmt.Println(len(review.UserReviews))
				writeReviewsToCSV("data.csv", review)
				executed = true
			}

		})

		url := fmt.Sprintf("https://www.seedoctor.com.hk/dr_detail-1.asp?dr_doctor=%d", doctorID)
		err := collector.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting doctor %d: %v\n", doctorID, err)
		}
		collector.Wait()
	}

	// Define the range of doctor IDs to scrape
	startID := 10
	endID := 5188 // Adjust this range as needed

	// Option 1: Sequential scraping
	for id := startID; id <= endID; id++ {
		scrapeDoctorPage(id)
		// Add a delay between requests to be polite to the server
		time.Sleep(2 * time.Second)
	}
}
