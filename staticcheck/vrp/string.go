package vrp

import (
	"fmt"
	"log"
	"math/big"

	"honnef.co/go/ssa"
)

// s1 + s2
// s[:]
// len(s)
// "" -> len("")
// s1 == s2
// len(s1) <cmp> x

type StringInterval struct {
	Length IntInterval
}

func (s StringInterval) Union(other Range) Range {
	i, ok := other.(StringInterval)
	if !ok {
		i = StringInterval{EmptyIntInterval}
	}
	if s.Length.Empty() || !s.Length.IsKnown() {
		return i
	}
	if i.Length.Empty() || !i.Length.IsKnown() {
		return s
	}
	return StringInterval{
		Length: s.Length.Union(i.Length).(IntInterval),
	}
}

func (s StringInterval) String() string {
	return s.Length.String()
}

func (s StringInterval) IsKnown() bool {
	return s.Length.IsKnown()
}

type StringSliceConstraint struct {
	aConstraint
	X     ssa.Value
	Lower ssa.Value
	Upper ssa.Value
}

func (c *StringSliceConstraint) String() string {
	var lname, uname string
	if c.Lower != nil {
		lname = c.Lower.Name()
	}
	if c.Upper != nil {
		uname = c.Upper.Name()
	}
	return fmt.Sprintf("%s[%s:%s]", c.X.Name(), lname, uname)
}

func (c *StringSliceConstraint) Eval(g *Graph) Range {
	lr := NewIntInterval(NewZ(&big.Int{}), NewZ(&big.Int{}))
	if c.Lower != nil {
		lr = g.Range(c.Lower).(IntInterval)
	}
	ur := g.Range(c.X).(StringInterval).Length
	if c.Upper != nil {
		ur = g.Range(c.Upper).(IntInterval)
	}
	if !lr.IsKnown() || !ur.IsKnown() {
		return StringInterval{}
	}

	ls := []Z{
		ur.lower.Sub(lr.lower),
		ur.upper.Sub(lr.lower),
		ur.lower.Sub(lr.upper),
		ur.upper.Sub(lr.upper),
	}
	// TODO(dh): if we don't truncate lengths to 0 we might be able to
	// easily detect slices with high < low. we'd need to treat -∞
	// specially, though.
	for i, l := range ls {
		if l.Sign() == -1 {
			ls[i] = NewZ(&big.Int{})
		}
	}

	return StringInterval{
		Length: NewIntInterval(MinZ(ls...), MaxZ(ls...)),
	}
}

func (c *StringSliceConstraint) Operands() []ssa.Value {
	vs := []ssa.Value{c.X}
	if c.Lower != nil {
		vs = append(vs, c.Lower)
	}
	if c.Upper != nil {
		vs = append(vs, c.Upper)
	}
	return vs
}

type StringIntersectionConstraint struct {
	aConstraint
	X ssa.Value
	I IntInterval
}

func (c *StringIntersectionConstraint) Operands() []ssa.Value {
	return []ssa.Value{c.X}
}

func (c *StringIntersectionConstraint) Eval(g *Graph) Range {
	log.Println(c.X)
	xi := g.Range(c.X).(StringInterval)
	if !xi.IsKnown() {
		return c.I
	}
	return StringInterval{
		Length: xi.Length.Intersection(c.I),
	}
}

func (c *StringIntersectionConstraint) String() string {
	return fmt.Sprintf("%s = %s.%t ⊓ %s", c.Y().Name(), c.X.Name(), c.Y().(*ssa.Sigma).Branch, c.I)
}
