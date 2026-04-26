package wikiindex

import "testing"

func TestBuildAndSearch(t *testing.T) {
	idx := New()
	stats := idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "guide/auth.md",
			FileName:   "auth.md",
			ChunkIndex: 0,
			Text:       "JWT token проверяется middleware перед доступом к API.",
		},
		{
			FilePath:   "guide/billing.md",
			FileName:   "billing.md",
			ChunkIndex: 0,
			Text:       "Платежи обрабатываются асинхронной очередью и retry worker.",
		},
	})

	if stats.Files != 2 {
		t.Fatalf("expected 2 files, got %d", stats.Files)
	}
	if stats.Chunks != 2 {
		t.Fatalf("expected 2 chunks, got %d", stats.Chunks)
	}

	results := idx.Search("Как проверяется JWT токен?", 3)
	if len(results) == 0 {
		t.Fatal("expected search results")
	}

	if results[0].Chunk.FilePath != "guide/auth.md" {
		t.Fatalf("expected auth chunk first, got %s", results[0].Chunk.FilePath)
	}
}

func TestBuildWithTokenLocaleRU(t *testing.T) {
	idx := New()
	stats := idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "a.md",
			FileName:   "a.md",
			ChunkIndex: 0, Text: "авторизация и проверка токена",
		},
	}, BuildOptions{
		TokenLocale: "ru",
	})
	if stats.TokenLocale != "ru" {
		t.Fatalf("stats token_locale: %q", stats.TokenLocale)
	}

	if stats.Stemming != "snowball_russian" {
		t.Fatalf("stats stemming: %q", stats.Stemming)
	}

	res := idx.Search("и", 3)
	if len(res) != 0 {
		t.Fatalf("expected no hit on stopword-only query, got %d", len(res))
	}

	hits := idx.Search("авторизация", 3)
	if len(hits) == 0 {
		t.Fatal("expected hit on content term")
	}
}

func TestStemRussianInflection(t *testing.T) {
	idx := New()
	idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "sec.md",
			FileName:   "sec.md",
			ChunkIndex: 0,
			Text:       "Проверки доступа к системе выполняются ежедневно.",
		},
	}, BuildOptions{TokenLocale: "ru"})

	res := idx.Search("проверка доступ система", 3)
	if len(res) == 0 {
		t.Fatal("expected stemmed match for RU inflections")
	}

	if res[0].Chunk.FilePath != "sec.md" {
		t.Fatalf("unexpected hit: %s", res[0].Chunk.FilePath)
	}
}

func TestStemRussianOffWithoutLocale(t *testing.T) {
	idx := New()
	idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "sec.md",
			FileName:   "sec.md",
			ChunkIndex: 0,
			Text:       "Проверки доступа к системе.",
		},
	})

	res := idx.Search("проверка доступ система", 3)
	if len(res) != 0 {
		t.Fatalf("without token_locale expected no stem match, got %d hits", len(res))
	}
}

func TestStemEnglishInflection(t *testing.T) {
	idx := New()
	idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "run.md",
			FileName:   "run.md",
			ChunkIndex: 0,
			Text:       "The service is running in production.",
		},
	}, BuildOptions{
		TokenLocale: "en",
	})

	res := idx.Search("runs production", 3)
	if len(res) == 0 {
		t.Fatal("expected stemmed match: running/run")
	}
}

func TestSearchNoMatch(t *testing.T) {
	idx := New()
	idx.Build("/tmp/wiki", []InputChunk{
		{
			FilePath:   "a.txt",
			FileName:   "a.txt",
			ChunkIndex: 0,
			Text:       "Документ о логировании и мониторинге.",
		},
	})

	results := idx.Search("квантовая физика", 5)
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}
