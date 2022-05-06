package redisearch

import (
	"fmt"
	"strings"
)

const (
	QUERY_TYPE_1 = 1
	QUERY_TYPE_2 = 2
	QUERY_TYPE_3 = 3
	QUERY_TYPE_4 = 4
	QUERY_TYPE_5 = 5
	QUERY_TYPE_6 = 6
	QUERY_TYPE_7 = 7
	QUERY_TYPE_8 = 8
)

// TODO: json entegrasyonu
type FtQuery struct {
	raw   string
	query []string
}

func NewFtQuery(raw string) *FtQuery {
	return &FtQuery{
		raw:   raw,
		query: make([]string, 0),
	}
}

// Single word query generator
func (ftq *FtQuery) GenerateSingleWordQuery(queryType int, word string) string {

	switch queryType {
	case QUERY_TYPE_1: // Exact phrase query - hello FOLLOWED BY world => "hello world"
		return word

	case QUERY_TYPE_2: // Not: documents containing hello but not world => hello -world
		return "-" + word

	case QUERY_TYPE_3: // Intersection of unions query => (hello|halo) (world|werld)
		return "(" + word + ")"

	case QUERY_TYPE_4: // Negation of union query => hello -(world|werld)
		return "-(" + word + ")"

	case QUERY_TYPE_5: // Prefix Queries: hell* world  => (hello|help|helm|...) world => hello worl* hel* worl* hello -worl*
		word := word[:len(word)-1]
		return word + "*"

	case QUERY_TYPE_6: // Prefix Queries: hell* world  => (hello|help|helm|...) world => hello worl* hel* worl* hello -worl*
		word := word[:len(word)-2]
		return word + "*"
	}

	return word
}

// Multiple words query generator
func (ftq *FtQuery) GenerateMultipleWordQuery(queryType int, words ...string) string {

	switch queryType {
	case QUERY_TYPE_7: // Union: documents containing either hello OR world => hello|world
		return strings.Join(words[:], "|")

	case QUERY_TYPE_8: // Optional terms with higher priority to ones containing more matches => obama ~barack ~michelle
		return strings.Join(words[:], "~")

	}

	return strings.Join(words[:], " ")
}

// %hello% world
func (ftq *FtQuery) GenerateFuzzyMatchQuery(ld int, words ...string) string {
	var nfq []string

	for _, word := range words {
		var prefix string
		for i := 0; i < ld; i++ {
			prefix += "%"
		}
		nfq = append(nfq, prefix+word+prefix)
	}

	return strings.Join(nfq, " ")
}

func (ftq *FtQuery) GenerateFieldModifyQuery(field string, query string) string {
	return "@" + field + ":" + query
}

func (ftq *FtQuery) GenerateMultiFieldsModifyQuery(fields []string, query string) string {
	return "@" + strings.Join(fields, "|") + ":" + query
}

// Too dangerous query: performance killer
func (ftq *FtQuery) AddWildcardQuery() *FtQuery {

	ftq.query = []string{"*"}

	return ftq
}

func (ftq *FtQuery) AddPureNegativeQuery(query string) *FtQuery {
	var nfq string

	nfq += "-" + query

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// @field:[{min} {max}]
// -inf , inf and +inf are acceptable numbers in a range. Thus greater-than 100 is expressed as [(100 inf]
// Numeric filters are inclusive. Exclusive min or max are expressed with ( prepended to the number, e.g. [(100 (200]
// It is possible to negate a numeric filter by prepending a - sign to the filter, e.g. returning a result where price differs from 100 is expressed as: @title:foo -@price:[100 100]
func (ftq *FtQuery) AddNumericFilterQuery(isNegative bool, field string, min, max int64, excludeMin, excludeMax bool, infMin, infMax bool) *FtQuery {
	var nfq string

	if isNegative {
		nfq += "-"
	}

	nfq += "@" + field + ":["

	if excludeMin {
		nfq += fmt.Sprintf("(%d", min)
	} else if infMin {
		nfq += "inf"
	} else {
		nfq += fmt.Sprintf("%d", min)
	}

	nfq += " "
	if excludeMax {
		nfq += fmt.Sprintf("(%d", max)
	} else if infMax {
		nfq += "inf"
	} else {
		nfq += fmt.Sprintf("%d", max)
	}

	nfq += "]"

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// @field:{ tag | tag | ...}
func (ftq *FtQuery) AddTagFilterQuery(isNegative bool, field string, tags ...string) *FtQuery {
	var nfq string

	if isNegative {
		nfq += "-"
	}

	nfq += "@" + field + ":{"

	// Via: https://stackoverflow.com/a/28799151
	nfq += strings.Join(tags[:], "|")

	nfq += "}"

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// @field:[{lon} {lat} {radius} {m|km|mi|ft}]
func (ftq *FtQuery) AddGeoFilterQuery(isNegative bool, field string, lon, lat, radius float64, unit Unit) *FtQuery {
	var nfq string

	if isNegative {
		nfq += "-"
	}

	nfq += fmt.Sprintf("@%s:[%.f %.f %.f %s]", field, lon, lat, radius, unit)

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// hel* world
func (ftq *FtQuery) AddPrefixMatchQuery(isNegative bool, field string, query string) *FtQuery {
	var nfq string

	if isNegative {
		nfq += "-"
	}

	if field != "" {
		nfq += "@" + field + ":"
	}

	nfq += query

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// @fileds|fields2|fields3:hell* world
func (ftq *FtQuery) AddMultiFieldsPrefixMatchQuery(isNegative bool, fields []string, query string) *FtQuery {
	var nfq string

	if isNegative {
		nfq += "-"
	}

	if len(fields) > 0 {
		nfq += "@" + strings.Join(fields, "|") + ":"
	}

	nfq += query

	ftq.query = append(ftq.query, nfq)

	return ftq
}

// Not ready to use
func (ftq *FtQuery) AddAttributeQuery() *FtQuery {
	return ftq
}

func (ftq *FtQuery) Serialize() string {
	ftq.raw = strings.Join(ftq.query, " ")
	return ftq.raw
}
