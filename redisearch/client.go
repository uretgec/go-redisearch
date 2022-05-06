// Via: https://oss.redis.com/redisearch/Commands/#ftcreate
// Master Slave: https://medium.com/@ashu.goldi/redis-cluster-with-redisearch-setup-227e542c6746
package redisearch

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

// NewUniversalClient returns a new multi client. The type of the returned client depends
// on the following conditions:
//
// 1. If the MasterName option is specified, a sentinel-backed FailoverClient is returned.
// 2. if the number of Addrs is two or more, a ClusterClient is returned.
// 3. Otherwise, a single-node Client is returned.
type RedisearchClient struct {
	Name    string
	UClient redis.UniversalClient
	Ctx     context.Context
}

func NewRedisearchClient(serviceName string, redisAddrs []string, redisPoolSizes, redisMinIdleConns, redisMaxRetries []int) *RedisearchClient {

	uClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:         redisAddrs,
		RouteRandomly: true,
		MaxRetries:    redisMaxRetries[0],
		PoolSize:      redisPoolSizes[0], // Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
		MinIdleConns:  redisMinIdleConns[0],
	})

	return &RedisearchClient{
		Name:    serviceName,
		UClient: uClient,
		Ctx:     context.Background(),
	}
}

func (rsc *RedisearchClient) CloseUniversalClient() error {
	return rsc.UClient.Close()
}

func (rsc *RedisearchClient) HealthCheckedUniversalClient() error {
	if _, err := rsc.UClient.Ping(rsc.Ctx).Result(); err != nil {
		return errors.New("redisearch cluster server not response")
	}

	return nil
}

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
func (rsc *RedisearchClient) Create(indexName string, args ...interface{}) (string, error) {
	return rsc.UClient.Do(rsc.Ctx, args...).Text()
}

/*
HSET, HMSET, HSETNX, HINCRBY, HINCRBYFLOAT, HDEL, DEL, SET, RENAME_FROM, RENAME_TO, TRIMMED, RESTORE, EXPIRED,
EVICTED, CHANGE, LOADED, JSON.SET, JSON.DEL, JSON.NUMINCRBY, JSON.ARRAPPEND, JSON.ARRINDEDX, JSON.ARRTRIM, JSON.ARRPOP

Beginning with RediSearch v2.0, you use native Redis commands to add, update or delete hashes. These include HSET , HINCRBY , HDEL .
*/
func (rsc *RedisearchClient) HSet(key string, values ...interface{}) (int64, error) {
	return rsc.UClient.HSet(rsc.Ctx, key, values...).Result()
}
func (rsc *RedisearchClient) HMSet(key string, values ...interface{}) (bool, error) {
	return rsc.UClient.HMSet(rsc.Ctx, key, values...).Result()
}
func (rsc *RedisearchClient) HDel(key string, fields ...string) (int64, error) {
	return rsc.UClient.HDel(rsc.Ctx, key, fields...).Result()
}
func (rsc *RedisearchClient) HGet(key string, field string) (string, error) {
	return rsc.UClient.HGet(rsc.Ctx, key, field).Result()
}
func (rsc *RedisearchClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	return rsc.UClient.HMGet(rsc.Ctx, key, fields...).Result()
}
func (rsc *RedisearchClient) HGetAll(key string) (map[string]string, error) {
	return rsc.UClient.HGetAll(rsc.Ctx, key).Result()
}
func (rsc *RedisearchClient) Del(keys ...string) (int64, error) {
	return rsc.UClient.Del(rsc.Ctx, keys...).Result()
}

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
*/
func (rsc *RedisearchClient) Search(indexName string, query ...interface{}) (interface{}, error) {
	return rsc.UClient.Do(rsc.Ctx, query...).Result()
}

/*
FT.AGGREGATE {index_name}
  {query_string}
  [VERBATIM]
  [LOAD {nargs} {identifier} [AS {property}] ...]
  [GROUPBY {nargs} {property} ...
    REDUCE {func} {nargs} {arg} ... [AS {name:string}]
    ...
  ] ...
  [SORTBY {nargs} {property} [ASC|DESC] ... [MAX {num}]]
  [APPLY {expr} AS {alias}] ...
  [LIMIT {offset} {num}] ...
  [FILTER {expr}] ...
*/
func (rsc *RedisearchClient) Aggregate(indexName string) error {
	return errors.New("not ready to use")
}

/*
FT.EXPLAIN {index} {query}
*/
func (rsc *RedisearchClient) Explain(indexName string, query string) error {
	return errors.New("not ready to use")
}

/*
FT.PROFILE {index} {[SEARCH, AGGREGATE]} [LIMITED] QUERY {query}
*/
func (rsc *RedisearchClient) Profile(indexName string) error {
	return errors.New("not ready to use")
}

/*
FT.ALTER {index} SCHEMA ADD {attribute} {options} ...
Adds a new attribute to the index.
*/
func (rsc *RedisearchClient) Alter(indexName string, values ...interface{}) (string, error) {
	values = append([]interface{}{"FT.ALTER", indexName, "SCHEMA", "ADD"}, values...)
	return rsc.UClient.Do(rsc.Ctx, values...).Text()
}

/*
FT.DROPINDEX {index} [DD]
Deletes the index.
By default, FT.DROPINDEX does not delete the document hashes associated with the index. Adding the DD option deletes the hashes as well.
*/
func (rsc *RedisearchClient) DropIndex(indexName string, deleteHash bool) (string, error) {
	if deleteHash {
		return rsc.UClient.Do(rsc.Ctx, "FT.DROPINDEX", indexName, "DD").Text()
	}

	return rsc.UClient.Do(rsc.Ctx, "FT.DROPINDEX", indexName).Text()
}

/*
FT.ALIASADD {name} {index}
FT.ALIASUPDATE {name} {index}
FT.ALIASDEL {name}

Indexes can have more than one alias, though an alias cannot refer to another alias.
*/
func (rsc *RedisearchClient) AliasAdd(name string, indexName string) (string, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.ALIASADD", name, indexName).Text()
}
func (rsc *RedisearchClient) AliasUpdate(name string, indexName string) (string, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.ALIASUPDATE", name, indexName).Text()
}
func (rsc *RedisearchClient) AliasDel(name string) (string, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.ALIASDEL", name).Text()
}

/*
FT.TAGVALS {index} {attribute_name}
*/
func (rsc *RedisearchClient) TagVals(indexName string, attr string) error {
	return errors.New("not ready to use")
}

/*
FT.SUGADD {key} {string} {score} [INCR] [PAYLOAD {payload}]
Adds a suggestion string to an auto-complete suggestion dictionary. This is disconnected from the index definitions, and leaves creating and updating suggestions dictionaries to the user.
*/
func (rsc *RedisearchClient) SugAdd(key string, val string) (int64, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.SUGADD", key, val, 1, "INCR").Int64()
}

/*
FT.SUGGET {key} {prefix} [FUZZY] [WITHSCORES] [WITHPAYLOADS] [MAX num]
Gets completion suggestions for a prefix.
*/
func (rsc *RedisearchClient) SugGet(key string, prefix string, max int) (interface{}, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.SUGGET", key, prefix, "FUZZY", "MAX", max).Result()
}

/*
FT.SUGDEL {key} {string}
Deletes a string from a suggestion index.
*/
func (rsc *RedisearchClient) SugDel(key string, val string) (int64, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.SUGLEN", key).Int64()
}

/*
FT.SUGLEN {key}
Gets the size of an auto-complete suggestion dictionary
*/
func (rsc *RedisearchClient) SugLen(key string) (int64, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.SUGLEN", key).Int64()
}

/*
FT.SYNUPDATE <index name> <synonym group id> [SKIPINITIALSCAN] <term1> <term2> ...
*/
func (rsc *RedisearchClient) SynUpdate(indexName string) error {
	return errors.New("not ready to use")
}

/*
FT.SYNDUMP <index name>
*/
func (rsc *RedisearchClient) SynDump(indexName string) error {
	return errors.New("not ready to use")
}

/*
FT.SPELLCHECK {index} {query}
    [DISTANCE dist]
    [TERMS {INCLUDE | EXCLUDE} {dict} [TERMS ...]]
*/
func (rsc *RedisearchClient) SpellCheck(indexName string, query string) error {
	return errors.New("not ready to use")
}

/*
FT.DICTADD {dict} {term} [{term} ...]
Adds terms to a dictionary.
*/
func (rsc *RedisearchClient) DictAdd(dict string, terms []interface{}) (int64, error) {
	terms = append([]interface{}{"FT.DICTADD", dict}, terms...)
	return rsc.UClient.Do(rsc.Ctx, terms...).Int64()
}

/*
FT.DICTDEL {dict} {term} [{term} ...]
Deletes terms from a dictionary.
*/
func (rsc *RedisearchClient) DictDel(dict string, terms []interface{}) (int64, error) {
	terms = append([]interface{}{"FT.DICTDEL", dict}, terms...)
	return rsc.UClient.Do(rsc.Ctx, terms...).Int64()
}

/*
FT.DICTDUMP {dict}
Dumps all terms in the given dictionary.
*/
func (rsc *RedisearchClient) DictDump(dict string) error {
	return errors.New("not ready to use")
}

/*
FT.INFO {index}
Returns information and statistics on the index.
*/
func (rsc *RedisearchClient) Info(indexName string) (string, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT.INFO", indexName).Text()
}

/*
FT._LIST
Returns a list of all existing indexes.
*/
func (rsc *RedisearchClient) List() (interface{}, error) {
	return rsc.UClient.Do(rsc.Ctx, "FT._LIST").Result()
}

/*
FT.CONFIG <GET|HELP> {option}
FT.CONFIG SET {option} {value}
*/
func (rsc *RedisearchClient) ConfigHelp(option string) error {
	return errors.New("not ready to use")
}
func (rsc *RedisearchClient) ConfigGet(option string) error {
	return errors.New("not ready to use")
}
func (rsc *RedisearchClient) ConfigSet(option string, val string) error {
	return errors.New("not ready to use")
}
