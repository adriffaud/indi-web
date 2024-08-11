package indiclient

type PropertySelector struct {
	Device string `json:"device"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}
