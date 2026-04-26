package wikiindex

type Retriever interface {
	Build(sourceDir string, chunks []InputChunk, opts ...BuildOptions) Stats

	Search(query string, topK int) []SearchResult

	Stats() Stats

	TokenLocale() string
}

type LexicalTFIDF struct {
	idx *Index
}

func NewLexicalTFIDF() *LexicalTFIDF {
	return &LexicalTFIDF{idx: New()}
}

func (l *LexicalTFIDF) Build(sourceDir string, chunks []InputChunk, opts ...BuildOptions) Stats {
	return l.idx.Build(sourceDir, chunks, opts...)
}

func (l *LexicalTFIDF) Search(query string, topK int) []SearchResult {
	return l.idx.Search(query, topK)
}

func (l *LexicalTFIDF) Stats() Stats {
	return l.idx.Stats()
}

func (l *LexicalTFIDF) TokenLocale() string {
	return l.idx.TokenLocale()
}

func (l *LexicalTFIDF) Index() *Index {
	return l.idx
}

type DenseEmbed struct {
	source string
	locale string
}

func NewDenseEmbed() *DenseEmbed {
	return &DenseEmbed{}
}

func (d *DenseEmbed) Build(sourceDir string, chunks []InputChunk, opts ...BuildOptions) Stats {
	d.source = sourceDir
	_ = chunks

	loc := ""
	if len(opts) > 0 {
		loc = normalizeLocale(opts[0].TokenLocale)
	}

	d.locale = loc

	return Stats{
		SourceDir:       sourceDir,
		Files:           0,
		Chunks:          0,
		UniqueTerms:     0,
		RetrievalMethod: "dense_embed_stub",
		TokenLocale:     loc,
		Stemming:        stemStatsLabel(loc),
	}
}

func (d *DenseEmbed) Search(_ string, _ int) []SearchResult {
	return nil
}

func (d *DenseEmbed) Stats() Stats {
	return Stats{
		SourceDir:       d.source,
		TokenLocale:     d.locale,
		Stemming:        stemStatsLabel(d.locale),
		RetrievalMethod: "dense_embed_stub",
	}
}

func (d *DenseEmbed) TokenLocale() string {
	return d.locale
}

var (
	_ Retriever = (*LexicalTFIDF)(nil)
	_ Retriever = (*DenseEmbed)(nil)
)
