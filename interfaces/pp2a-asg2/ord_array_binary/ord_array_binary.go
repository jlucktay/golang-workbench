// ord_array_binary is an ordered slice with binary search
package ord_array_binary

const (
	SUCCESS = iota
	FAILURE
)

const WCSIZE = 250000

type OrdArrayBinary struct {
	words []string
	size  uint
}

/*
 * MakeCollection initialises the WordCollection whose pointer it is given as
 * a parameter, up to a size defined by WCSIZE, with NULL entries. It returns
 * SUCCESS upon successful completion.
 */
func (o *OrdArrayBinary) MakeCollection() int {
	o.words = make([]string, 0, WCSIZE)
	o.size = 0

	return SUCCESS
}

/*
 * FreeCollection frees the memory dynamically allocated to the WordCollection
 * parameter.
 */
func (o *OrdArrayBinary) FreeCollection() {
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
func (o *OrdArrayBinary) AddCollection(word string) int {
	return FAILURE
}

/*
 * SearchCollection searches for the string, given by the second parameter, in
 * the WordCollection given by the first parameter. This utilises a binary
 * search algorithm, and returns SUCCESS or FAILURE depending upon the outcome
 * of the search.
 */
func (o *OrdArrayBinary) SearchCollection(string) int {
	return FAILURE
}

/*
 * SizeCollection returns the number of words in the WordCollection given as
 * the parameter.
 */
func (o *OrdArrayBinary) SizeCollection() int {
	return FAILURE
}

/*
 * DisplayCollection prints the contents of the WordCollection given as the
 * parameter to standard output.
 */
func (o *OrdArrayBinary) DisplayCollection() {}

//  /*
//   * StrCmpWrap is a simple wrapper function for the strcmp function (defined in
//   * string.h). This allows it to be compatible with the bsearch function
//   * (defined in stdlib.h) that is utilised by SearchCollection.
//   */
//  int StrCmpWrap(const void *, const void *);
