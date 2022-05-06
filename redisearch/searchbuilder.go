package redisearch

import (
	"encoding/json"
	"regexp"
	"strings"
)

type SearchBuilder struct {
	Raw    string
	Query  string
	Cats   string
	Tags   string
	SortBy string // all|latest(updated)
	Attr   []string
}

func NewSearchBuilder() *SearchBuilder {
	return &SearchBuilder{}
}

func (sb *SearchBuilder) MarshalBinary() ([]byte, error) {
	return json.Marshal(sb)
}

func (sb *SearchBuilder) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &sb); err != nil {
		return err
	}

	return nil
}

// Additional Methods: Only for send special query
func (sb *SearchBuilder) Encode(attr []string) string {
	var query []string
	for i := 0; i < len(attr); i += 2 {
		query = append(query, "@"+attr[i]+":"+attr[i+1])
	}

	return strings.Join(query, " ")
}

func (sb *SearchBuilder) Decode(query string) []string {
	// Regex: `([^@:]+):([^@]+)?`g
	re := regexp.MustCompile(`([^@:]+):([^@]+)?`)
	res := re.FindAllStringSubmatch(query, -1)

	if len(res) > 0 {
		var query []string
		for _, item := range res {
			query = append(query, item[1], strings.TrimSpace(item[2]))
		}

		return query
	}

	return nil
}

func (sb *SearchBuilder) Decode2Map(query string) map[string]string {
	// Regex: `([^@:]+):([^@]+)?`g
	re := regexp.MustCompile(`([^@:]+):([^@]+)?`)
	res := re.FindAllStringSubmatch(query, -1)

	if len(res) > 0 {
		query := make(map[string]string)
		for _, item := range res {
			query[item[1]] = strings.TrimSpace(item[2])
		}

		return query
	}

	return nil
}
