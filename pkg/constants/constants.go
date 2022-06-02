package constants

const (
	AVAILABLE XrefStatus = "AVAILABLE"
	USED      XrefStatus = "USED"
	NEW       XrefStatus = "new"
	EXISTING  XrefStatus = "existing"
)

type XrefStatus string
