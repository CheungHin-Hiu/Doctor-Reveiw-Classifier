package scraper

type DoctorInfo struct {
	Name, Specialty string
}

type UserReview struct {
	TreatmentOutput, CustomerService, Professionalism, ReasonablePrice int
	Author, VisitReason, CommentDate, VisitDate                        string
	Heading, Comment                                                   string
	MedicalFee                                                         int
}

type DoctorReview struct {
	DocInfo DoctorInfo
	Reviews []UserReview
}
