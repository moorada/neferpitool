package notification

type TemplateData struct {
	H1            string
	TextStatus    string
	TextWhois     string
	HeadersStatus []string
	HeadersWhois  []string
	DatasStatus   [][]string
	DatasWhois    [][]string
	TextExpiry    string
	HeadersExpiry []string
	DatasExpiry   [][]string
}

type Request struct {
	From     string
	Password string
	To       []string
	Subject  string
	body     string
}
