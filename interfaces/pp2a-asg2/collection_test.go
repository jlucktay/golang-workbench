package main

import (
	"reflect"
	"testing"

	"github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2/ord_array_linear"
)

const (
	FAILURE = iota
	SUCCESS
)

func TestDriver(t *testing.T) {
	names := []string{"Peter", "Sathish", "Wade", "Don", "Indrajit", "Rahul", "Sam", "Kevin"}
	testCases := []struct {
		desc       string
		collection WordCollection
	}{
		{
			desc:       "Ordered slice with linear search",
			collection: &ord_array_linear.OrdArrayLinear{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s\n", reflect.TypeOf(tC.collection))

			if tC.collection.MakeCollection() == FAILURE {
				t.Fatal("Unable to initialise WordCollection")
			}

			for _, name := range names {
				if tC.collection.AddCollection(name) == FAILURE {
					t.Fatal("AddCollection failed")
				}
			}

			t.Logf("Collection contains %d names\n", tC.collection.SizeCollection())

			t.Log("The following names are in the Collection:")
			tC.collection.DisplayCollection()

			t.Log(`Searching for "Sathish": `)
			if result := tC.collection.SearchCollection("Sathish"); result == SUCCESS {
				t.Log("FOUND")
			} else {
				t.Log("NOT FOUND")
			}

			tC.collection.FreeCollection()
		})
	}
}
