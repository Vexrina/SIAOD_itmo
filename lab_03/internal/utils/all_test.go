package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

var (
	index bleve.Index
)

type doc struct {
	id   int
	text string
	tags []string
}

// индексация рандомного текста
func Benchmark_IndexingThroughput(b *testing.B) {
	mapping := bleve.NewIndexMapping()
	tmpIndex, err := bleve.NewMemOnly(mapping)
	if err != nil {
		b.Fatal(err)
	}
	var metric []float64
	index = tmpIndex
	docs := generateDocs(b.N)

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		startIter := time.Now()
		document := docs[i]
		err := index.Index(fmt.Sprint(document.id), docs[i])
		if err != nil {
			b.Fatal(err)
		}
		metric = append(metric, float64(time.Now().Sub(startIter).Microseconds()))
	}

	elapsed := time.Since(start)
	docsPerSec := float64(b.N) / elapsed.Seconds()
	metrics := getMetric(metric)

	b.ReportMetric(docsPerSec, "docs/sec")
	b.ReportMetric(float64(elapsed.Microseconds())/float64(b.N), "microsec/op")
	b.ReportMetric(metrics.q25, "q25_microsec/doc")
	b.ReportMetric(metrics.median, "median_microsec/doc")
	b.ReportMetric(metrics.q75, "q75_microsec/doc")
	b.ReportMetric(metrics.q9, "q9_microsec/doc")
	b.ReportMetric(metrics.me, "mean_microsec/doc")
	b.ReportMetric(metrics.std, "std_microsec/doc")
}

func BenchmarkSearchLatency(b *testing.B) {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		err = idx.Index(fmt.Sprint(i), map[string]any{
			"text": generateRandomText(10),
			"tags": []string{fmt.Sprintf("tag%d", rand.Intn(100)), fmt.Sprintf("tag%d", rand.Intn(100))},
		})
	}
	err = idx.Index("100", map[string]any{
		"text": "the quick brown fox jumps over the lazy dog",
		"tags": []string{"tag1", "tag2"},
	})
	for i := 101; i <= 200; i++ {
		err = idx.Index(fmt.Sprint(i), map[string]any{
			"text": generateRandomText(10),
			"tags": []string{fmt.Sprintf("tag%d", rand.Intn(100)), fmt.Sprintf("tag%d", rand.Intn(100))},
		})
	}
	if err != nil {
		b.Fatal(err)
	}

	benchCases := []struct {
		name  string
		query query.Query
	}{
		{
			name:  "match",
			query: query.NewMatchQuery("quick"),
		},
		{
			name:  "term",
			query: query.NewTermQuery("tag1"),
		},
		{
			name:  "phrase",
			query: query.NewPhraseQuery([]string{"quick", "brown"}, "text"),
		},
	}
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var metrics []float64
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				start := time.Now()
				sr := bleve.NewSearchRequest(bc.query)

				_, serr := idx.Search(sr)

				latency := time.Since(start)

				if serr != nil {
					b.Fatal(serr)
				}
				metrics = append(metrics, float64(latency.Microseconds()))
			}
			m := getMetric(metrics)
			b.ReportMetric(m.q25, "q25_microsec/doc")
			b.ReportMetric(m.median, "median_microsec/doc")
			b.ReportMetric(m.q75, "q75_microsec/doc")
			b.ReportMetric(m.q9, "q9_microsec/doc")
			b.ReportMetric(m.me, "mean_microsec/doc")
			b.ReportMetric(m.std, "std_microsec/doc")
		})
	}
}

func Benchmark_SaveCSVIndex(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.N = 10
	cases := []struct {
		fileName  string
		benchName string
	}{
		{fileName: "../../../csv/example.csv", benchName: "example"},
		{fileName: "../../../csv/trc.csv", benchName: "trc"},
		{fileName: "../../../csv/tdvach.csv", benchName: "tdvach"},
	}
	for _, c := range cases {
		b.Run(c.benchName, func(b *testing.B) {
			var metrics []float64
			mapping := bleve.NewIndexMapping()
			idx, err := bleve.NewMemOnly(mapping)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				start := time.Now()
				err = ReadCsv(c.fileName, idx)
				if err != nil {
					b.Fatal(err)
				}
				end := time.Since(start)
				metrics = append(metrics, float64(end.Microseconds()))
			}
			m := getMetric(metrics)
			b.ReportMetric(m.q25, "q25_microsec/file")
			b.ReportMetric(m.median, "median_microsec/file")
			b.ReportMetric(m.q75, "q75_microsec/file")
			b.ReportMetric(m.q9, "q9_microsec/file")
			b.ReportMetric(m.me, "mean_microsec/file")
			b.ReportMetric(m.std, "std_microsec/file")
			mapping = nil
			runtime.GC()
		})
	}
}

func Benchmark_SearchIndex(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	cases := []struct {
		name     string
		filePath string
		terms    string
	}{
		{
			name:     "small file, 1 term",
			filePath: "../../../csv/example.csv",
			terms:    "барабан",
		},
		{
			name:     "small file, half of text",
			filePath: "../../../csv/example.csv",
			terms:    "Разработчики дайте щиты, ничего не понимаю,",
		},
		{
			name:     "small file, full text",
			filePath: "../../../csv/example.csv",
			terms:    "Разработчики дайте щиты, ничего не понимаю, последний барабан глючит постоянно",
		},
		{
			name:     "medium file, 1 term",
			filePath: "../../../csv/trc.csv",
			terms:    "колобки",
		},
		{
			name:     "medium file, half of text",
			filePath: "../../../csv/trc.csv",
			terms:    "ну что за хрень,уже на дальнем востоке народ знает в какой жопе страна а вы все телевизор,соловей и прочие киселевы",
		},
		{
			name:     "medium file, full text",
			filePath: "../../../csv/trc.csv",
			terms:    "пидр, россия единственная мире страна которая первая в мире произвела вакцину от эболы! ты, пидр не путай с окраиной, которая не тол ко вакцину от кори не способна изобрести, но и купить ее!",
		},
		{
			name:     "large file, 1 term",
			filePath: "../../../csv/tdvach.csv",
			terms:    "ватсаппом",
		},
		{
			name:     "large file, half of text",
			filePath: "../../../csv/tdvach.csv",
			terms:    "Если копипастите статью из интернета, хоть до конца ее дочитайте. С 17 года ",
		},
		{
			name:     "large file, full text",
			filePath: "../../../csv/tdvach.csv",
			terms:    "высшая школа экономики гайдар БРЕКИНГ НЬЮС - фашистская параша, обслуживающая интересы трансграничного олигархата и поклоняющаяся людям которые уничтожали наш народ в 90-ых - не любят авторитарного лидера который давил эту эксплуотаторскую гниль... действительно, почему народ так его любит? непонятно...",
		},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			var metrics []float64
			mapping := bleve.NewIndexMapping()
			idx, err := bleve.NewMemOnly(mapping)
			if err != nil {
				b.Fatal(err)
			}
			err = ReadCsv(c.filePath, idx)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				start := time.Now()
				_, serr := SearchIndex(idx, strings.Split(c.terms, " ")...)
				if serr != nil {
					b.Fatal(serr)
				}
				end := time.Since(start)
				metrics = append(metrics, float64(end.Microseconds()))
			}
			m := getMetric(metrics)
			b.ReportMetric(m.q25, "q25_microsec/query")
			b.ReportMetric(m.median, "median_microsec/query")
			b.ReportMetric(m.q75, "q75_microsec/query")
			b.ReportMetric(m.q9, "q9_microsec/query")
			b.ReportMetric(m.me, "mean_microsec/query")
			b.ReportMetric(m.std, "std_microsec/query")
		})
	}
}

func generateDocs(count int) []doc {
	docs := make([]doc, count)
	for i := 0; i < count; i++ {
		docs[i] = doc{
			id:   i,
			text: generateRandomText(100),
			tags: []string{"tag1", "tag2"},
		}
	}
	return docs
}

func generateRandomText(wordCount int) string {
	words := []string{"lorem", "ipsum", "dolor", "sit", "amet", "comicrosecectetur"}
	builder := strings.Builder{}

	for i := 0; i < wordCount; i++ {
		builder.WriteString(words[i%len(words)])
		builder.WriteString(" ")
	}

	return builder.String()
}

type Metric struct {
	q25    float64
	median float64
	q75    float64
	q9     float64
	me     float64
	std    float64
}

func getMetric(metric []float64) Metric {
	q25 := quantile(metric, 0.25)
	median := quantile(metric, 0.5)
	q75 := quantile(metric, 0.75)
	q9 := quantile(metric, 0.9)
	me := mean(metric)
	stdd := std(metric)

	return Metric{q25, median, q75, q9, me, stdd}
}

func quantile(data []float64, q float64) float64 {
	sort.Float64s(data)
	n := len(data)
	idx := q * float64(n-1)

	if idx == float64(int64(idx)) {
		return data[int(idx)]
	} else {

		lower := int(idx)
		upper := lower + 1
		weight := idx - float64(lower)
		return data[lower]*(1-weight) + data[upper]*weight
	}
}

func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func std(data []float64) float64 {
	m := mean(data)
	variance := 0.0
	for _, v := range data {
		variance += (v - m) * (v - m)
	}
	variance /= float64(len(data))
	return math.Sqrt(variance)
}
