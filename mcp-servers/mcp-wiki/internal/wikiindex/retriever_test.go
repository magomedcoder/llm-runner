package wikiindex

import "testing"

func TestLexicalTFIDFRetriever(t *testing.T) {
	r := NewLexicalTFIDF()
	stats := r.Build("/w", []InputChunk{
		{
			FilePath:   "a.md",
			FileName:   "a.md",
			ChunkIndex: 0,
			Text:       "alpha beta uniqueword",
		},
	})

	if stats.Chunks != 1 {
		t.Fatalf("chunks: %d", stats.Chunks)
	}

	if len(r.Search("uniqueword", 3)) != 1 {
		t.Fatal("expected search hit")
	}
}

func TestDenseEmbedStub(t *testing.T) {
	d := NewDenseEmbed()
	stats := d.Build("/w", []InputChunk{
		{
			FilePath:   "a.md",
			FileName:   "a.md",
			ChunkIndex: 0, Text: "ignored until embed pipeline",
		},
	}, BuildOptions{
		TokenLocale: "en",
	})
	if stats.RetrievalMethod != "dense_embed_stub" || stats.Chunks != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	if d.TokenLocale() != "en" {
		t.Fatalf("token locale: %q", d.TokenLocale())
	}

	if len(d.Search("anything", 5)) != 0 {
		t.Fatal("stub search must be empty")
	}
}
