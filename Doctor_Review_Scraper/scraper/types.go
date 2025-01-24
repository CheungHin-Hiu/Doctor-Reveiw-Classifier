package scraper

type UserReview struct {
	TreatmentPerformance, ServiceAttitude, Professionalism, ReasonablePrice int
	Author, VisitReason, CommentDate, VisitDate                             string
	Heading, Comment                                                        string
	MedicalFee                                                              int
}
