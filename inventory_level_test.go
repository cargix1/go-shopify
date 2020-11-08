package goshopify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestInventoryLevelListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	listURL := fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels.json", client.pathPrefix)

	// The strconv.Atoi error changed in go 1.8, 1.7 is still being tested/supported.
	limitConversionErrorMessage := `strconv.Atoi: parsing "invalid": invalid syntax`
	if runtime.Version()[2:5] == "1.7" {
		limitConversionErrorMessage = `strconv.ParseInt: parsing "invalid": invalid syntax`
	}

	cases := []struct {
		body                    string
		linkHeader              string
		expectedInventoryLevels []InventoryLevel
		expectedPagination      *Pagination
		expectedErr             error
	}{
		// Expect empty pagination when there is no link header
		{
			`{"inventory_levels": [{"available":1},{"available":2}]}`,
			"",
			[]InventoryLevel{{Available: 1}, {Available: 2}},
			new(Pagination),
			nil,
		},
		// Invalid link header responses
		{
			"{}",
			"invalid link",
			[]InventoryLevel(nil),
			nil,
			ResponseDecodingError{Message: "could not extract pagination link header"},
		},
		{
			"{}",
			`<:invalid.url>; rel="next"`,
			[]InventoryLevel(nil),
			nil,
			ResponseDecodingError{Message: "pagination does not contain a valid URL"},
		},
		{
			"{}",
			`<http://valid.url?%invalid_query>; rel="next"`,
			[]InventoryLevel(nil),
			nil,
			errors.New(`invalid URL escape "%in"`),
		},
		{
			"{}",
			`<http://valid.url>; rel="next"`,
			[]InventoryLevel(nil),
			nil,
			ResponseDecodingError{Message: "page_info is missing"},
		},
		{
			"{}",
			`<http://valid.url?page_info=foo&limit=invalid>; rel="next"`,
			[]InventoryLevel(nil),
			nil,
			errors.New(limitConversionErrorMessage),
		},
		// Valid link header responses
		{
			`{"inventory_levels": [{"available":1}]}`,
			`<http://valid.url?page_info=foo&limit=2>; rel="next"`,
			[]InventoryLevel{{Available: 1}},
			&Pagination{
				NextPageOptions: &ListOptions{PageInfo: "foo", Limit: 2},
			},
			nil,
		},
		{
			`{"inventory_levels": [{"available":2}]}`,
			`<http://valid.url?page_info=foo>; rel="next", <http://valid.url?page_info=bar>; rel="previous"`,
			[]InventoryLevel{{Available: 2}},
			&Pagination{
				NextPageOptions:     &ListOptions{PageInfo: "foo"},
				PreviousPageOptions: &ListOptions{PageInfo: "bar"},
			},
			nil,
		},
	}
	for i, c := range cases {
		response := &http.Response{
			StatusCode: 200,
			Body:       httpmock.NewRespBodyFromString(c.body),
			Header: http.Header{
				"Link": {c.linkHeader},
			},
		}

		httpmock.RegisterResponder("GET", listURL, httpmock.ResponderFromResponse(response))

		inventoryLevels, pagination, err := client.InventoryLevel.ListWithPagination(nil)
		if !reflect.DeepEqual(inventoryLevels, c.expectedInventoryLevels) {
			t.Errorf("test %d InventoryLevel.ListWithPagination inventory levels returned %+v, expected %+v", i, inventoryLevels, c.expectedInventoryLevels)
		}

		if !reflect.DeepEqual(pagination, c.expectedPagination) {
			t.Errorf(
				"test %d InventoryLevel.ListWithPagination pagination returned %+v, expected %+v",
				i,
				pagination,
				c.expectedPagination,
			)
		}

		if (c.expectedErr != nil || err != nil) && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"test %d InventoryLevel.ListWithPagination err returned %+v, expected %+v",
				i,
				err,
				c.expectedErr,
			)
		}
	}
}
