package main

import (
	"fmt"

	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
)

func main() {

	data := map[string]string{
		"name": "Afiff",
		"time": "08:01",
	}

	err := utils.SendEmail(utils.EmailParams{
		FromName:  "Absensi System",
		FromEmail: "afifu5881@gmail.com",
		Password:  "hgdomzeuvlrmknbz",
		ToName:    "Afif",
		ToEmail:   "afifu5882@gmail.com",
		Subject:   "Notifikasi Absensi",
		Template:  "templates/email.html",
		Data:      data,
	})

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Email berhasil dikirim")
}
