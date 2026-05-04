package response

import (
	"testing"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

func TestToAttendanceResponse_FileURLs_NilSafeAndCorrect(t *testing.T) {
	date := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)

	t.Run("nil files", func(t *testing.T) {
		got := ToAttendanceResponse(model.Attendance{
			ID:     "a1",
			UserID: "u1",
			Status: model.AttendanceStatusPresent,
			Date:   date,
		})

		if got.CheckInFileURL != nil {
			t.Fatalf("CheckInFileURL = %v, want nil", *got.CheckInFileURL)
		}
		if got.CheckOutFileURL != nil {
			t.Fatalf("CheckOutFileURL = %v, want nil", *got.CheckOutFileURL)
		}
	})

	t.Run("maps check-in and check-out separately", func(t *testing.T) {
		checkInURL := "https://example.test/files/check_in.png"
		checkOutURL := "https://example.test/files/check_out.png"

		got := ToAttendanceResponse(model.Attendance{
			ID:     "a2",
			UserID: "u2",
			Status: model.AttendanceStatusPresent,
			Date:   date,
			CheckInFile: &model.File{
				FileURL: checkInURL,
			},
			CheckOutFile: &model.File{
				FileURL: checkOutURL,
			},
		})

		if got.CheckInFileURL == nil || *got.CheckInFileURL != checkInURL {
			if got.CheckInFileURL == nil {
				t.Fatalf("CheckInFileURL = nil, want %q", checkInURL)
			}
			t.Fatalf("CheckInFileURL = %q, want %q", *got.CheckInFileURL, checkInURL)
		}
		if got.CheckOutFileURL == nil || *got.CheckOutFileURL != checkOutURL {
			if got.CheckOutFileURL == nil {
				t.Fatalf("CheckOutFileURL = nil, want %q", checkOutURL)
			}
			t.Fatalf("CheckOutFileURL = %q, want %q", *got.CheckOutFileURL, checkOutURL)
		}
	})
}

