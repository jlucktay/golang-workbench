package collection_test

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"time"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
	"github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2/ord_array_linear"
)

const (
	SUCCESS = iota
	FAILURE
)

func TestDriver(t *testing.T) {
	names := []string{"Peter", "Sathish", "Wade", "Don", "Indrajit", "Rahul", "Sam", "Kevin"}
	testCases := []struct {
		desc       string
		collection p2.WordCollection
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
			b := new(bytes.Buffer)
			tC.collection.DisplayCollection(b)
			t.Logf("DisplayCollection buffer:\n%s", b)

			rand.Seed(time.Now().UnixNano())
			needle := names[rand.Intn(len(names))]
			t.Logf(`Searching for "%s": `, needle)
			if result := tC.collection.SearchCollection(needle); result == SUCCESS {
				t.Log("FOUND")
			} else {
				t.Fatal("NOT FOUND")
			}

			tC.collection.FreeCollection()
		})
	}
}

func BenchmarkOAL(b *testing.B) {
	benchmarkCollection(&ord_array_linear.OrdArrayLinear{}, b)
}

func benchmarkCollection(wc WordCollection, b *testing.B) {
	b.Logf("Implementation: %s\n", reflect.TypeOf(wc))

	for i := 0; i < b.N; i++ {
	}
}
