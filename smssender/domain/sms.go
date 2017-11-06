package domain

// CHUCKSIZE is defined to 160 - udh length
const CHUCKSIZE = 153

// SMS is our representation of an SMS message
type SMS struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
	cur        int
}

// Split splits the message in chunks of CHUCKSIZE size
func (s *SMS) Split() []string {

	slices := make([]string, 0)
	msglen := len(s.Message)
	for i := 0; i < msglen; i += CHUCKSIZE {

		ulimit := msglen
		if CHUCKSIZE+i <= ulimit {
			ulimit = CHUCKSIZE + i
		}

		slices = append(slices, string(s.Message[i:ulimit]))
	}

	return slices
}
