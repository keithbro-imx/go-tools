//go:build go1.18 && ignore

package pkg

type S1[T any] struct {
	// flag, 'any' is too permissive
	F T `json:",string"` // want `the JSON string option`
}

type S2[T int | string] struct {
	// don't flag, all types in T are okay
	F T `json:",string"`
}

type S3[T int | complex128] struct {
	// flag, can't use ,string on complex128
	F T `json:",string"` // want `the JSON string option`
}

type S4[T int | string] struct {
	// don't flag, pointers to stringable types are also stringable
	F *T `json:",string"`
}

type S5[T int | string, PT *T] struct {
	// don't flag, pointers to stringable types are also stringable
	F PT `json:",string"`
}

type S6[T int | complex128] struct {
	// flag, pointers to non-stringable types aren't stringable, either
	F *T `json:",string"` // want `the JSON string option`
}

type S7[T int | complex128, PT *T] struct {
	// flag, pointers to non-stringable types aren't stringable, either
	F PT `json:",string"` // want `the JSON string option`
}

type S8[T int, PT *T | complex128] struct {
	// do flag, variation of S7
	F PT `json:",string"` // want `the JSON string option`
}

type S9[T int | *bool, PT *T | float64, PPT *PT | string] struct {
	// don't flag, multiple levels of pointers are fine
	F PPT `json:",string"`
}

type S10[T1 *T2, T2 *T1] struct {
	// flag? don't flag? who knows. just don't get stuck in an infinite loop
	F T1 `json:",string"`
}
