package dto

type TargetServiceDto struct {
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
	BaseURL     string `json:"baseUrl"`
	IsActive    bool   `json:"isActive"`
}
