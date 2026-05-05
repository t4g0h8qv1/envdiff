package export

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", Status: diff.StatusMissing, LeftValue: "localhost", RightValue: ""},
	{Key: "API_KEY", Status: diff.StatusConflict, LeftValue: "abc", RightValue: "xyz"},
	{Key: "PORT", Status: diff.StatusMissing, LeftValue: "", RightValue: "8080"},
}

func TestWrite_CSV(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, sampleResults, FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,status,left_value,right_value") {
		t.Error("missing CSV header")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("missing DB_HOST row")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("missing API_KEY row")
	}
}

func TestWrite_Markdown(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, sampleResults, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Status") {
		t.Error("missing markdown header")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("missing DB_HOST row")
	}
	if !strings.Contains(out, "---") {
		t.Error("missing markdown separator")
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, sampleResults, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"Key"`) && !strings.Contains(out, `"key"`) {
		t.Error("expected JSON output with key field")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("missing DB_HOST in JSON output")
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, sampleResults, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWrite_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, []diff.Result{}, FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,status") {
		t.Error("expected header even for empty results")
	}
}
