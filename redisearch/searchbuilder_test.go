package redisearch

import (
	"reflect"
	"testing"
)

func TestSearchBuilder_Encode(t *testing.T) {
	type fields struct {
		Query string
		Attr  []string
	}
	type args struct {
		attr []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Simple Query Encode",
			fields: fields{
				Query: "",
				Attr:  []string{"dene", "obaraks", "query", "ne dedin sen", "type", "12"},
			},
			args: args{
				attr: []string{"dene", "obaraks", "query", "ne dedin sen", "type", "12"},
			},
			want: "@dene:obaraks @query:ne dedin sen @type:12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := SearchBuilder{
				Query: tt.fields.Query,
				Attr:  tt.fields.Attr,
			}
			if got := sb.Encode(tt.args.attr); got != tt.want {
				t.Errorf("SearchBuilder.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchBuilder_Decode(t *testing.T) {
	type fields struct {
		Query string
		Attr  []string
	}
	type args struct {
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Simple Query Decode",
			fields: fields{
				Query: "@dene:obaraks @query:ne dedin sen @type:12",
				Attr:  []string{"dene", "obaraks", "query", "ne dedin sen", "type", "12"},
			},
			args: args{
				query: "@dene:obaraks @query:ne dedin sen @type:12",
			},
			want: []string{"dene", "obaraks", "query", "ne dedin sen", "type", "12"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := SearchBuilder{
				Query: tt.fields.Query,
				Attr:  tt.fields.Attr,
			}
			if got := sb.Decode(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchBuilder.Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchBuilder_Decode2Map(t *testing.T) {
	type fields struct {
		Query string
		Attr  []string
	}
	type args struct {
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "Simple Query Decode2Map",
			fields: fields{
				Query: "@dene:obaraks @query:ne dedin sen @type:12",
				Attr:  []string{"dene", "obaraks", "query", "ne dedin sen", "type", "12"},
			},
			args: args{
				query: "@dene:obaraks @query:ne dedin sen @type:12",
			},
			want: map[string]string{"dene": "obaraks", "query": "ne dedin sen", "type": "12"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := SearchBuilder{
				Query: tt.fields.Query,
				Attr:  tt.fields.Attr,
			}
			if got := sb.Decode2Map(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchBuilder.Decode2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
