package seeder

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/google/uuid"
)

func SeedPayrollSettings() {
	db := config.GetDB()

	payrollSettings := []model.PayrollSetting{
		{ID: uuid.NewString(), ConfigName: "Hourly Overtime Rate", ConfigKey: "hourly_overtime_rate", Value: 12000},
		{ID: uuid.NewString(), ConfigName: "BPJS Health Rate", ConfigKey: "bpjs_health_rate", Value: 100000},
		{ID: uuid.NewString(), ConfigName: "BPJS Employment JHT Rate", ConfigKey: "bpjs_employment_jht_rate", Value: 150000},
		{ID: uuid.NewString(), ConfigName: "BPJS Employment JP Rate", ConfigKey: "bpjs_employment_jp_rate", Value: 50000},
	}

	for _, setting := range payrollSettings {
		var existing model.PayrollSetting
		err := db.Where("config_name = ?", setting.ConfigName).Attrs(setting).FirstOrCreate(&existing).Error
		if err != nil {
			continue
		}
	}

}
