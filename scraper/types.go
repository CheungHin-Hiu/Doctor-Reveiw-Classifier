package main

type DoctorInfo struct {
	ChineseName, EnglishName, RegisteredSpecialty string
}

type UserReview struct {
	OverallScore                                                         float64
	Professionalism, TargetedTreatment, ReasonablePrice                  int
	WorthRecommending                                                    string
	ConsultationDate                                                     string
	Author, WaitingTime, ConsultationPurpose, Symptoms, Fee, VisitReason string
	Sentiment, Experience, Recommendation, OtherComment                  string
}

type DoctorReview struct {
	Info        DoctorInfo
	UserReviews []UserReview
}
