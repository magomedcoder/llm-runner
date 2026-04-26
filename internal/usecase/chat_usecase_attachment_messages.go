package usecase

import (
	"fmt"
	"strings"
)

const hydratedAttachmentExcerptRunes = 320

func buildCompactHydratedAttachmentMessage(attachmentName, extractedText, userMessage string, includeUserMessage bool) string {
	extracted := strings.TrimSpace(extractedText)
	excerpt := truncateStringRunes(extracted, hydratedAttachmentExcerptRunes)

	var b strings.Builder
	fmt.Fprintf(&b, "[attachment_ref: %s]\n", attachmentName)
	if excerpt != "" {
		b.WriteString("Краткое содержание вложения:\n")
		b.WriteString(excerpt)
		b.WriteString("\n")
	}

	userMessage = strings.TrimSpace(userMessage)
	if includeUserMessage && userMessage != "" {
		b.WriteString("\nСообщение пользователя:\n")
		b.WriteString(userMessage)
	}

	return strings.TrimSpace(b.String())
}
