package main

import (
	"regexp"
	"strings"
	"time"
)

//Store data of form 'what: name (email)'
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

type ValueDate struct {
	val  string
	time *time.Time
}

// Get the original value of this ValueDate.
func (d ValueDate) Val() string { return d.val }

// Get the *time.Time pointer parsed form the value.
func (d ValueDate) Time() *time.Time { return d.time }

// Set the value of this ValueDate and parse the date format.
func (d *ValueDate) SetValue(v string) {
	d.val = v
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
