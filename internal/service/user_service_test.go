package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	getAllFn     func(ctx context.Context) ([]model.User, error)
	getByIDFn    func(ctx context.Context, id string) (*model.User, error)
	getByEmailFn func(ctx context.Context, email string) (*model.User, error)
	createFn     func(ctx context.Context, user *model.User, profile *model.UserProfile) error
	updateFn     func(ctx context.Context, user *model.User, profile *model.UserProfile) error
	deleteFn     func(ctx context.Context, id string) error

	getAllCalls  int
	getByIDCalls int
	createCalls  int
	updateCalls  int
	deleteCalls  int

	lastCreatedUser    *model.User
	lastCreatedProfile *model.UserProfile
	lastUpdatedUser    *model.User
	lastUpdatedProfile *model.UserProfile
	lastDeleted        string
}

func (m *mockUserRepo) GetAll(ctx context.Context, loadProfile bool) ([]model.User, error) {
	m.getAllCalls++
	return m.getAllFn(ctx)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string, loadProfile bool) (*model.User, error) {
	m.getByIDCalls++
	return m.getByIDFn(ctx, id)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.getByEmailFn == nil {
		return nil, errors.New("GetByEmail not implemented in mock")
	}
	return m.getByEmailFn(ctx, email)
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	m.createCalls++
	m.lastCreatedUser = user
	m.lastCreatedProfile = profile
	return m.createFn(ctx, user, profile)
}

func (m *mockUserRepo) Update(ctx context.Context, user *model.User, profile *model.UserProfile) error {
	m.updateCalls++
	m.lastUpdatedUser = user
	m.lastUpdatedProfile = profile
	return m.updateFn(ctx, user, profile)
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	m.deleteCalls++
	m.lastDeleted = id
	return m.deleteFn(ctx, id)
}

func strPtr(v string) *string { return &v }

func TestUserService_Create_HashesPasswordAndPersists(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(ctx context.Context, user *model.User, profile *model.UserProfile) error {
			return nil
		},
	}
	svc := &userService{userRepo: repo}

	in := request.UserCreateRequest{
		FullName: "Jane Doe",
		Email:    "jane@example.com",
		Password: "supersecret123",
		Role:     model.UserRoleUser,
	}

	user, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if user.ID == "" {
		t.Fatalf("Create() user.ID is empty")
	}
	if repo.createCalls != 1 {
		t.Fatalf("repo.Create calls = %d, want %d", repo.createCalls, 1)
	}
	if repo.lastCreatedUser == nil {
		t.Fatalf("repo.lastCreatedUser is nil")
	}
	if repo.lastCreatedUser.Password == "" || repo.lastCreatedUser.Password == in.Password {
		t.Fatalf("password was not hashed")
	}
	if !utils.CheckPassword(in.Password, repo.lastCreatedUser.Password) {
		t.Fatalf("stored password hash does not match input password")
	}
	if repo.lastCreatedProfile == nil {
		t.Fatalf("repo.lastCreatedProfile is nil")
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(ctx context.Context, user *model.User, profile *model.UserProfile) error {
			return errors.New("Error 1062 (23000): Duplicate entry x for key uq_users_email")
		},
	}
	svc := &userService{userRepo: repo}

	_, err := svc.Create(context.Background(), request.UserCreateRequest{
		FullName: "Jane Doe",
		Email:    "jane@example.com",
		Password: "supersecret123",
		Role:     model.UserRoleUser,
	})
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf("Create() err = %v, want ErrEmailAlreadyRegistered", err)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := &userService{userRepo: repo}

	_, err := svc.GetByID(context.Background(), "missing")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetByID() err = %v, want ErrUserNotFound", err)
	}
}

func TestUserService_Update_AppliesFieldsAndPersists(t *testing.T) {
	oldName := "Old Name"
	newName := "New Name"
	roleAdmin := model.UserRoleAdmin
	phone := "08123456789"

	existing := &model.User{
		ID:       "u1",
		FullName: oldName,
		Email:    "old@example.com",
		Role:     model.UserRoleUser,
		Profile:  &model.UserProfile{UserID: "u1"},
	}

	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) { return existing, nil },
		updateFn:  func(ctx context.Context, user *model.User, profile *model.UserProfile) error { return nil },
	}
	svc := &userService{userRepo: repo}

	updated, err := svc.Update(context.Background(), "u1", request.UserUpdateRequest{
		FullName:    &newName,
		Role:        &roleAdmin,
		PhoneNumber: &phone,
	})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.FullName != newName {
		t.Fatalf("Update() FullName = %q, want %q", updated.FullName, newName)
	}
	if updated.Role != model.UserRoleAdmin {
		t.Fatalf("Update() Role = %q, want %q", updated.Role, model.UserRoleAdmin)
	}
	if updated.Profile == nil || updated.Profile.PhoneNumber == nil || *updated.Profile.PhoneNumber != phone {
		t.Fatalf("Update() Profile.PhoneNumber = %v, want %q", updated.Profile, phone)
	}
	if repo.updateCalls != 1 {
		t.Fatalf("repo.Update calls = %d, want %d", repo.updateCalls, 1)
	}
}

func TestUserService_Update_DuplicateEmail(t *testing.T) {
	existing := &model.User{ID: "u1", FullName: "A", Email: "a@a.com", Role: model.UserRoleUser, Profile: &model.UserProfile{UserID: "u1"}}
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) { return existing, nil },
		updateFn: func(ctx context.Context, user *model.User, profile *model.UserProfile) error {
			return errors.New("duplicate entry")
		},
	}
	svc := &userService{userRepo: repo}

	_, err := svc.Update(context.Background(), "u1", request.UserUpdateRequest{FullName: strPtr("B")})
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Fatalf("Update() err = %v, want ErrEmailAlreadyRegistered", err)
	}
}

func TestUserService_UpdateProfile_HashesPassword(t *testing.T) {
	existing := &model.User{ID: "u1", FullName: "A", Email: "a@a.com", Role: model.UserRoleUser, Password: "oldhash", Profile: &model.UserProfile{UserID: "u1"}}
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) { return existing, nil },
		updateFn:  func(ctx context.Context, user *model.User, profile *model.UserProfile) error { return nil },
	}
	svc := &userService{userRepo: repo}

	newPassword := "newpass12345"
	newEmail := "new@example.com"
	updated, err := svc.UpdateProfile(context.Background(), "u1", request.UserUpdateProfileRequest{
		Email:    &newEmail,
		Password: &newPassword,
	})
	if err != nil {
		t.Fatalf("UpdateProfile() err = %v", err)
	}
	if updated.Email != newEmail {
		t.Fatalf("UpdateProfile() Email = %q, want %q", updated.Email, newEmail)
	}
	if updated.Password == "" || updated.Password == newPassword {
		t.Fatalf("UpdateProfile() password was not hashed")
	}
	if !utils.CheckPassword(newPassword, updated.Password) {
		t.Fatalf("stored password hash does not match new password")
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		deleteFn:  func(ctx context.Context, id string) error { return nil },
	}
	svc := &userService{userRepo: repo}

	err := svc.Delete(context.Background(), "missing")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("Delete() err = %v, want ErrUserNotFound", err)
	}
}

func TestUserService_Delete_DeletesExisting(t *testing.T) {
	existing := &model.User{ID: "u1", FullName: "A", Email: "a@a.com", Role: model.UserRoleUser, Profile: &model.UserProfile{UserID: "u1"}}
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) { return existing, nil },
		deleteFn:  func(ctx context.Context, id string) error { return nil },
	}
	svc := &userService{userRepo: repo}

	if err := svc.Delete(context.Background(), "u1"); err != nil {
		t.Fatalf("Delete() err = %v", err)
	}
	if repo.deleteCalls != 1 || repo.lastDeleted != "u1" {
		t.Fatalf("repo.Delete calls=%d id=%q, want calls=1 id=%q", repo.deleteCalls, repo.lastDeleted, "u1")
	}
}

func TestUserService_GetAll(t *testing.T) {
	repo := &mockUserRepo{
		getAllFn: func(ctx context.Context) ([]model.User, error) {
			return []model.User{{ID: "u1"}, {ID: "u2"}}, nil
		},
	}
	svc := &userService{userRepo: repo}

	users, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() err = %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("GetAll() len = %d, want %d", len(users), 2)
	}
	if repo.getAllCalls != 1 {
		t.Fatalf("repo.GetAll calls = %d, want %d", repo.getAllCalls, 1)
	}
}

func TestUserService_Create_PassesOptionalFields(t *testing.T) {
	employeeCode := "EMP001"
	basicSalary := 1234.56
	birthDate := time.Date(1995, 1, 2, 0, 0, 0, 0, time.UTC)

	repo := &mockUserRepo{
		createFn: func(ctx context.Context, user *model.User, profile *model.UserProfile) error {
			if profile == nil {
				return errors.New("profile not passed")
			}
			if profile.EmployeeCode == nil || *profile.EmployeeCode != employeeCode {
				return errors.New("employee_code not passed")
			}
			if profile.BasicSalary == nil || *profile.BasicSalary != basicSalary {
				return errors.New("basic_salary not passed")
			}
			if profile.BirthDate == nil || !profile.BirthDate.Equal(birthDate) {
				return errors.New("birth_date not passed")
			}
			return nil
		},
	}
	svc := &userService{userRepo: repo}

	_, err := svc.Create(context.Background(), request.UserCreateRequest{
		FullName:     "A",
		Email:        "a@a.com",
		Password:     "password123",
		Role:         model.UserRoleUser,
		EmployeeCode: &employeeCode,
		BasicSalary:  &basicSalary,
		BirthDate:    &birthDate,
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
}
