package redisearch

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

// Search Aggregate Builder
type FtAggregate struct {
}
