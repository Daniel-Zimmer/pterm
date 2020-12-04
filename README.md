# pterm

Go package for printing structs as tables in the terminal.

---

pterm uses reflection to get struct field names and tags to define formatting behaviour.

## Example

```go
package main

import (
	"fmt"
	"github.com/Daniel-Zimmer/pterm"
	"math/rand"
	"time"
)

type Example struct {
	ID   string `json:"id" pterm:"2, right"`
	Name string `json:"name" pterm:"2, elastic"`
	Size string `json:"size" pterm:"left"`
}

func main() {
	sizes := []string{"BIG", "MEDIUM", "SMALL", "TINY"}
	names := []string{"Albert", "John", "Richard", "Mary", "Cleo", "Bernard", "Zimmer"}

	examples := make([]Example, 0)

	rand.Seed(time.Now().Unix())
	for i := 0; i < 15; i++ {
		examples = append(examples, Example{
			ID:   fmt.Sprintf("%x", rand.Int()),
			Name: names[rand.Intn(len(names))],
			Size: sizes[rand.Intn(len(sizes))],
		},
		)
	}

	pterm.PrintTable(examples)

}

```

## Tags

### "left"
Aligns to the **left**.


### "right"
Aligns to the **right**.


### any number
How much **padding** to leave before the next column.

Example: "2"

### "ignore"
**Ignores** the field.

### "elastic"
Trims the strings in this column that would not fit in the terminal.
Only one field can be elastic.

(99% sure this does not work on Windows)

Example: if one of the entries in the table is: **"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."** and the terminal is not wide enough to fit the entire string, the displayed value will be something like: **"Lorem ipsum dolor sit..."**.

### any name:

**Name** to use as **column header** instead of the name of the struct field.
If you want to use any of the above keywords as struct names just add a "$" in front of the name.

Example: "$left"

If you want the column header to start with a "$" just use "$$" instead.
