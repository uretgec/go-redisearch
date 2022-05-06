# go-redisearch
Redisearch go client with go-redis options

For usage example go to examples folder

## Include
- Search Builder
- Raw query
- Go-redis integration

## Install

```
go get  github.com/uretgec/go-redisearch
```

## Examples

```
// Constants
const INDEX_DREAMS = "index_dreams"
const DREAM_DIC_KEY = "dreamdic"

// Redis Options
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

```

## TODO
- Use client methods
- Add test files
- Add new examples

## Links

Go-Redis (https://github.com/go-redis/redis)