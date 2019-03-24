package sql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var queryString = `{
  "condition": "OR",
  "rules": [
    {
      "id": "price",
      "field": "price",
      "type": "double",
      "input": "number",
      "operator": "equal",
      "value": 10.25
    },
    {
      "id": "category",
      "field": "category",
      "type": "integer",
      "input": "select",
      "operator": "not_equal",
      "value": 3
    },
    {
      "condition": "AND",
      "rules": [
        {
          "id": "price",
          "field": "price",
          "type": "double",
          "input": "number",
          "operator": "equal",
          "value": 1
        }
      ]
    },
    {
      "id": "category",
      "field": "category",
      "type": "integer",
      "input": "select",
      "operator": "not_equal",
      "value": null
    }
  ]
}`

var expectedSQL = `("price" = 10.25 OR "category" != 3 OR ("price" = 1) OR "category" != null)`

func TestNode_SQL(t *testing.T) {
	var node Node
	err := json.Unmarshal([]byte(queryString), &node)
	if err != nil {
		t.Fatal(err)
	}

	sql, err := node.SQL()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedSQL, sql)
}
