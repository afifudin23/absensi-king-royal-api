package request

type PayrollSettingRequest struct {
	ConfigName string  `json:"config_name" binding:"required,min=3,max=255"`
	Value      float64 `json:"value" binding:"required"`
}

type PayrollSettingByKeyRequest struct {
	ConfigKey string  `json:"config_key" binding:"required,min=3,max=255"`
	Value     float64 `json:"value" binding:"required"`
}
