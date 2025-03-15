package model

import "gorm.io/gorm"

// Banner represents a banner to be displayed in Danta.

type Banner struct {

	// Content is the content of the banner.
	Content string

	// ApplicantEmail is the email of the applicant who submitted the banner.
	ApplicantEmail string

	// Status is the status of the banner.
	// It can be one of the following values:
	//  - "pending"
	//  - "approved"
	//  - "disapproved"
	Status string

	// gorm.Model provides fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	gorm.Model
}
