package constants

const (
	AVAILABLE   XrefStatus = "AVAILABLE"
	UNAVAILABLE XrefStatus = "UNAVAILABLE"
	USED        XrefStatus = "USED"
	NEW         XrefStatus = "new"
	EXISTING    XrefStatus = "existing"
)

type XrefStatus string
