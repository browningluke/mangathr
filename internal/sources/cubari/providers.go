package cubari

type Provider struct {
	name        string
	sourceURL   string
	regex       string
	registrable bool
}

var IMGUR = Provider{
	name:        "imgur",
	sourceURL:   "https://imgur.com/a",
	regex:       `^((http(s)?://)?imgur\.com/a/)([a-zA-Z0-9]{4,8})/?$`,
	registrable: false,
}

var GIST = Provider{
	name:        "gist",
	sourceURL:   "https://cubari.moe/read/gist",
	regex:       ``,
	registrable: true,
}

var NHENTAI = Provider{
	name:        "nhentai",
	sourceURL:   "https://nhentai.net/g",
	regex:       `^((http(s)?://)?nhentai\.net/g/)?(\d{1,8})/?$`,
	registrable: false,
}

var MANGASEE = Provider{
	name:        "mangasee",
	sourceURL:   "https://mangasee123.com/manga",
	regex:       `^((http(s)?://)?mangasee123\.com/manga/)(.*?)/?$`,
	registrable: true,
}

var PROVIDERBYSTR = map[string]Provider{
	"imgur":    IMGUR,
	"gist":     GIST,
	"nhentai":  NHENTAI,
	"mangasee": MANGASEE,
}
