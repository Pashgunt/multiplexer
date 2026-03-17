package dto

type TargetServiceDto struct {
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
	BaseUrl     string `json:"baseUrl"`
	IsActive    bool   `json:"isActive"`
}
