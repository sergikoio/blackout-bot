package bot

type status string

var minutesForms = []string{"хвилину", "хвилини", "хвилин"}
var hoursForms = []string{"годину", "години", "годин"}

const (
	offlineStatus status = "offline"
	onlineStatus  status = "online"
)

func (s status) Validate() bool {
	if s == offlineStatus || s == onlineStatus {
		return true
	}

	return false
}

func (s status) ToString() string {
	return string(s)
}

type ServerConfig struct {
	IsEmergency bool `json:"is_emergency"`
	IsBotOff    bool `json:"is_send_reject"`
}
