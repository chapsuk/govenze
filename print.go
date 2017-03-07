package main

import "fmt"
import "strings"

type table struct {
	rows []row
}

type row struct {
	columns []string
}

// Col1 | Col2 | Col3 | Col4
// --- | --- | --- | ---
// row1 | row1 | row1 | row1
func prettyPrint(repo string, res []*result) {
	fmt.Printf("\nResult for %s\n\n", repo)
	t := &table{
		rows: make([]row, 4), // 1 - head, 2 - size, 3 - time, 4 - build
	}

	t.rows[0].columns = append(t.rows[0].columns, "Tool")
	t.rows[1].columns = append(t.rows[1].columns, "Size")
	t.rows[2].columns = append(t.rows[2].columns, "Time")
	t.rows[3].columns = append(t.rows[3].columns, "Builded")

	for _, r := range res {
		if r == nil {
			continue
		}
		t.rows[0].columns = append(t.rows[0].columns, r.manager)
		t.rows[1].columns = append(t.rows[1].columns, fmt.Sprintf("%.2fMb", r.size))
		t.rows[2].columns = append(t.rows[2].columns, fmt.Sprintf("%.2fs", r.elapsed.Seconds()))
		t.rows[3].columns = append(t.rows[3].columns, fmt.Sprintf("%t", r.build))
	}

	for i, r := range t.rows {
		for j, val := range r.columns {
			fmt.Printf("%s", val)
			if j != len(r.columns)-1 {
				fmt.Printf(" | ")
			}
		}
		fmt.Print("\n")
		if i == 0 {
			sep := strings.Repeat(" --- |", len(t.rows[i].columns))
			fmt.Print(sep[:len(sep)-2])
			fmt.Print("\n")
		}
	}

	fmt.Print("\n")
}
