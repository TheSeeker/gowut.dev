// Copyright 2013 Andras Belicza. All rights reserved.

// ListBox component interface and implementation.

package gwu

import (
	"net/http"
	"strconv"
	"strings"
)

// ListBox interface defines a component which allows selecting one or multiple values
// from a predefined list.
// 
// Suggested event type to handle changes: ETYPE_CHANGE
// 
// Default style class: "gwu-ListBox"
type ListBox interface {
	// ListBox is a component
	Comp

	// ListBox can be enabled/disabled.
	HasEnabled

	// Multi tells if multiple selections are allowed.
	Multi() bool

	// SetMulti sets whether multiple selections are allowed.
	SetMulti(multi bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this ListBox a dropdown list (if multi is false!).
	// Note that if rows is greater than 1, most browsers enforce a visual minimum size
	// (about 4 rows) even if rows is less than that.
	SetRows(rows int)

	// SelectedValue retruns the first selected value.
	// Empty string is returned if nothing is selected.
	SelectedValue() string

	// SelectedValues retruns all the selected values.
	SelectedValues() []string

	// Selected tells if the value at index i is selected.
	Selected(i int) bool

	// SelectedIndices returns a slice of the indices of the selected values.
	SelectedIndices() []int

	// SetSelected sets the selection state of the value at index i.
	SetSelected(i int, selected bool)

	// SetSelectedIndices sets the (only) selected values.
	// Only values will be selected that are contained in the specified indices slice.
	SetSelectedIndices(indices []int)

	// ClearSelected deselects all values.
	ClearSelected()
}

// ListBox implementation.
type listBoxImpl struct {
	compImpl       // Component implementation 
	hasEnabledImpl // Has enabled implementation

	values   []string // Values to choose from
	multi    bool     // Allow multiple selection
	selected []bool   // Array of selection state of the values
	rows     int      // Number of displayed rows
}

// NewListBox creates a new ListBox.
func NewListBox(values []string) ListBox {
	c := &listBoxImpl{newCompImpl("getSelectedIndices(this)"), newHasEnabledImpl(), values, false, make([]bool, len(values)), 1}
	c.AddSyncOnETypes(ETYPE_CHANGE)
	c.Style().AddClass("gwu-ListBox")
	return c
}

func (c *listBoxImpl) Multi() bool {
	return c.multi
}

func (c *listBoxImpl) SetMulti(multi bool) {
	c.multi = multi
}

func (c *listBoxImpl) Rows() int {
	return c.rows
}

func (c *listBoxImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *listBoxImpl) SelectedValue() string {
	for i, s := range c.selected {
		if s {
			return c.values[i]
		}
	}

	return ""
}

func (c *listBoxImpl) SelectedValues() (sv []string) {
	for i, s := range c.selected {
		if s {
			sv = append(sv, c.values[i])
		}
	}
	return
}

func (c *listBoxImpl) Selected(i int) bool {
	return c.selected[i]
}

func (c *listBoxImpl) SelectedIndices() (si []int) {
	for i, s := range c.selected {
		if s {
			si = append(si, i)
		}
	}
	return
}

func (c *listBoxImpl) SetSelected(i int, selected bool) {
	c.selected[i] = selected
}

func (c *listBoxImpl) SetSelectedIndices(indices []int) {
	// First clear selected slice
	for i, _ := range c.selected {
		c.selected[i] = false
	}

	// And now select that needs to be selected
	for _, idx := range indices {
		c.selected[idx] = true
	}
}

func (c *listBoxImpl) ClearSelected() {
	for i, _ := range c.selected {
		c.selected[i] = false
	}
}

func (c *listBoxImpl) preprocessEvent(event Event, r *http.Request) {
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) == 0 {
		return
	}

	// Set selected indices
	c.ClearSelected()
	for _, sidx := range strings.Split(value, ",") {
		if idx, err := strconv.Atoi(sidx); err == nil {
			c.selected[idx] = true
		}
	}
}

func (c *listBoxImpl) Render(w writer) {
	w.Writes("<select")
	if c.multi {
		w.Writes(" multiple=\"multiple\"")
	}
	w.WriteAttr("size", strconv.Itoa(c.rows))
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	for i, value := range c.values {
		if c.selected[i] {
			w.Writes("<option selected=\"selected\">")
		} else {
			w.Writes("<option>")
		}
		w.Writees(value)
		w.Writes("</option>")
	}

	w.Writes("</select>")
}
