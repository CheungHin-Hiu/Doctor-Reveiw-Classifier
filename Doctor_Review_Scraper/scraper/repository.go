package scraper

import (

	"github.com/gocolly/colly/v2"
)

func GetDoctorInto(e *colly.HTMLElement) DoctorInfo{
	//  not implemented yet
	doctorInfo := DoctorInfo{};
	
	return doctorInfo;
}

func GetUserReviews(e *colly.HTMLElement) []UserReview {
	userReviews := []UserReview{}
}






