package compgeo

// TypeError is returned when some input to be
// read is improperly formatted for the expected type.
type TypeError struct{}

func (fte TypeError) Error() string {
	return "The input was not of the expected type, or was malformed"
}

// EmptyError is returned when some input to be read
// did not have any contents.
type EmptyError struct{}

func (ee EmptyError) Error() string {
	return "The input was empty"
}

// NotManifoldError is returned when it is detected
// that the input shape to ReadOFF was not possible
// Euclidean geometry.
type NotManifoldError struct{}

func (nme NotManifoldError) Error() string {
	return "The given shape was not manifold"
}

// BadEdgeError is returned from edge-processing functions
// if an edge is expected to have access to some field, or
// be initialized, when it does or is not. I.E. an edge has
// no twin for FullEdge.
type BadEdgeError struct{}

func (bee BadEdgeError) Error() string {
	return "The input edge was invalid"
}

// A RangeError represents when some query attempted
// on a structure falls out of the range of the structure's
// span.
type RangeError struct{}

func (re RangeError) Error() string {
	return "The query value was not in range of the structure"
}

// A BadDimensionError is returned when some function that takes
// an input dimension is given a dimension which is not defined
// for the query structure.
type BadDimensionError struct{}

func (bde BadDimensionError) Error() string {
	return "The query dimension does not exist in the query structure."
}

// An InsufficientDimensionsError is returned when some
// function takes n input dimensions and was given less than
// n dimensions to work with.
type InsufficientDimensionsError struct{}

func (ide InsufficientDimensionsError) Error() string {
	return "Not enough dimensions were supplied to the function"
}

// BadDCELError is returned when a query needs access to eleemnts
// of a DCEL that are not defined on a given DCEL. A good example
// would be that most queries would expect a DCEL to have at least
// three vertices.
type BadDCELError struct{}

func (bde BadDCELError) Error() string {
	return "The input DCEL was not valid"
}

// DivideByZero is returned by functions that attempt to Divide by zero.
type DivideByZero struct{}

func (dbz DivideByZero) Error() string {
	return "Division by zero. Default value to Infinity if unavoidable"
}
