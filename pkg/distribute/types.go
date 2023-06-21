package distribute

type Config struct {
	ExcelConfig     ExcelConfig
	EmailConfig     EmailConfig
	MailSubject     string
	MailBody        string
	DocumentDir     string
	FileTypes       []FileType
	RecursiveSearch bool
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func (c Config) AsEmailMessage(to string, document_paths []string) *Message {
	message := Message{
		To:      []string{to},
		Subject: c.MailSubject,
		Body:    c.MailBody,
	}
	for _, document_path := range document_paths {
		message.AttachFile(document_path)
	}
	return &message
}



type EmailConfig struct {
	SenderEmail  string
	SMTPHost     string 
	SMTPPort     string
	Username string 
	Password string 
}

type ExcelConfig struct {
	FilePath     string
	FirstNameCol uint64
	LastNameCol  *uint64
	EmailCol     uint64
	SheetName    *string
	StartRow     *uint64
	EndRow       *uint64
}

type LexedData struct {
	NameEmails []NameEmail
	Anomalies  []LexerAnomaly
}

type ParsedData struct {
	NameEmailsToDocuments map[NameEmail][]DocumentPath
	Anomalies             []ParserAnomaly
}

type DocumentPath = string

type NameEmail struct {
	Name  string
	Email string
}

type LexerAnomaly struct {
	Type        LexerAnomalyType
	LineNumber  uint64
	LineContent []string
}

type LexerAnomalyType int

const (
	DuplicateNames LexerAnomalyType = iota
	DuplicateEmails
	NamesWithNoEmails
	EmailsWithNoNames
	NoEmailAndMail
	WrongEmailFormat
)

type ParserAnomaly struct {
	Type ParserAnomalyType
	Info []string
}

type ParserAnomalyType int

const (
	NoNameMatchForDocument ParserAnomalyType = iota
	NoDocumentMatchForName
	MultipleDocumentMatchesForName
)

type FileType = string

const (
	PDF FileType = ".pdf"
	TXT FileType = ".txt"
)

type Void struct{}
