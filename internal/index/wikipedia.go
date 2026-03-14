package index

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var indexPages = map[string]string{
	"sp500":     "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies",
	"nasdaq100": "https://en.wikipedia.org/wiki/Nasdaq-100",
	"dowjones":  "https://en.wikipedia.org/wiki/Dow_Jones_Industrial_Average",
}

// FetchAllConstituents fetches constituents from all three indices.
func FetchAllConstituents() (*Constituents, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	sp500, err := fetchAndParse(client, indexPages["sp500"], parseSP500)
	if err != nil {
		return nil, fmt.Errorf("S&P 500: %w", err)
	}

	nasdaq100, err := fetchAndParse(client, indexPages["nasdaq100"], parseNasdaq100)
	if err != nil {
		return nil, fmt.Errorf("NASDAQ 100: %w", err)
	}

	dowjones, err := fetchAndParse(client, indexPages["dowjones"], parseDowJones)
	if err != nil {
		return nil, fmt.Errorf("Dow Jones: %w", err)
	}

	merged := mergeSymbols(sp500, nasdaq100, dowjones)

	return &Constituents{
		UpdatedAt: time.Now().UTC(),
		SP500:     sp500,
		Nasdaq100: nasdaq100,
		DowJones:  dowjones,
		Merged:    merged,
	}, nil
}

func fetchAndParse(client *http.Client, rawURL string, parser func(io.Reader) ([]string, error)) ([]string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "earnings-cal-cli/1.0 (https://github.com/oinkywaddles/earnings-cal-cli)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, rawURL)
	}

	return parser(resp.Body)
}

// parseSymbolTable parses an HTML page to find a table with a "Symbol" or "Ticker"
// column header, then extracts all valid symbols from that column.
func parseSymbolTable(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	tables := findElements(doc, "table")
	if len(tables) == 0 {
		return nil, fmt.Errorf("no tables found")
	}

	for _, t := range tables {
		rows := findElements(t, "tr")
		if len(rows) < 2 {
			continue
		}

		// Find the symbol column index from the header row
		colIndex := -1
		headers := findElements(rows[0], "th")
		for i, h := range headers {
			text := strings.ToLower(strings.TrimSpace(textContent(h)))
			if text == "symbol" || text == "ticker" || text == "ticker symbol" {
				colIndex = i
				break
			}
		}
		if colIndex == -1 {
			continue
		}

		// Extract symbols from data rows
		var symbols []string
		for _, row := range rows[1:] {
			cells := findCells(row)
			if colIndex >= len(cells) {
				continue
			}
			text := strings.TrimSpace(textContent(cells[colIndex]))
			if text != "" && isValidSymbol(text) {
				symbols = append(symbols, text)
			}
		}

		if len(symbols) > 0 {
			return symbols, nil
		}
	}

	return nil, fmt.Errorf("no table with Symbol/Ticker column found")
}

// parseSP500 parses the S&P 500 wiki page.
func parseSP500(r io.Reader) ([]string, error) {
	return parseSymbolTable(r)
}

// parseNasdaq100 parses the NASDAQ 100 wiki page.
func parseNasdaq100(r io.Reader) ([]string, error) {
	return parseSymbolTable(r)
}

// parseDowJones parses the Dow Jones wiki page.
func parseDowJones(r io.Reader) ([]string, error) {
	return parseSymbolTable(r)
}

// findCells returns all direct th and td children of a tr (in order).
func findCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "th" || n.Data == "td") {
			cells = append(cells, n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(row)
	return cells
}

func findElements(n *html.Node, tag string) []*html.Node {
	var result []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			result = append(result, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return result
}

func textContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(textContent(c))
	}
	return sb.String()
}

func isValidSymbol(s string) bool {
	if len(s) == 0 || len(s) > 5 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

func mergeSymbols(lists ...[]string) []string {
	seen := make(map[string]bool)
	var merged []string
	for _, list := range lists {
		for _, s := range list {
			if !seen[s] {
				seen[s] = true
				merged = append(merged, s)
			}
		}
	}
	return merged
}
