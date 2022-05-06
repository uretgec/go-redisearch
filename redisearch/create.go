package redisearch

/*
FT.CREATE {index}
    [ON {data_type}]
       [PREFIX {count} {prefix} [{prefix} ..]
       [FILTER {filter}]
       [LANGUAGE {default_lang}]
       [LANGUAGE_FIELD {lang_attribute}]
       [SCORE {default_score}]
       [SCORE_FIELD {score_attribute}]
       [PAYLOAD_FIELD {payload_attribute}]
    [MAXTEXTFIELDS] [TEMPORARY {seconds}] [NOOFFSETS] [NOHL] [NOFIELDS] [NOFREQS] [SKIPINITIALSCAN]
    [STOPWORDS {num} {stopword} ...]
    SCHEMA {identifier} [AS {attribute}]
        [TEXT [NOSTEM] [WEIGHT {weight}] [PHONETIC {matcher}] | NUMERIC | GEO | TAG [SEPARATOR {sep}] [CASESENSITIVE]
        [SORTABLE [UNF]] [NOINDEX]] ...
*/

// ON {data_type}
const (
	HASH = "HASH"
	// JSON = "JSON" // Not ready to use
)

// Phonetic Matchers
const (
	PhoneticDoubleMetaphoneEnglish    string = "dm:en"
	PhoneticDoubleMetaphoneFrench     string = "dm:fr"
	PhoneticDoubleMetaphonePortuguese string = "dm:pt"
	PhoneticDoubleMetaphoneSpanish    string = "dm:es"
)

const (
	FieldTypeText    string = "TEXT"
	FieldTypeNumeric string = "NUMERIC"
	FieldTypeTag     string = "TAG"
	FieldTypeGeo     string = "GEO"
)

// TODO: json entegrasyonu
type FtSchema struct {
	identifier string // field name
	attribute  string // AS
	fieldtype  string // TEXT, NUMERIC, TAG, GEO
	sortable   bool
	option     FtSchemaOption
}

type FtSchemaOption struct {
	unf           bool
	nostem        bool
	noindex       bool
	phonetic      string
	weight        float32
	separator     string
	casesensitive bool
}

// Search Index Builder
type FtCreate struct {
	indexname        string
	datatype         string
	prefix           []string
	filterexp        string
	language         string
	languagefield    string
	score            float64
	scorefield       string
	payloadfiled     string
	maxtextfields    bool
	temporary        bool
	temporaryseconds int
	nooffsets        bool
	nofields         bool
	nofreqs          bool
	skipinitialscan  bool
	stopwords        []string
	schema           []FtSchema
}

func NewFtCreate(indexName string) *FtCreate {
	return &FtCreate{
		indexname: indexName,
	}
}

func (ftc *FtCreate) AddIndexName(name string) *FtCreate {
	ftc.indexname = name

	return ftc
}

func (ftc *FtCreate) AddDataType(dataType string) *FtCreate {
	ftc.datatype = dataType

	return ftc
}

func (ftc *FtCreate) AddPrefix(fields ...string) *FtCreate {
	ftc.prefix = append(ftc.prefix, fields...)

	return ftc
}

func (ftc *FtCreate) AddFilterExp(value string) *FtCreate {
	ftc.filterexp = value

	return ftc
}

// Arabic, Basque, Catalan, Danish, Dutch, English, Finnish, French, German, Greek, Hungarian, Indonesian, Irish, Italian,
// Lithuanian, Nepali, Norwegian, Portuguese, Romanian, Russian, Spanish, Swedish, Tamil, Turkish, Chinese
func (ftc *FtCreate) AddLanguage(lang string) *FtCreate {
	ftc.language = lang

	return ftc
}

func (ftc *FtCreate) AddLanguageField(field string) *FtCreate {
	ftc.languagefield = field

	return ftc
}

func (ftc *FtCreate) AddScore(score float64) *FtCreate {
	ftc.score = score

	return ftc
}

func (ftc *FtCreate) AddScoreField(field string) *FtCreate {
	ftc.scorefield = field

	return ftc
}

func (ftc *FtCreate) AddPayloadField(field string) *FtCreate {
	ftc.payloadfiled = field

	return ftc
}

func (ftc *FtCreate) AddMaxTextFields(active bool) *FtCreate {
	ftc.maxtextfields = active

	return ftc
}

func (ftc *FtCreate) AddTemporarySeconds(active bool, seconds int) *FtCreate {
	ftc.temporary = active
	ftc.temporaryseconds = seconds

	return ftc
}

func (ftc *FtCreate) AddNoOffsets(active bool) *FtCreate {
	ftc.nooffsets = active

	return ftc
}

func (ftc *FtCreate) AddNoFields(active bool) *FtCreate {
	ftc.nofields = active

	return ftc
}

func (ftc *FtCreate) AddNoFreqs(active bool) *FtCreate {
	ftc.nofreqs = active

	return ftc
}

func (ftc *FtCreate) AddSkipInitialScan(active bool) *FtCreate {
	ftc.skipinitialscan = active

	return ftc
}

func (ftc *FtCreate) AddStopWords(values ...string) *FtCreate {
	ftc.stopwords = append(ftc.stopwords, values...)

	return ftc
}

func (ftc *FtCreate) AddSchemaTextOption(weight float32, nostem, noindex bool, phonetic string) FtSchemaOption {
	return FtSchemaOption{
		unf:           false,
		nostem:        nostem,
		noindex:       noindex,
		phonetic:      phonetic,
		weight:        weight,
		separator:     "",
		casesensitive: false,
	}
}

func (ftc *FtCreate) AddSchemaTagOption(noindex bool, separator string) FtSchemaOption {
	return FtSchemaOption{
		unf:           false,
		nostem:        false,
		noindex:       noindex,
		phonetic:      "",
		weight:        0,
		separator:     separator,
		casesensitive: false,
	}
}

func (ftc *FtCreate) AddSchemaNumericOption(noindex bool) FtSchemaOption {
	return FtSchemaOption{
		unf:           false,
		nostem:        false,
		noindex:       noindex,
		phonetic:      "",
		weight:        0,
		separator:     "",
		casesensitive: false,
	}
}

func (ftc *FtCreate) AddSchemaGeoOption(noindex bool) FtSchemaOption {
	return FtSchemaOption{
		unf:           false,
		nostem:        false,
		noindex:       noindex,
		phonetic:      "",
		weight:        0,
		separator:     "",
		casesensitive: false,
	}
}

func (ftc *FtCreate) AddSchema(fieldType string, identifier string, attr string, sortable bool, option FtSchemaOption) *FtCreate {
	ftc.schema = append(ftc.schema, FtSchema{
		identifier: identifier,
		attribute:  attr,
		fieldtype:  fieldType,
		sortable:   sortable,
		option:     option,
	})

	return ftc
}

func (ftc *FtCreate) Serialize() []interface{} {

	var queryCode []interface{}

	queryCode = append(queryCode, "FT.CREATE")

	queryCode = append(queryCode, ftc.indexname)

	queryCode = append(queryCode, "ON", ftc.datatype)

	if ftc.prefix != nil && len(ftc.prefix) > 0 {
		queryCode = append(queryCode, "PREFIX", len(ftc.prefix))

		for _, p := range ftc.prefix {
			queryCode = append(queryCode, p)
		}
	}

	if ftc.filterexp != "" {
		queryCode = append(queryCode, "FILTER", ftc.filterexp)
	}

	if ftc.language != "" {
		queryCode = append(queryCode, "LANGUAGE", ftc.language)
	}

	if ftc.languagefield != "" {
		queryCode = append(queryCode, "LANGUAGE_FIELD", ftc.languagefield)
	}

	if ftc.score > 0 {
		queryCode = append(queryCode, "SCORE", ftc.language)
	}

	if ftc.scorefield != "" {
		queryCode = append(queryCode, "SCORE_FIELD", ftc.scorefield)
	}

	if ftc.payloadfiled != "" {
		queryCode = append(queryCode, "PAYLOAD_FIELD", ftc.payloadfiled)
	}

	if ftc.maxtextfields {
		queryCode = append(queryCode, "MAXTEXTFIELDS")
	}

	if ftc.temporary && ftc.temporaryseconds > 0 {
		queryCode = append(queryCode, "TEMPORARY", ftc.temporaryseconds)
	}

	if ftc.nooffsets {
		queryCode = append(queryCode, "NOOFFSETS")
	}

	if ftc.nofields {
		queryCode = append(queryCode, "NOFIELDS")
	}

	if ftc.nofreqs {
		queryCode = append(queryCode, "NOFREQS")
	}

	if ftc.skipinitialscan {
		queryCode = append(queryCode, "SKIPINITIALSCAN")
	}

	if ftc.stopwords != nil && len(ftc.stopwords) > 0 {
		queryCode = append(queryCode, "STOPWORDS", len(ftc.stopwords))

		for _, sw := range ftc.stopwords {
			queryCode = append(queryCode, sw)
		}
	}

	if ftc.schema != nil && len(ftc.schema) > 0 {
		queryCode = append(queryCode, "SCHEMA")

		for _, sc := range ftc.schema {

			queryCode = append(queryCode, sc.identifier)

			switch sc.fieldtype {
			case FieldTypeText:

				if sc.attribute != "" {
					queryCode = append(queryCode, "AS", sc.attribute)
				}

				queryCode = append(queryCode, "TEXT")

				if sc.option.nostem {
					queryCode = append(queryCode, "NOSTEM")
				}

				if sc.sortable {
					queryCode = append(queryCode, "SORTABLE")

					if sc.option.unf {
						queryCode = append(queryCode, "UNF")
					}
				}

				if sc.option.weight != 0 && sc.option.weight != 1 {
					queryCode = append(queryCode, "WEIGHT", sc.option.weight)
				}

				if sc.option.phonetic != "" {
					queryCode = append(queryCode, "PHONETIC", sc.option.phonetic)
				}

			case FieldTypeNumeric:

				if sc.attribute != "" {
					queryCode = append(queryCode, "AS", sc.attribute)
				}

				queryCode = append(queryCode, "NUMERIC")

				if sc.sortable {
					queryCode = append(queryCode, "SORTABLE")

					if sc.option.unf {
						queryCode = append(queryCode, "UNF")
					}
				}

			case FieldTypeTag:

				if sc.attribute != "" {
					queryCode = append(queryCode, "AS", sc.attribute)
				}

				queryCode = append(queryCode, "TAG")

				if sc.option.separator != "" {
					queryCode = append(queryCode, "SEPARATOR", sc.option.separator)
				}

				if sc.sortable {
					queryCode = append(queryCode, "SORTABLE")

					if sc.option.unf {
						queryCode = append(queryCode, "UNF")
					}
				}

			case FieldTypeGeo:

				if sc.attribute != "" {
					queryCode = append(queryCode, "AS", sc.attribute)
				}

				queryCode = append(queryCode, "GEO")

			}

			if sc.option.noindex {
				queryCode = append(queryCode, "NOINDEX")
			}

		}
	}

	return queryCode
}
