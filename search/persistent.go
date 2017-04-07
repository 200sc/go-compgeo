package search

// Persistent types have a concept of time instants elapsing.
type Persistent interface {
	MinInstant() float64
	MaxInstant() float64
	SetInstant(float64)
}

// StaticPersistent types are Persistent types that are searchable as
// static types. At a given instant of time, a sub-static type can be
// returned from a StaticPersistent type.
type StaticPersistent interface {
	Persistent
	Static
	AtInstant(float64) Static
}

// DynamicPersistent types are StaticPersitent types that also allow
// modification as Dynamic types. AtInstant on a DynamicPersistent type
// will return a sub-dynamic type instead of a sub-static type, and
// being implicitly an extension on StaticPersistent, these can be converted
// back to StaticPersistent.
type DynamicPersistent interface {
	Persistent
	Dynamic
	AtInstant(float64) Dynamic
	ToStaticPersitent() StaticPersistent
}

// Persistable types are dynamic types, convertible to PersistentDynamic.
type Persistable interface {
	Dynamic
	ToPersistent() DynamicPersistent
}

// Why is there no PersistableStatic type? Because static types cannot be
// modified. There's no point in converting something which is static to
// a persistent type as it will only ever have the one instant which was
// the existing Static type.
//
// There is justification in converting a DynamicPersitent type to
// StaticPersistent, as you may be done making modifications to the time
// structure.
//
// Why isn't there a StaticPersistent -> DynamicPersistent function?
// See "Why isn't there a Static -> Dynamic function" in search.go