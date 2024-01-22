package requests

type UpdateProductRequest struct {
	Id         string `json:"id" validate:"required"`
	Name       string `json:"name" validate:"required,min=3,max=12"`
	Category   string `json:"category" validate:"required, min=3, max=12"`
	OptionName string `json:"optionName" validate:"required, min=3, max= 12"`
}
