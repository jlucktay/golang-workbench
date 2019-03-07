// ord_array_linear is an ordered slice with linear search
package ord_array_linear

import "fmt"

const (
	FAILURE = iota
	SUCCESS
)

const WCSIZE = 250000

type OrdArrayLinear struct {
	words []string
	size  int
}

/*
 * MakeCollection initialises the WordCollection whose pointer it is given as
 * a parameter, up to a size defined by WCSIZE, with NULL entries. It returns
 * SUCCESS upon successful completion.
 */
func (o *OrdArrayLinear) MakeCollection() int {
	o.words = make([]string, 0, WCSIZE)
	o.size = 0

	return SUCCESS
}

/*
 * FreeCollection frees the memory dynamically allocated to the WordCollection
 * parameter.
 */
func (o *OrdArrayLinear) FreeCollection() {
	o.words = nil
	o.size = 0
}

/*
 * AddCollection adds the string, given by the second parameter, to the
 * WordCollection given by the first parameter. It returns SUCCESS or FAILURE,
 * depending on whether or not there is space for the string in the
 * WordCollection, and also on the outcome of the dynamic allocation of
 * memory. The string is added so that the WordCollection is in alphabetical
 * order.
 */
func (o *OrdArrayLinear) AddCollection(word string) int {
	i := o.size

	fmt.Printf("Inserting '%s'...\n", word)
	fmt.Println("i:", i)

	for i > 0 && o.words[i-1] > word {
		fmt.Println("i:", i)
		fmt.Println("o.words[i]:", o.words[i])
		fmt.Println("o.words[i-1]:", o.words[i-1])
		o.words[i] = o.words[i-1]
		i--
	}

	o.words = append(o.words, "")
	copy(o.words[i+1:], o.words[i:])
	o.words[i] = word
	o.size++

	return SUCCESS
}

/*
 * SearchCollection searches for the string, given by the second parameter, in
 * the WordCollection given by the first parameter. This utilises a linear
 * search algorithm, and returns SUCCESS or FAILURE depending upon the outcome
 * of the search.
 */
func (o *OrdArrayLinear) SearchCollection(word string) int {
	for i := 0; i < o.size; i++ {
		if o.words[i] == word {
			return SUCCESS
		} else if o.words[i] > word {
			return FAILURE
		}
	}

	return FAILURE
}

/*
 * SizeCollection returns the number of words in the WordCollection given as
 * the parameter.
 */
func (o *OrdArrayLinear) SizeCollection() int {
	return o.size
}

/*
 * DisplayCollection prints the contents of the WordCollection given as the
 * parameter to standard output.
 */
func (o *OrdArrayLinear) DisplayCollection() {
	for i := 0; i < o.size; i++ {
		fmt.Printf("Element %d:\t%s\n", i, o.words[i])
	}
}
