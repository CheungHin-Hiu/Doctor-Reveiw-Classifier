package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var HEADERS = []string{
	"DocChineseName", "DocEnglishName", "DocRegisteredSpecialty",
	"Author", " WaitingTime", "ConsultationPurpose", "Symptoms", "Fee", "VisitReason", "ConsultationDate",
	"Sentiment", "Experience", "Recommendation", "OtherComment",
	"OverallScore", "Professionalism", "TargetedTreatment", "ReasonablePrice", "WorthRecommending",
}

// Extract content between the starting substring and ending substring
func extractContent(text, start, end string) string {
	startIndex := strings.Index(text, start) //Return the index of the first instance of the substring in text
	if startIndex == -1 {                    //If we cannot find the starting substring, return an empty string
		return ""
	}
	startIndex += len(start) //Go to the char after the start substring

	if end == "" { //Go all char after the starting substring
		return strings.TrimSpace(text[startIndex:])
	}

	endIndex := strings.Index(text[startIndex:], end)
	if endIndex == -1 {
		return strings.TrimSpace(text[startIndex:])
	}
	// Get the the substring between the starting substring and ending substring
	return strings.TrimSpace(text[startIndex : startIndex+endIndex])
}

// Convert the number extracted in /image/comment?.gif to the words that the gif show
func worthRecommendingCovert(score string) string {
	switch score {
	case "3":
		return "值得推介"
	case "2":
		return "可以一試"
	default:
		return "不予置評"
	}
}

// Create a new collector
func newCollector(allowedDomain, userAgent string, timeout time.Duration) *colly.Collector {
	collector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)

	collector.SetRequestTimeout(timeout)
	collector.UserAgent = userAgent

	return collector
}

func writeReviewsToCSV(fileName string, review DoctorReview) {
	// Check if the file exists or not
	_, err := os.Stat(fileName)
	fileExist := err == nil

	// Open the file in append mode if the file exists, else create the new file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open or create the output csv file", err)
	}

	defer file.Close()

	// Initialize the file writer
	writer := csv.NewWriter(file)

	// Write the headers if the file is newly created
	if !fileExist {
		writer.Write(HEADERS)
	}

	// Write each user review to the CSV file
	for _, userReview := range review.UserReviews {
		record := []string{
			review.Info.ChineseName,
			review.Info.EnglishName,
			review.Info.RegisteredSpecialty,
			userReview.Author,
			userReview.WaitingTime,
			userReview.ConsultationPurpose,
			userReview.Symptoms,
			userReview.Fee,
			userReview.VisitReason,
			userReview.ConsultationDate,
			userReview.Sentiment,
			userReview.Experience,
			userReview.Recommendation,
			userReview.OtherComment,
			strconv.FormatFloat(userReview.OverallScore, 'f', -1, 64),
			strconv.Itoa(userReview.Professionalism),
			strconv.Itoa(userReview.TargetedTreatment),
			strconv.Itoa(userReview.ReasonablePrice),
			userReview.WorthRecommending,
		}
		writer.Write(record)
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatalln("Error writing to CSV:", err)
	}
}

// Extract the doctor information from the HTMElement
func getDoctorInfo(e *colly.HTMLElement) DoctorInfo {
	doctorInfo := DoctorInfo{}

	// Scrape the target doctor basic information

	//Extract the chinese name
	chineseNameElem := e.ChildText("h1.f20.green.all5.left5 b")
	doctorInfo.ChineseName = strings.TrimSpace(chineseNameElem)

	//Extract the english name
	englishNameElem := e.ChildText("h2.green.f20.all5.left5 b")
	doctorInfo.EnglishName = strings.TrimSpace(englishNameElem)

	//Extract specialty
	specialtyElem := e.ChildText("span[itemprop='medicalSpecialty']")
	doctorInfo.RegisteredSpecialty = strings.TrimSpace(specialtyElem)

	return doctorInfo
}

// Extract the user review from the HTMLElement
func getUserReviews(e *colly.HTMLElement) []UserReview {
	userReviews := []UserReview{}

	//Get every user review in the page
	e.ForEach("div[itemscope][itemtype='http://schema.org/Review']", func(i int, h *colly.HTMLElement) {
		review := UserReview{}

		// Get the basic information of the review
		review.Author = h.ChildText("div#guest_detail strong span[itemprop='author']")
		review.ConsultationDate = extractContent(h.ChildText("div#guest_detail div:nth-child(3)"), "• 診症日期 : ", "")
		review.ConsultationPurpose = extractContent(h.ChildText("div#guest_detail div:nth-child(4)"), "• 求診目的 : ", "")
		review.Fee = extractContent(h.ChildText("div#guest_detail div:nth-child(5)"), "• 診金收費 : ", "")
		review.WaitingTime = extractContent(h.ChildText("div#guest_detail div:nth-child(6)"), "• 輪候時間 : ", "")
		review.VisitReason = extractContent(h.ChildText("div#guest_detail div:nth-child(7)"), "• 選尋原因 : ", "")

		//Get score and comment of the guest
		review.OverallScore, _ = strconv.ParseFloat(h.ChildText("div#guest_point div:nth-child(1) span[itemprop='ratingValue']"), 64)
		review.WorthRecommending = worthRecommendingCovert(extractContent(h.ChildAttr("div#guest_point div:nth-child(2) img", "src"), "/image/comment", ".gif"))
		review.Professionalism, _ = strconv.Atoi(h.ChildText("div#guest_point div:nth-child(3) span"))
		review.TargetedTreatment, _ = strconv.Atoi(h.ChildText("div#guest_point div:nth-child(4) span"))
		review.ReasonablePrice, _ = strconv.Atoi(h.ChildText("div#guest_point div:nth-child(5) span"))

		comment := h.ChildText("div#LayoutDiv1 span[itemprop='reviewBody']")

		// Extract the comment part
		review.OtherComment = extractContent(comment, "", "簡評:個人觀感- ")
		review.Sentiment = extractContent(comment, "個人觀感- ", "經歷過程- ")
		review.Experience = extractContent(comment, "經歷過程- ", "推介意見- ")
		review.Recommendation = extractContent(comment, "推介意見- ", "")

		//Append to the review list
		userReviews = append(userReviews, review)
	})

	return userReviews
}
