package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type mockAttendanceRequestRepo struct {
	getByIDFn func(ctx context.Context, id string, loadFile bool) (*model.AttendanceRequest, error)
	updateFn  func(ctx context.Context, req *model.AttendanceRequest) error
}

func (m *mockAttendanceRequestRepo) Create(ctx context.Context, req *model.AttendanceRequest) error {
	return errors.New("not implemented")
}
func (m *mockAttendanceRequestRepo) GetAll(ctx context.Context, loadFile bool) ([]model.AttendanceRequest, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAttendanceRequestRepo) GetByID(ctx context.Context, id string, loadFile bool) (*model.AttendanceRequest, error) {
	return m.getByIDFn(ctx, id, loadFile)
}
func (m *mockAttendanceRequestRepo) GetByUserID(ctx context.Context, userID string) ([]model.AttendanceRequest, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAttendanceRequestRepo) Update(ctx context.Context, req *model.AttendanceRequest) error {
	return m.updateFn(ctx, req)
}
func (m *mockAttendanceRequestRepo) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

type mockAttendanceRepo struct {
	byUserDate map[string]*model.Attendance
	createCalls int
	updateCalls int
}

func newMockAttendanceRepo() *mockAttendanceRepo {
	return &mockAttendanceRepo{byUserDate: map[string]*model.Attendance{}}
}

func (m *mockAttendanceRepo) key(userID string, date time.Time) string {
	return userID + "|" + date.Format("2006-01-02")
}

func (m *mockAttendanceRepo) GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*model.Attendance, error) {
	if v, ok := m.byUserDate[m.key(userID, date)]; ok {
		clone := *v
		return &clone, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAttendanceRepo) Create(ctx context.Context, attendance *model.Attendance) error {
	m.createCalls++
	m.byUserDate[m.key(attendance.UserID, attendance.Date)] = attendance
	return nil
}

func (m *mockAttendanceRepo) Update(ctx context.Context, attendance *model.Attendance) error {
	m.updateCalls++
	m.byUserDate[m.key(attendance.UserID, attendance.Date)] = attendance
	return nil
}

func (m *mockAttendanceRepo) GetLogsByUserID(ctx context.Context, userID string) ([]model.Attendance, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAttendanceRepo) GetByID(ctx context.Context, id string) (*model.Attendance, error) {
	return nil, gorm.ErrRecordNotFound
}

type noopFileRepo struct{}

func (noopFileRepo) Create(ctx context.Context, image *model.File) error { return nil }
func (noopFileRepo) GetByID(ctx context.Context, id string) (*model.File, error) {
	return nil, gorm.ErrRecordNotFound
}
func (noopFileRepo) Delete(ctx context.Context, id string) error { return nil }

func TestAttendanceRequestService_UpdateStatus_Approved_UpsertsAttendance(t *testing.T) {
	day := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)

	reqRepo := &mockAttendanceRequestRepo{
		getByIDFn: func(ctx context.Context, id string, loadFile bool) (*model.AttendanceRequest, error) {
			return &model.AttendanceRequest{
				ID:        id,
				UserID:    "u1",
				Type:      model.AttendanceRequestTypeSick,
				Status:    model.AttendanceRequestStatusPending,
				StartDate: day,
				EndDate:   day,
				Reason:    "flu",
			}, nil
		},
		updateFn: func(ctx context.Context, req *model.AttendanceRequest) error {
			return nil
		},
	}

	attRepo := newMockAttendanceRepo()
	svc := NewAttendanceRequestService(reqRepo, attRepo, noopFileRepo{}).(*attendanceRequestService)

	_, err := svc.UpdateStatus(context.Background(), "admin-1", "req-1", request.AttendanceRequestUpdateStatusRequest{
		Status: model.AttendanceRequestStatusApproved,
	})
	if err != nil {
		t.Fatalf("UpdateStatus() err = %v", err)
	}

	finalDay := startOfDay(day)
	got, err := attRepo.GetByUserAndDate(context.Background(), "u1", finalDay)
	if err != nil {
		t.Fatalf("GetByUserAndDate() err = %v", err)
	}
	if got.Status != model.AttendanceStatusSick {
		t.Fatalf("Status = %q, want %q", got.Status, model.AttendanceStatusSick)
	}
	if got.Source != model.AttendanceSourceApprovedRequest {
		t.Fatalf("Source = %q, want %q", got.Source, model.AttendanceSourceApprovedRequest)
	}
	if got.UpdatedBy == nil || *got.UpdatedBy != "admin-1" {
		if got.UpdatedBy == nil {
			t.Fatalf("UpdatedBy = nil, want %q", "admin-1")
		}
		t.Fatalf("UpdatedBy = %q, want %q", *got.UpdatedBy, "admin-1")
	}
	if got.Note == nil || *got.Note != "flu" {
		if got.Note == nil {
			t.Fatalf("Note = nil, want %q", "flu")
		}
		t.Fatalf("Note = %q, want %q", *got.Note, "flu")
	}
	if got.CheckInAt != nil || got.CheckOutAt != nil {
		t.Fatalf("CheckInAt/CheckOutAt should be nil for sick request")
	}
}

