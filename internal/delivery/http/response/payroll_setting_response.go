package response

import "github.com/afifudin23/absensi-king-royal-api/internal/model"

type PayrollSettingResponse struct {
	ID         string  `json:"id"`
	ConfigName string  `json:"config_name"`
	ConfigKey  string  `json:"config_key"`
	Value      float64 `json:"value"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type PayrollSettingDeleteResponse struct {
	Total        int `json:"total"`
	DeletedCount int `json:"deleted_count"`
	SkippedCount int `json:"skipped_count"`
}

func ToPayrollSettingResponse(payrollSetting model.PayrollSetting) PayrollSettingResponse {
	return PayrollSettingResponse{
		ID:         payrollSetting.ID,
		ConfigName: payrollSetting.ConfigName,
		ConfigKey:  payrollSetting.ConfigKey,
		Value:      payrollSetting.Value,
		CreatedAt:  payrollSetting.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  payrollSetting.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToPayrollSettingListResponse(payrollSettings []model.PayrollSetting) []PayrollSettingResponse {
	response := make([]PayrollSettingResponse, 0, len(payrollSettings))
	for _, payrollSetting := range payrollSettings {
		response = append(response, ToPayrollSettingResponse(payrollSetting))
	}
	return response
}

func ToPayrollSettingDeleteResponse(total, deletedCount, skippedCount int) PayrollSettingDeleteResponse {
	return PayrollSettingDeleteResponse{
		Total:        total,
		DeletedCount: deletedCount,
		SkippedCount: skippedCount,
	}
}
