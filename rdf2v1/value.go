package rdf2v1

import (
	"regexp"
	"strings"
)

func Str_ValStr(v string) ValueStr { return ValueStr{v} }

// interface for elements
type Value interface {
	Val() string
}

// string
type ValueStr struct {
	Val string
}

func Str(v string) ValueStr              { return ValueStr{v} }
func (v ValueStr) V() string             { return v.Val }
func (v ValueStr) Equal(w ValueStr) bool { return v.Val == w.Val }

// bool
type ValueBool struct {
	Val bool
}

func Bool(v bool) ValueBool { return ValueBool{v} }
func (v ValueBool) Value() string {
	if v.Val {
		return "true"
	}
	return "false"
}

type ValueCreator struct {
	val   string
	what  string
	name  string
	email string
}

// Get the original value of this ValueCreator
func (c ValueCreator) Val() string { return c.val }

// Get the `what` part from the format `what: name (email)`.
func (c ValueCreator) What() string { return c.what }

// Get the `name` part from the format `what: name (email)`
func (c ValueCreator) Name() string { return c.name }

// Get the `email` part from the format `what: name (email)`
func (c ValueCreator) Email() string { return c.email }

// parses and populates values of value creator
func (c *ValueCreator) SetValue(v string) {
	c.val = v
	RegexCreator := regexp.MustCompile("^([^:]*):([^\\(]*)(\\((.*)\\))?$")
	match := RegexCreator.FindStringSubmatch(v)
	if len(match) == 5 {
		c.what = strings.TrimSpace(match[1])
		c.name = strings.TrimSpace(match[2])
		c.email = strings.TrimSpace(match[4])
	}
}

// Create and populate a new ValueCreator.
func ValueCreatorNew(val string) ValueCreator {
	var valuecreator ValueCreator
	(&valuecreator).SetValue(val)
	return valuecreator
}

type ValueDate struct {
	val string
}

func (d ValueDate) V() string          { return d.val }
func (d *ValueDate) SetValue(v string) { d.val = v }

// New ValueDate.
func ValueDateNew(val string) ValueDate {
	var valuedate ValueDate
	(&valuedate).SetValue(val)
	return valuedate
}
