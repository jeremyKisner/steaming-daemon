package record

type Audio struct {
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	PickupURL   string `json:"pickup_url"`
	Description string `json:"description"`
	Plays       int    `json:"plays"`
}
