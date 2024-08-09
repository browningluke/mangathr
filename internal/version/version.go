package version

var (
	version = "dev"
	sha     = "N/A"
)

func GetVersion() string {
	return version
}

func GetSHA() string {
	return sha
}
