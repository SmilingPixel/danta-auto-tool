package service

import (

	"github.com/rs/zerolog/log"
)

// DantaServiceIntf defines the interface for DantaService.
type DantaServiceIntf interface {
	// UpdateBannerAndNotify updates the banner and notifies the applicants.
	UpdateBannerAndNotify(content string, toEmailList []string) error
}


// DantaService provides methods to handle business logic related to Danta.
type DantaService struct {
	// larkDocService is used to interact with Lark documents
	larkDocService LarkDocServiceIntf

	// larkEmailService is used to interact with Lark emails
	larkEmailService LarkEmailServiceIntf
}


// NewDantaService creates a new instance of DantaService.
// It takes a LarkDocServiceIntf as a parameter and returns a pointer to DantaService.
func NewDantaService(
	larkDocService LarkDocServiceIntf,
	larkEmailService LarkEmailServiceIntf,
) *DantaService {
	return &DantaService{
		larkDocService: larkDocService,
		larkEmailService: larkEmailService,
	}
}

// UpdateBannerAndNotify do the following things:
//  1. Edit banner config file (in Github repo)
//  2. Send email to applicants
func (s *DantaService) UpdateBannerAndNotify(
	content string,
	toEmailList []string,
) error {
	log.Info().Msgf("[UpdateBannerAndNotify] Start updating banner and notifying applicants, content: %+v, toEmailList: %+v", content, toEmailList)



	for _, email := range toEmailList {
		log.Info().Msgf("[UpdateBannerAndNotify] Sending email to: %s", email)
		// TODO: Send email @xunzhou24
	}

	return nil
}


