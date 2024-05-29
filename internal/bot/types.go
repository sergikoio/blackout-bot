package bot

type status string

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
