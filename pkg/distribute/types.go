package distribute

type RawData struct {
	NameEmails []NameEmail
	Anomalies  []Anomaly
}

type NameEmail struct {
	Name  string
	Email string
}

type Anomaly struct {
	Type        AnomalyType
	LineNumber  uint64
	LineContent []string
}

type AnomalyType int

const (
	DuplicateNames AnomalyType = iota
	DuplicateEmails
	NamesWithNoEmails
	EmailsWithNoNames
	NoEmailAndMail
)

type Void struct{}
