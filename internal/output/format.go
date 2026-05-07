package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type Format string

const (
	FormatJSON   Format = "json"
	FormatTable  Format = "table"
	FormatPretty Format = "pretty"
)

func FormatOutput(w io.Writer, data interface{}, format Format) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case FormatPretty:
		return formatPretty(w, data)
	case FormatTable:
		return formatTable(w, data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintTable(w io.Writer, rows []map[string]interface{}) {
	if len(rows) == 0 {
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	headers := []string{}
	for k := range rows[0] {
		headers = append(headers, k)
	}

	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	for _, row := range rows {
		for i, h := range headers {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			val := row[h]
			if val != nil {
				fmt.Fprint(tw, fmt.Sprintf("%v", val))
			}
		}
		fmt.Fprintln(tw)
	}

	tw.Flush()
}

func formatPretty(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("  ", "  ")
	return enc.Encode(data)
}

func formatTable(w io.Writer, data interface{}) error {
	switch v := data.(type) {
	case []map[string]interface{}:
		PrintTable(w, v)
	case map[string]interface{}:
		PrintTable(w, []map[string]interface{}{v})
	default:
		return FormatOutput(w, data, FormatJSON)
	}
	return nil
}

func ErrWithHint(exitCode int, code, msg, hint string) error {
	return fmt.Errorf("%s: %s\nHint: %s", code, msg, hint)
}

func PrintJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func Print(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

func Println(w io.Writer, args ...interface{}) {
	fmt.Fprintln(w, args...)
}

func ExitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
