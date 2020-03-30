package globals

var (
	Conf       Config
	Hostname   string
	Conn_cache = make(map[string]int64)
)

type Config struct {
	Ports       []int
	SMTP_Server string
	SMTP_User   string
	SMTP_Passwd string
	Mail_From   string
	Mail_To     []string
	Timeout     int64
}
