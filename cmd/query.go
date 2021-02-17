package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satoqz/lyr/query"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"get", "find", "search"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    queryExec,
}

func queryExec(_ *cobra.Command, args []string) error {

	q := query.New(strings.Join(args, " "))
	res, err := q.Search()
	if err != nil {
		return err
	}
	data := res.Collect()

	count := 0
	for i, v := range data {
		fmt.Printf("%d. %s\n", i+1, v.Name)
		count++
	}
	if count == 0 {
		return errors.New("No songs found. Exiting.\n")
	}
	fmt.Printf("Select a song to fetch lyrics of: [1-%d] >> ", count)

	var sel int

	for sel < 1 || sel > count || err != nil {
		_, err = fmt.Scanln(&sel)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("Fetching...")

	song := data[sel-1]

	lyrics, err := song.ScrapeLyrics()

	fmt.Print("\033[H\033[2J")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n\n%s\n", song.Name, lyrics)
	return nil
}
