package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
)

type jsonListOutput struct {
	Range    jsonRange          `json:"range"`
	Filter   string             `json:"filter"`
	Count    int                `json:"count"`
	Earnings []finnhub.Earning  `json:"earnings"`
}

type jsonRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// PrintListJSON outputs the list in JSON format.
func PrintListJSON(w io.Writer, earnings []finnhub.Earning, from, to, filterDesc string) {
	out := jsonListOutput{
		Range:    jsonRange{From: from, To: to},
		Filter:   filterDesc,
		Count:    len(earnings),
		Earnings: earnings,
	}

	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Fprintln(w, string(data))
}

// PrintDetailJSON outputs detail view in JSON format.
func PrintDetailJSON(w io.Writer, earnings []finnhub.Earning) {
	data, _ := json.MarshalIndent(earnings, "", "  ")
	fmt.Fprintln(w, string(data))
}
