package main

import (
	"regexp"
	"strings"
	"time"
)
// interface for elements
type Value interface {
	Val() string
}

// string
type ValueStr struct {
	Val  string
}

func Str(v string) ValueStr { return ValueStr{v} }
func (v ValueStr) V() string { return v.Val }
func (v ValueStr) Equal(w ValueStr) bool { return v.Val == w.Val }


// bool
type ValueBool struct {
	Val  bool
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

// Set the value of this ValueCreator. It parses the format `what: name (email)`
// and populates the relevant fields.
func (c *ValueCreator) SetValue(v string) {
	c.val = v
	CreatorRegex := regexp.MustCompile("^([^:]*):([^\\(]*)(\\((.*)\\))?$")
	match := CreatorRegex.FindStringSubmatch(v)
	if len(match) == 5 {
		c.what = strings.TrimSpace(match[1])
		c.name = strings.TrimSpace(match[2])
		c.email = strings.TrimSpace(match[4])
	}
}
// Create and populate a new ValueCreator.
func CreateValue(val string, m *Meta) ValueCreator {
	value := ValueCreator{Meta: m}
	(&value).SetValue(val)
	return value
}
type ValueDate struct {
	val  string
	time *time.Time
}


// Set the value of this ValueDate and parse the date format.
func SetDateValue(v string){
	 = v

}
func SetDateTime(t string){
	time, err := time.Parse(time.RFC3339, v)
	if err == nil {
		d.time = &time
	}
}

// Create and populate a new ValueDate.
func NewValueDate(val string) ValueDate {
	valuedate := ValueDate{}
	(&valuedate).SetValue(val)
	return valuedate
}
