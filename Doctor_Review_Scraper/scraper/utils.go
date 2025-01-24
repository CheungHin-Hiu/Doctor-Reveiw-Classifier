package scraper

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Remove the prefixing "$" sign and convert string into int
func convertMedicalFee(s string) int {
	// Handle empty string
	if s == "" {
		return -1
	}

	// Remove the preceding dollars sign in the string
	s, _ = strings.CutPrefix(s, "$")

	// Convert the string into int
	fee, err := strconv.Atoi(s)

	if err != nil {
		return -1
	} else {
		return fee
	}
}

// Function to count the star in different scoring categories
// Return the score in integer type
func starCounting(h *colly.HTMLElement) int {
	var counter int
	// Count the number of bright star
	h.ForEach(".star-box img", func(_ int, img *colly.HTMLElement) {
		if strings.Contains(img.Attr("src"), "star@3x.png") {
			counter++
		}
	})
	return counter
}

// Assign the score to associated struct field of UserReview according to the metricName
func assignReviewScore(review *UserReview, metricName string, starCount int) {
	switch metricName {
	case "醫療效果":
		review.TreatmentPerformance = starCount
	case "服務態度":
		review.ServiceAttitude = starCount
	case "專業知識":
		review.Professionalism = starCount
	case "費用恰當":
		review.ReasonablePrice = starCount
	}
}
