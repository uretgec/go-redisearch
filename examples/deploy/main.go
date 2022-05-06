package main

import (
	"errors"
	"fmt"

	"github.com/uretgec/go-redisearch/redisearch"
)

const (
	INDEX_DREAMS = "index_dreams"
	INDEX_TERMS  = "index_terms"
)

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

	// Create Dream Index
	args, err := generateIndexQuery(INDEX_DREAMS)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	err = client.UClient.Do(client.Ctx, args...).Err()
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	fmt.Println("Dream index created")

	// Create Term Index
	args, err = generateIndexQuery(INDEX_TERMS)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	err = client.UClient.Do(client.Ctx, args...).Err()
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}

	fmt.Println("Term index created")

	fmt.Println("bye bye")
}

func generateIndexQuery(indexName string) ([]interface{}, error) {
	switch indexName {
	case INDEX_DREAMS:
		indexCreator := redisearch.NewFtCreate(indexName)
		//indexCreator.AddTemporarySeconds(true, 3600)
		indexCreator.AddDataType(redisearch.HASH)
		indexCreator.AddPrefix("drd:")
		//indexCreator.AddScoreField("score")

		// uid
		schemaTextOpt := indexCreator.AddSchemaTextOption(0, false, false, "")
		indexCreator.AddSchema(redisearch.FieldTypeText, "uid", "", false, schemaTextOpt)

		// name
		indexCreator.AddSchema(redisearch.FieldTypeText, "name", "", false, schemaTextOpt)

		// slug
		//indexCreator.AddSchema(redisearch.FieldTypeText, "slug", "", false, schemaTextOpt)

		// description
		indexCreator.AddSchema(redisearch.FieldTypeText, "description", "", false, schemaTextOpt)

		// cats
		indexCreator.AddSchema(redisearch.FieldTypeText, "cats", "", false, schemaTextOpt)

		// tags
		indexCreator.AddSchema(redisearch.FieldTypeText, "tags", "", false, schemaTextOpt)

		// updated
		schemaNumericOpt := indexCreator.AddSchemaNumericOption(false)
		indexCreator.AddSchema(redisearch.FieldTypeNumeric, "updated", "", true, schemaNumericOpt)

		return indexCreator.Serialize(), nil
	case INDEX_TERMS:
		indexCreator := redisearch.NewFtCreate(indexName)
		//indexCreator.AddTemporarySeconds(true, 3600)
		indexCreator.AddDataType(redisearch.HASH)
		indexCreator.AddPrefix("drd:")
		//indexCreator.AddScoreField("score")

		// uid
		schemaTextOpt := indexCreator.AddSchemaTextOption(0, false, false, "")
		indexCreator.AddSchema(redisearch.FieldTypeText, "uid", "", false, schemaTextOpt)

		// cats
		schemaTagOpt := indexCreator.AddSchemaTagOption(false, ",")
		indexCreator.AddSchema(redisearch.FieldTypeTag, "cats", "", false, schemaTagOpt)

		// tags
		indexCreator.AddSchema(redisearch.FieldTypeTag, "tags", "", false, schemaTagOpt)

		// updated
		schemaNumericOpt := indexCreator.AddSchemaNumericOption(false)
		indexCreator.AddSchema(redisearch.FieldTypeNumeric, "updated", "", true, schemaNumericOpt)

		return indexCreator.Serialize(), nil

	}

	return nil, errors.New("index name not found")
}
