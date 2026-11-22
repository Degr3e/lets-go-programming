package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var departments int
	if _, err := fmt.Fscan(in, &departments); err != nil {
		fmt.Fprintln(os.Stderr, "reading departments:", err)

		return
	}

	for range departments {
		var employees int
		if _, err := fmt.Fscan(in, &employees); err != nil {
			fmt.Fprintln(os.Stderr, "reading employees:", err)

			return
		}

		low, high := 15, 30

		for range employees {
			var (
				sign string
				temp int
			)

			if _, err := fmt.Fscan(in, &sign, &temp); err != nil {
				fmt.Fprintln(os.Stderr, "reading preference:", err)

				return
			}

			switch sign {
			case ">=":
				low = max(low, temp)
			case "<=":
				high = min(high, temp)
			default:
				fmt.Fprintln(os.Stderr, "invalid sign:", sign)

				return
			}

			if low > high {
				fmt.Println(-1)
			} else {
				fmt.Println(low)
			}
		}
	}
}
