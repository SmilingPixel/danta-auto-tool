package service


// DantaServiceIntf defines the interface for DantaService.
type DantaServiceIntf interface {
}


// DantaService provides methods to handle business logic related to Danta.
type DantaService struct {
	// larkDocService is used to interact with Lark documents
	larkDocService LarkDocServiceIntf
}


// NewDantaService creates a new instance of DantaService.
// It takes a LarkDocServiceIntf as a parameter and returns a pointer to DantaService.
func NewDantaService(
	larkDocService LarkDocServiceIntf,
) *DantaService {
	return &DantaService{
		larkDocService: larkDocService,
		
	}
}


