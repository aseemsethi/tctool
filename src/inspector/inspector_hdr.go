package inspector

const (
	crUser                      = iota
	crArn                       = iota
	crUserCreationTime          = iota
	crPasswordEnabled           = iota
	crPasswordLastUsed          = iota
	crPasswordLastChanged       = iota
	crPasswordNextRotation      = iota
	crMfaActive                 = iota
	crAccessKey1Active          = iota
	crAccessKey1LastRotated     = iota
	crAccessKey1LastUsedDate    = iota
	crAccessKey1LastUsedRegion  = iota
	crAccessKey1LastUsedService = iota
	crAccessKey2Active          = iota
	crAccessKey2LastRotated     = iota
	crAccessKey2LastUsedDate    = iota
	crAccessKey2LastUsedRegion  = iota
	crAccessKey2LastUsedService = iota
	crCert1Active               = iota
	crCert1LastRotated          = iota
	crCert2Active               = iota
	crCert2LastRotated          = iota
)
