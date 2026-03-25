package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/magomedcoder/gen/internal/domain"
)

type mapFileRepo map[int64]*domain.File

func (m mapFileRepo) Create(context.Context, *domain.File) error {
	return nil
}

func (m mapFileRepo) UpdateStoragePath(context.Context, int64, string) error {
	return nil
}

func (m mapFileRepo) GetById(_ context.Context, id int64) (*domain.File, error) {
	f, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("файл в репозитории: запись с id=%d не найдена", id)
	}

	return f, nil
}

func TestHydrateImageAttachmentsForRunner_loadsFromDisk(t *testing.T) {
	dir := t.TempDir()
	sid := int64(7)
	sessDir := filepath.Join(dir, strconv.FormatInt(sid, 10))
	if err := os.MkdirAll(sessDir, 0o755); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(sessDir, "42_img.png")
	want := []byte{0x89, 0x50, 0x4e, 0x47, 0x01, 0x02}
	if err := os.WriteFile(path, want, 0o644); err != nil {
		t.Fatal(err)
	}

	fid := int64(42)
	c := &ChatUseCase{
		fileRepo: mapFileRepo{
			fid: {
				Id:          fid,
				Filename:    "img.png",
				StoragePath: path,
			},
		},
		attachmentsSaveDir: dir,
	}

	msgs := []*domain.Message{
		{
			SessionId:        sid,
			AttachmentName:   "img.png",
			AttachmentFileID: &fid,
			Content:          "caption",
		},
	}

	if err := c.hydrateImageAttachmentsForRunner(context.Background(), msgs); err != nil {
		t.Fatal(err)
	}

	if string(msgs[0].AttachmentContent) != string(want) {
		t.Fatalf("содержимое вложения не совпадает с файлом на диске")
	}
}

func TestHydrateImageAttachmentsForRunner_rejectsPathOutsideSession(t *testing.T) {
	dir := t.TempDir()
	sid := int64(7)
	other := filepath.Join(dir, "other", "evil.png")
	_ = os.MkdirAll(filepath.Dir(other), 0o755)
	_ = os.WriteFile(other, []byte{1, 2, 3}, 0o644)

	fid := int64(99)
	c := &ChatUseCase{
		fileRepo: mapFileRepo{
			fid: {
				Id:          fid,
				Filename:    "x.png",
				StoragePath: other,
			},
		},
		attachmentsSaveDir: dir,
	}
	msgs := []*domain.Message{
		{
			SessionId:        sid,
			AttachmentName:   "x.png",
			AttachmentFileID: &fid,
		},
	}
	if err := c.hydrateImageAttachmentsForRunner(context.Background(), msgs); err == nil {
		t.Fatal("ожидалась ошибка: путь к файлу вне каталога сессии")
	}
}
