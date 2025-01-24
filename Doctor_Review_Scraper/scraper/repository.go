package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

var HEADERS = []string{
	"TreatmentPerformance", "ServiceAttitude", "Professionalism", "ReasonablePrice",
	"Author", "VisitReason", " CommentDate", "VisitDate", "Heading", "Comment", "MedicalFee",
}

func NewCollector(allowedDomain, userAgent string, timeout time.Duration) *colly.Collector {
	collector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)

	collector.SetRequestTimeout(timeout)
	collector.UserAgent = userAgent

	return collector
}

// Extract the user review from the HTML element
func GetUserReviews(e *colly.HTMLElement) []UserReview {
	userReviews := []UserReview{}

	e.ForEach("div.doctor-info-item.d-flex", func(i int, h *colly.HTMLElement) {
		review := UserReview{} //Initialize a UserReview object

		// Get the author of the comment
		review.Author = h.ChildText("div.name.num-family-12px")
		fmt.Printf("Review author: %s \n", review.Author)

		// Get the visit reason of the appointment
		review.VisitReason = h.ChildText("a.reason-text")
		fmt.Printf("Visit reason: %s \n", review.VisitReason)

		// Get the comment date
		review.CommentDate = h.ChildText("div.row1.d-flex.position-relative span.num-family-12px.datetime")
		fmt.Printf("Comment Date: %s \n", review.CommentDate)

		// Get the title of the comment
		review.Heading = h.ChildText("span.title-s")
		fmt.Printf("Comment heading: %s \n", review.Heading)

		// Get the content of the comment
		review.Comment = h.ChildText("p.text-align-justify")
		fmt.Printf("Comment: %s \n", review.Comment)

		// Get the visit date of the comment
		review.VisitDate = h.ChildText("div.right-right div.other-box span.num-family-12px[style='display: inline-block;']")
		fmt.Printf("Visit Date: %s \n", review.VisitDate)

		// Get the medical fee
		review.MedicalFee = convertMedicalFee(h.ChildText("div.right-right div.other-box span.num-family-12px.ml-2"))
		fmt.Printf("Medical fee: %d \n", review.MedicalFee)

		// Get the score assigned by the commenter
		h.ForEach("ul.rating-box li.rating-box-item", func(j int, h2 *colly.HTMLElement) {
			metricName := h2.ChildText("span:first-of-type") // Access the first span tag existed

			starCount := starCounting(h2)

			assignReviewScore(&review, metricName, starCount)
		})

		fmt.Printf("TreatmentPerformance score: %d\n", review.TreatmentPerformance)
		fmt.Printf("ServiceAttitude score: %d\n", review.ServiceAttitude)
		fmt.Printf("Professionalism score: %d\n", review.Professionalism)
		fmt.Printf("ReasonablePrice score: %d\n", review.ReasonablePrice)

		userReviews = append(userReviews, review)
	})
	return userReviews
}

func GetNextPageUrl(h *colly.HTMLElement) string {
	var nextPageURL string = ""
	h.ForEach("div.fenye ul.mobilie-pagination li", func(i int, li *colly.HTMLElement) {
		// Get the text inside a href if it exists
		linkText := li.ChildText("a")

		if linkText == "下一頁" {
			nextPageURL = li.ChildAttr("a", "href")
			return
		}

	})
	return nextPageURL
}

func WriteToCSV(filename string, reviews []UserReview) {
	var fileExist bool = true
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fileExist = false
	}

	// Open the file in append mode if the file exists, create the file otherwise
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open or create the file")
	}

	defer file.Close()

	// Initialize the file writer
	writer := csv.NewWriter(file)

	// If the file is newly created, write the corresponding header to it
	if !fileExist {
		writer.Write(HEADERS)
	}

	// Write the collected user review to the CSV file
	for _, review := range reviews {
		record := []string{
			strconv.Itoa(review.TreatmentPerformance),
			strconv.Itoa(review.ServiceAttitude),
			strconv.Itoa(review.Professionalism),
			strconv.Itoa(review.ReasonablePrice),
			review.Author,
			review.VisitReason,
			review.CommentDate,
			review.VisitReason,
			review.Heading,
			review.Comment,
			strconv.Itoa(review.MedicalFee),
		}
		writer.Write(record)
	}

	// Ensure all data is already written to our CSV file
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatalln("Error writing to CSV:", err)
	}

}
