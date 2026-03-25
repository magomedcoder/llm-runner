package service

import (
	"fmt"
	"strings"

	"github.com/magomedcoder/llm-runner/domain"
)

func errIfVisionAttachments(messages []*domain.AIChatMessage, mmprojPath, canonicalGGUF, modelsDir string) error {
	disp := DisplayModelName(canonicalGGUF)
	for _, m := range messages {
		if m == nil || len(m.AttachmentContent) == 0 {
			continue
		}

		name := m.AttachmentName
		if name == "" {
			name = "attachment"
		}

		if strings.TrimSpace(mmprojPath) == "" {
			return fmt.Errorf("llama: изображение %q (%d байт): не найден mmproj для модели %q - положите в %s файл %s-mmproj.gguf (или %s.mmproj.gguf) либо задайте mmproj_path в config.yaml", name, len(m.AttachmentContent), disp, modelsDir, disp, disp)
		}
	}

	return nil
}
