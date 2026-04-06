package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type mockFileRepo struct {
	createCalls int
	lastCreated *model.File
	createFn    func(ctx context.Context, file *model.File) error
}

func (m *mockFileRepo) Create(ctx context.Context, file *model.File) error {
	m.createCalls++
	m.lastCreated = file
	if m.createFn != nil {
		return m.createFn(ctx, file)
	}
	return nil
}

func (m *mockFileRepo) GetByID(ctx context.Context, id string) (*model.File, error) { return nil, nil }
func (m *mockFileRepo) Delete(ctx context.Context, id string) error                 { return nil }

func TestFileService_Upload_RejectsNonImage(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() err = %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("os.Chdir() err = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	fileHeader := newMultipartFileHeader(t, "file.txt", []byte("hello world"))

	repo := &mockFileRepo{}
	svc := NewFileService(repo, "http://localhost:8080")

	_, uploadErr := svc.Upload(context.Background(), fileHeader, model.FileTypeProfilePicture, "user-1")
	appErr, ok := uploadErr.(*common.AppError)
	if !ok {
		t.Fatalf("err type = %T, want *common.AppError (err=%v)", uploadErr, uploadErr)
	}
	if appErr.StatusCode != 400 {
		t.Fatalf("status = %d, want %d", appErr.StatusCode, 400)
	}
	if repo.createCalls != 0 {
		t.Fatalf("repo.Create calls = %d, want %d", repo.createCalls, 0)
	}
}

func TestFileService_Upload_AcceptsPNG(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() err = %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("os.Chdir() err = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	// 1x1 transparent PNG
	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xe5, 0x27, 0xd4, 0xa2, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
		0x42, 0x60, 0x82,
	}

	fileHeader := newMultipartFileHeader(t, "image.png", pngBytes)

	repo := &mockFileRepo{}
	baseURL := "https://example.test"
	svc := NewFileService(repo, baseURL)

	got, err := svc.Upload(context.Background(), fileHeader, model.FileTypeProfilePicture, "user-1")
	if err != nil {
		t.Fatalf("Upload() err = %v", err)
	}
	if repo.createCalls != 1 {
		t.Fatalf("repo.Create calls = %d, want %d", repo.createCalls, 1)
	}
	if got.MimeType != "image/png" {
		t.Fatalf("MimeType = %q, want %q", got.MimeType, "image/png")
	}
	if !strings.HasSuffix(got.FileName, ".png") {
		t.Fatalf("FileName = %q, want suffix %q", got.FileName, ".png")
	}
	if !strings.HasPrefix(got.FileURL, baseURL+"/files/") {
		t.Fatalf("FileURL = %q, want prefix %q", got.FileURL, baseURL+"/files/")
	}
	if _, err := os.Stat(filepath.Clean(got.FilePath)); err != nil {
		t.Fatalf("os.Stat(FilePath) err = %v", err)
	}

	// Sanity check folder naming includes today's date.
	dateFolder := time.Now().Format("2006-01-02")
	if !strings.Contains(got.FilePath, filepath.Join("files", "profile_picture", dateFolder)) {
		t.Fatalf("FilePath = %q, want contain %q", got.FilePath, filepath.Join("files", "profile_picture", dateFolder))
	}
}

func newMultipartFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile() err = %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("part.Write() err = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() err = %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		t.Fatalf("ParseMultipartForm() err = %v", err)
	}

	_, fh, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("FormFile() err = %v", err)
	}
	return fh
}

