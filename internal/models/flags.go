package models

type Config struct {
	Interactive bool
	Mode        string // "auto" | "manual"
	ProbeOnly   bool
	Discover    bool

	SelectCache int
	AutoCache   bool
	UseCache    bool
	ForgetCache string // "", "all", or IP
	ListCache   bool

	DeepSearch bool
	Subnet     string

	TIP      string // TV IP
	TPort    string // TV SOAP port
	TPath    string // SOAP path
	TVVendor string // TV vendor

	LIP       string // local IP
	LFile     string // local file path (used only for MediaURL)
	LDir      string // directory to serve
	ServePort string // local HTTP port

	CachedConnMgrURL string
	CachedControlURL string
	ServerUp         bool
}

var DefaultConfig = Config{
	Interactive: false,
	SelectCache: -1,
	AutoCache:   false,
	UseCache:    true,
	ProbeOnly:   false,
	Discover:    false,
	Mode:        "auto",
	ServePort:   "8000",
	LDir:        "./directory",
}
