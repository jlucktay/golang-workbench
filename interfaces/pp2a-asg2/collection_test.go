package collection_test

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"time"

	p2 "github.com/jlucktay/golang-workbench/interfaces/pp2a-asg2"
	"github.com/matryer/is"
)

func TestDriver(t *testing.T) {
	names := []string{"Peter", "Sathish", "Wade", "Don", "Indrajit", "Rahul", "Sam", "Kevin"}
	testCases := []struct {
		desc       string
		collection p2.WordCollection
	}{
		{
			desc:       "Ordered slice with linear search",
			collection: &p2.OrdArrayLinear{},
		},
		{
			desc:       "Ordered linked list with linear search",
			collection: &p2.OrdLinkedList{},
		},
		{
			desc:       "Ordered slice with binary search",
			collection: &p2.OrdArrayBinary{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Logf("Current implementation based on: %s", reflect.TypeOf(tC.collection))
			i := is.New(t)

			i.Equal(tC.collection.MakeCollection(), SUCCESS)

			for _, name := range names {
				i.Equal(tC.collection.AddCollection(name), SUCCESS)
			}

			t.Logf("Collection contains %d names", tC.collection.SizeCollection())
			i.Equal(len(names), tC.collection.SizeCollection())

			t.Log("The following names are in the Collection:")
			b := new(bytes.Buffer)
			tC.collection.DisplayCollection(b)
			t.Logf("DisplayCollection buffer:\n%s", b)

			rand.Seed(time.Now().UnixNano())
			needle := names[rand.Intn(len(names))]
			t.Logf(`Searching for "%s"...`, needle)
			i.Equal(tC.collection.SearchCollection(needle), SUCCESS)

			tC.collection.FreeCollection()
		})
	}
}
