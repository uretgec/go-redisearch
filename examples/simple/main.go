package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uretgec/go-redisearch/redisearch"
)

const (
	INDEX_DREAMS = "index_dreams"
	INDEX_TERMS  = "index_terms"
)

const DRSEARCH_PATTERN = "[uid]"
const DREAM_DIC_KEY = "dreamdic"

// Dream
// CMD: HGETALL
// Key: drsearch:[uid]
const DRSEARCH_DETAIL_KEY = "drd:[uid]"

type DreamSearch struct {
	UID         string `json:"uid" redis:"uid"`
	Name        string `json:"name" redis:"name"`
	Slug        string `json:"slug" redis:"slug"`
	Updated     int64  `json:"updated" redis:"updated"`
	Description string `json:"description" redis:"description"`
	Categories  string `json:"cats" redis:"cats"`
	Tags        string `json:"tags" redis:"tags"`
}

func (bbi *DreamSearch) MarshalBinary() ([]byte, error) {
	return json.Marshal(bbi)
}

func (bbi *DreamSearch) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &bbi); err != nil {
		return err
	}

	return nil
}

func main() {

	// Connect Redisearch
	redisAddrs := []string{"127.0.0.1:6481"}
	redisPoolSizes := []int{1000}
	redisMinIdleConns := []int{2}
	redisMaxRetries := []int{2}
	client := redisearch.NewRedisearchClient("myredisearch", redisAddrs, redisPoolSizes, redisMinIdleConns, redisMaxRetries)
	defer client.CloseUniversalClient()

	// Check Redis is here :)
	err := client.HealthCheckedUniversalClient()
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	// Generate New Object and Add Something to search
	item := &DreamSearch{
		UID:         fmt.Sprintf("%d", 123456),
		Updated:     time.Now().Unix(),
		Name:        fmt.Sprintf("%s %d", "Different Test", 123456),
		Slug:        fmt.Sprintf("%s-%d", "different-test", 123456),
		Description: "Test different model set action.",
		Categories:  "Test Category, Different model, set action",
		Tags:        "Test, Different, Default",
	}

	paramsMap := convert2RedisMap(item)
	if err := client.UClient.HSet(
		client.Ctx,
		generateRedisKey(DRSEARCH_DETAIL_KEY, DRSEARCH_PATTERN, item.UID),
		paramsMap,
	).Err(); err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	for i := 0; i < 10; i++ {
		item := &DreamSearch{
			UID:         fmt.Sprintf("%d", i),
			Updated:     time.Now().Unix(),
			Name:        fmt.Sprintf("%s %d", "Dream Test", i),
			Slug:        fmt.Sprintf("%s-%d", "dream-test", i),
			Description: "Test dream model set action.",
			Categories:  "Test Category, Dream model, set action",
			Tags:        "Test, Dream, Default",
		}

		paramsMap := convert2RedisMap(item)
		if err := client.UClient.HSet(
			client.Ctx,
			generateRedisKey(DRSEARCH_DETAIL_KEY, DRSEARCH_PATTERN, item.UID),
			paramsMap,
		).Err(); err != nil {
			fmt.Printf("error: %v\n", err)
			panic(err)
		}
	}

	// Add dic
	pipe := client.UClient.Pipeline()

	for _, tag := range []string{"Test", "Dream", "Diffirent", "Tag", "Default"} {
		pipe.Do(client.Ctx, "FT.SUGADD", DREAM_DIC_KEY, tag, 1, "INCR")
	}

	_, err = pipe.Exec(client.Ctx)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	// Suggest Search
	suggestions, err := client.UClient.Do(client.Ctx, "FT.SUGGET", DREAM_DIC_KEY, "Dif", "FUZZY", "MAX", 5).StringSlice()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	fmt.Printf("Suggestions: %v\n", suggestions)

	// Search Builder Init
	searchBuilder := redisearch.NewSearchBuilder()
	searchBuilder.Query = "different" // return only one result
	//searchBuilder.Query = "dream" //  return all result no special one
	searchBuilder.SortBy = "latest"

	searchQuery, err := generateSearchQuery(INDEX_DREAMS, searchBuilder, int64(0), int64(10))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	results, err := client.UClient.Do(client.Ctx, searchQuery...).Slice()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	var total int64
	uids := []string{}

	if len(results) > 0 {
		for i := 0; i < len(results); i++ {
			if i == 0 {
				total = results[0].(int64)
				if total == 0 {
					break
				} else {
					continue
				}
			}

			uids = append(uids, strings.Replace(results[i].(string), "drd:", "", 1))
		}
	}

	fmt.Printf("Total: %d, Results: %v\n", int(total), uids)
}

func generateSearchQuery(indexName string, searchBuilder *redisearch.SearchBuilder, offset, limit int64) ([]interface{}, error) {

	switch indexName {
	case INDEX_DREAMS:

		// query builder
		indexQuery := redisearch.NewFtQuery("")

		// "@name:(asd|dsa|"asd dsa") @cats|tags:{asd|dsa|"asd dsa"} @is_purchased:true|false"
		addQuery := searchBuilder.Raw
		if searchBuilder.Raw == "" {
			var q string
			words := strings.Split(searchBuilder.Query, "+")
			if len(words) > 1 {
				words = append(words, indexQuery.GenerateSingleWordQuery(redisearch.QUERY_TYPE_1, strings.Join(words, " ")))
				word := indexQuery.GenerateMultipleWordQuery(redisearch.QUERY_TYPE_7, words...)
				q = indexQuery.GenerateSingleWordQuery(redisearch.QUERY_TYPE_3, word)
			} else {
				q = indexQuery.GenerateSingleWordQuery(redisearch.QUERY_TYPE_5, searchBuilder.Query)
			}

			indexQuery.AddMultiFieldsPrefixMatchQuery(false, []string{"name", "slug", "description", "cats", "tags"}, q)

			addQuery = indexQuery.Serialize()
		}

		// search builder
		indexSearcher := redisearch.NewFtSearch(indexName)
		indexSearcher.AddQuery(addQuery)

		if searchBuilder.SortBy == "latest" {
			indexSearcher.AddSortBy("updated", false) // DESC
		}

		indexSearcher.AddNoContent(true)
		indexSearcher.AddNoStopWords(true)
		indexSearcher.AddLimit(offset, limit)

		return indexSearcher.Serialize(), nil
	case INDEX_TERMS:
		// query builder
		indexQuery := redisearch.NewFtQuery("")

		addQuery := searchBuilder.Raw
		if searchBuilder.Raw == "" {
			isDummyQuery := true
			if searchBuilder.Cats != "" {
				catQuery := indexQuery.GenerateSingleWordQuery(redisearch.QUERY_TYPE_1, searchBuilder.Cats)
				indexQuery.AddTagFilterQuery(false, "cats", catQuery)

				isDummyQuery = false
			}

			if searchBuilder.Tags != "" {
				tagQuery := indexQuery.GenerateSingleWordQuery(redisearch.QUERY_TYPE_1, searchBuilder.Tags)
				indexQuery.AddTagFilterQuery(false, "tags", tagQuery)

				isDummyQuery = false
			}

			if isDummyQuery {
				if searchBuilder.SortBy == "latest" {
					now := time.Now()
					min := now.AddDate(-1, 0, 0)

					indexQuery.AddNumericFilterQuery(false, "updated", min.Unix(), -1, false, false, false, true)
				}
			}

			addQuery = indexQuery.Serialize()
		}

		// search builder
		indexSearcher := redisearch.NewFtSearch(indexName)
		indexSearcher.AddQuery(addQuery)

		if searchBuilder.SortBy == "latest" {
			indexSearcher.AddSortBy("updated", false) // DESC
		}

		indexSearcher.AddNoContent(true)
		indexSearcher.AddNoStopWords(true)
		indexSearcher.AddLimit(offset, limit)
		//indexSearcher.AddReturnFields([]string{"uid", "name", "slug", "description", "cats"}...)

		return indexSearcher.Serialize(), nil
	}

	return nil, errors.New("index name not found")
}

func convert2RedisMap(item interface{}) map[string]interface{} {
	var itemMap map[string]interface{}
	data, _ := json.Marshal(item)
	json.Unmarshal(data, &itemMap)

	return itemMap
}

func generateRedisKey(key, pattern, val string) string {
	return strings.ReplaceAll(key, pattern, val)
}
