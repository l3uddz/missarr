package build

var (
	Version   string
	Timestamp string
	GitCommit string

	UserAgent string
)

func init() {
	UserAgent = "missarr/" + Version
}
