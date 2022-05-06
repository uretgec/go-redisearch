// Via: https://oss.redis.com/redisearch/Commands/#ftsearch
package redisearch

import (
	"fmt"
	"math"
)

/*
FT.SEARCH {index} {query} [NOCONTENT] [VERBATIM] [NOSTOPWORDS] [WITHSCORES] [WITHPAYLOADS] [WITHSORTKEYS]
  [FILTER {numeric_attribute} {min} {max}] ...
  [GEOFILTER {geo_attribute} {lon} {lat} {radius} m|km|mi|ft]
  [INKEYS {num} {key} ... ]
  [INFIELDS {num} {attribute} ... ]
  [RETURN {num} {identifier} [AS {property}] ... ]
  [SUMMARIZE [FIELDS {num} {attribute} ... ] [FRAGS {num}] [LEN {fragsize}] [SEPARATOR {separator}]]
  [HIGHLIGHT [FIELDS {num} {attribute} ... ] [TAGS {open} {close}]]
  [SLOP {slop}] [INORDER]
  [LANGUAGE {language}]
  [EXPANDER {expander}]
  [SCORER {scorer}] [EXPLAINSCORE]
  [PAYLOAD {payload}]
  [SORTBY {attribute} [ASC|DESC]]
  [LIMIT offset num]

Parameters
index : The index name. The index must be first created with FT.CREATE .
query : the text query to search. If it's more than a single word, put it in quotes. Refer to query syntax for more details.

NOCONTENT : If it appears after the query, we only return the document ids and not the content. This is useful if RediSearch is only an index on an external document collection

VERBATIM : if set, we do not try to use stemming for query expansion but search the query terms verbatim.
NOSTOPWORDS : If set, we do not filter stopwords from the query.
WITHSCORES : If set, we also return the relative internal score of each document. this can be used to merge results from multiple instances
WITHPAYLOADS : If set, we retrieve optional document payloads (see FT.ADD). the payloads follow the document id, and if WITHSCORES was set, follow the scores.
WITHSORTKEYS : Only relevant in conjunction with SORTBY . Returns the value of the sorting key, right after the id and score and /or payload if requested. This is usually not needed by users, and exists for distributed search coordination purposes.

FILTER numeric_attribute min max : If set, and numeric_attribute is defined as a numeric attribute in FT.CREATE, we will limit results to those having numeric values ranging between min and max. min and max follow ZRANGE syntax, and can be -inf , +inf and use ( for exclusive ranges. Multiple numeric filters for different attributes are supported in one query.

GEOFILTER {geo_attribute} {lon} {lat} {radius} m|km|mi|ft : If set, we filter the results to a given radius from lon and lat. Radius is given as a number and units. See GEORADIUS for more details.
INKEYS {num} {attribute} ... : If set, we limit the result to a given set of keys specified in the list. the first argument must be the length of the list, and greater than zero. Non-existent keys are ignored - unless all the keys are non-existent.
INFIELDS {num} {attribute} ... : If set, filter the results to ones appearing only in specific attributes of the document, like title or URL . You must include num , which is the number of attributes you're filtering by. For example, if you request title and URL , then num is 2.

RETURN {num} {identifier} AS {property} ... : Use this keyword to limit which attributes from the document are returned. num is the number of attributes following the keyword. If num is 0, it acts like NOCONTENT . identifier is either an attribute name (for hashes and JSON) or a JSON Path expression for (JSON). property is an optional name used in the result. If not provided, the identifier is used in the result.

SUMMARIZE ... : Use this option to return only the sections of the attribute which contain the matched text. See Highlighting for more details
HIGHLIGHT ... : Use this option to format occurrences of matched text. See Highlighting for more details
SLOP {slop} : If set, we allow a maximum of N intervening number of unmatched offsets between phrase terms. (i.e the slop for exact phrases is 0)
INORDER : If set, and usually used in conjunction with SLOP, we make sure the query terms appear in the same order in the document as in the query, regardless of the offsets between them.
LANGUAGE {language} : If set, we use a stemmer for the supplied language during search for query expansion. If querying documents in Chinese, this should be set to chinese in order to properly tokenize the query terms. Defaults to English. If an unsupported language is sent, the command returns an error. See FT.ADD for the list of languages.

EXPANDER {expander} : If set, we will use a custom query expander instead of the stemmer. See Extensions .

SCORER {scorer} : If set, we will use a custom scoring function defined by the user. See Extensions .
EXPLAINSCORE : If set, will return a textual description of how the scores were calculated. Using this options requires the WITHSCORES option.
PAYLOAD {payload} : Add an arbitrary, binary safe payload that will be exposed to custom scoring functions. See Extensions .

SORTBY {attribute} [ASC|DESC] : If specified, the results are ordered by the value of this attribute. This applies to both text and numeric attributes.

LIMIT first num : Limit the results to the offset and number of results given. Note that the offset is zero-indexed. The default is 0 10, which returns 10 items starting from the first result.
*/

// units of Radius
type Unit string

const (
	KILOMETERS Unit = "km"
	METERS     Unit = "m"
	FEET       Unit = "ft"
	MILES      Unit = "mi"
)

// TODO: json entegrasyonu
type FtFilter struct {
	field        string
	min          float64
	exclusiveMin bool
	max          float64
	exclusiveMax bool
}

// Query Builder
type FtSearch struct {
	indexname    string
	query        string
	nocontent    bool
	verbatim     bool
	nostopwords  bool
	withscores   bool
	withpayloads bool
	withsortkeys bool
	filters      []FtFilter
	geofilter    struct {
		field  string
		lon    float64
		lat    float64
		radius float64
		unit   Unit
	}
	inkeys       []string
	infields     []string
	returnfields []string
	summarize    struct {
		fields    []string
		fragnum   int
		fragsize  int
		separator string
	}
	highlight struct {
		fields []string
		tags   struct {
			open  string
			close string
		}
	}
	slop     *int
	inorder  bool
	language string
	expander string
	scorer   string
	//explainscore bool
	payload []byte
	sortby  struct {
		attribute string
		asc       bool
	}
	limit struct {
		offset int64
		num    int64
	}
}

func NewFtSearch(indexName string) *FtSearch {
	return &FtSearch{
		indexname: indexName,
	}
}

func (fts *FtSearch) AddIndexName(name string) *FtSearch {
	fts.indexname = name

	return fts
}

func (fts *FtSearch) AddQuery(query string) *FtSearch {
	fts.query = query

	return fts
}

func (fts *FtSearch) AddNoContent(active bool) *FtSearch {
	fts.nocontent = active

	return fts
}

func (fts *FtSearch) AddVerbatim(active bool) *FtSearch {
	fts.verbatim = active

	return fts
}

func (fts *FtSearch) AddNoStopWords(active bool) *FtSearch {
	fts.nostopwords = active

	return fts
}

func (fts *FtSearch) AddWithScores(active bool) *FtSearch {
	fts.withscores = active

	return fts
}

func (fts *FtSearch) AddWithPayloads(active bool) *FtSearch {
	fts.withpayloads = active

	return fts
}

func (fts *FtSearch) AddWithSortKeys(active bool) *FtSearch {
	fts.withsortkeys = active

	return fts
}

func (fts *FtSearch) AddFilter(field string, min, max float64, excludeMin, excludeMax bool) *FtSearch {
	fts.filters = append(fts.filters, FtFilter{
		field:        field,
		min:          min,
		exclusiveMin: excludeMin,
		max:          max,
		exclusiveMax: excludeMax,
	})

	return fts
}

func (fts *FtSearch) AddGeoFilter(field string, lon, lat, radius float64, unit Unit) *FtSearch {
	fts.geofilter.field = field
	fts.geofilter.lon = lon
	fts.geofilter.lat = lat
	fts.geofilter.radius = radius
	fts.geofilter.unit = unit

	return fts
}

func (fts *FtSearch) AddInKeys(keys ...string) *FtSearch {
	fts.inkeys = keys

	return fts
}

func (fts *FtSearch) AddInFields(fields ...string) *FtSearch {
	fts.infields = fields

	return fts
}

func (fts *FtSearch) AddReturnFields(fields ...string) *FtSearch {
	fts.returnfields = fields

	return fts
}

func (fts *FtSearch) AddSummarize(fields []string, fragNum, fragSize int, separator string) *FtSearch {
	fts.summarize.fields = fields
	fts.summarize.fragnum = fragNum
	fts.summarize.fragsize = fragSize
	fts.summarize.separator = separator

	return fts
}

func (fts *FtSearch) AddHighlight(fields []string, openTag, closeTag string) *FtSearch {
	fts.highlight.fields = fields
	fts.highlight.tags.open = openTag
	fts.highlight.tags.close = closeTag

	return fts
}

func (fts *FtSearch) AddSlop(slop *int) *FtSearch {
	fts.slop = slop

	return fts
}

func (fts *FtSearch) AddInOrder(active bool) *FtSearch {
	fts.inorder = active

	return fts
}

func (fts *FtSearch) AddLanguage(lang string) *FtSearch {
	fts.language = lang

	return fts
}

func (fts *FtSearch) AddExpander(exp string) *FtSearch {
	fts.expander = exp

	return fts
}

func (fts *FtSearch) AddScorer(scorer string) *FtSearch {
	fts.scorer = scorer

	return fts
}

/*func (fts *FtSearch) AddExplainScore() {
	fts.explainscore = false
}*/

func (fts *FtSearch) AddPayload(payload []byte) *FtSearch {
	fts.payload = payload

	return fts
}

func (fts *FtSearch) AddSortBy(attr string, asc bool) *FtSearch {
	fts.sortby.attribute = attr
	fts.sortby.asc = asc

	return fts
}

func (fts *FtSearch) AddLimit(offset, num int64) *FtSearch {
	fts.limit.offset = offset
	fts.limit.num = num

	return fts
}

func (fts *FtSearch) Serialize() []interface{} {

	var queryCode []interface{}

	queryCode = append(queryCode, "FT.SEARCH")

	queryCode = append(queryCode, fts.indexname)

	if fts.query != "" {
		queryCode = append(queryCode, fts.query)
	}

	// NOTE: sadece kaç adet sonuç var bilmek istiyorsak limit parametresini 0 0 olarak yollamamız yeterli.
	if fts.nocontent {
		queryCode = append(queryCode, "NOCONTENT")
	}

	if fts.verbatim {
		queryCode = append(queryCode, "VERBATIM")
	}

	if fts.nostopwords {
		queryCode = append(queryCode, "NOSTOPWORDS")
	}

	if fts.withscores {
		queryCode = append(queryCode, "WITHSCORES")
	}

	if fts.withpayloads {
		queryCode = append(queryCode, "WITHPAYLOADS")
	}

	if fts.withsortkeys {
		queryCode = append(queryCode, "WITHSORTKEYS")
	}

	if fts.filters != nil && len(fts.filters) > 0 {

		for _, f := range fts.filters {
			queryCode = append(queryCode, "FILTER", f.field)

			if f.exclusiveMin {
				queryCode = append(queryCode, fmt.Sprintf("(%.f", f.min))
			} else if math.IsInf(f.min, 1) {
				queryCode = append(queryCode, "+inf")
			} else {
				queryCode = append(queryCode, fmt.Sprintf("%.f", f.min))
			}

			if f.exclusiveMax {
				queryCode = append(queryCode, fmt.Sprintf("(%.f", f.max))
			} else if math.IsInf(f.max, -1) {
				queryCode = append(queryCode, "-inf")
			} else {
				queryCode = append(queryCode, fmt.Sprintf("%.f", f.max))
			}
		}
	}

	if fts.geofilter.field != "" {
		queryCode = append(queryCode, "GEOFILTER", fts.geofilter.field, fts.geofilter.lon, fts.geofilter.lat, fts.geofilter.radius, fts.geofilter.unit)
	}

	if fts.inkeys != nil && len(fts.inkeys) > 0 {
		queryCode = append(queryCode, "INKEYS", len(fts.inkeys))
		for _, ik := range fts.inkeys {
			queryCode = append(queryCode, ik)
		}
	}

	if fts.returnfields != nil {
		queryCode = append(queryCode, "RETURN", len(fts.returnfields))
		for _, rf := range fts.returnfields {
			queryCode = append(queryCode, rf)
		}
	}

	if fts.summarize.fields != nil && len(fts.summarize.fields) > 0 {
		queryCode = append(queryCode, "SUMMIRIZE")

		queryCode = append(queryCode, "FIELDS", len(fts.summarize.fields))
		for _, sf := range fts.summarize.fields {
			queryCode = append(queryCode, sf)
		}

		if fts.summarize.fragnum > 0 {
			queryCode = append(queryCode, "FRAGS", fts.summarize.fragnum)
		}

		if fts.summarize.fragsize > 0 {
			queryCode = append(queryCode, "LEN", fts.summarize.fragsize)
		}

		if fts.summarize.separator != "" {
			queryCode = append(queryCode, "SEPARATOR", fts.summarize.separator)
		}
	}

	if fts.highlight.fields != nil && len(fts.highlight.fields) > 0 {
		queryCode = append(queryCode, "HIGHLIGHT")

		queryCode = append(queryCode, "FIELDS", len(fts.highlight.fields))
		for _, hf := range fts.highlight.fields {
			queryCode = append(queryCode, hf)
		}

		queryCode = append(queryCode, "TAGS", fts.highlight.tags.open, fts.highlight.tags.close)
	}

	if fts.slop != nil {
		queryCode = append(queryCode, "SLOP", *fts.slop)
	}

	if fts.language != "" {
		queryCode = append(queryCode, "LANGUAGE", fts.language)
	}

	if fts.expander != "" {
		queryCode = append(queryCode, "EXPANDER", fts.expander)
	}

	if fts.scorer != "" {
		queryCode = append(queryCode, "SCORER", fts.scorer)
	}

	// TODO: payload

	if fts.sortby.attribute != "" {
		queryCode = append(queryCode, "SORTBY", fts.sortby.attribute)
		if fts.sortby.asc {
			queryCode = append(queryCode, "ASC")
		} else {
			queryCode = append(queryCode, "DESC")
		}
	}

	//if !fts.nocontent {
	if fts.limit.num > 0 {
		queryCode = append(queryCode, "LIMIT", fts.limit.offset, fts.limit.num)
	} else {
		queryCode = append(queryCode, "LIMIT", 0, 10)
	}
	//}

	return queryCode
}
