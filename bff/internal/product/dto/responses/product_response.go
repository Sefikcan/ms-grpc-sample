package responses

type ProductResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	OptionName string `json:"optionName"`
}
