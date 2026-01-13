package main_test

import (
	"reflect"
	"testing"

	"github.com/hreshchyshynt/gator"
)

func TestNewCommand(t *testing.T) {
	// - `gator browse --sort=oldest` (see historical posts first)
	// - `gator browse --sort=title` (alphabetical by title)
	// - `gator browse --sort=feed` (group by feed source)
	//
	// ### Filtering Flags
	//
	// Filtering flags let you narrow down what you see:
	// - `gator browse --limit=10` (show only 10 posts instead of 2)
	// - `gator browse --feed="Boot.dev Blog"` (only posts from specific feed)
	// - `gator browse --since="2024-01-01"` (posts after a date)
	// - `gator browse --unread` (if you track read status)
	//
	// ### Pagination
	//
	// Pagination handles large result sets:
	// - `gator browse --page=2` (show next page of results)
	// - `gator browse --limit=20 --offset=40` (skip first 40, show next 20)

	tests := []struct {
		// Named input parameters for target function.
		name string
		args []string
		want main.Command
	}{
		{
			name: "browse",
			args: []string{"4"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Value: "4"}},
			},
		},

		// ---- Sorting flags ----
		{
			name: "browse",
			args: []string{"--sort=oldest"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "sort", Value: "oldest"}},
			},
		},
		{
			name: "browse",
			args: []string{"--sort=title"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "sort", Value: "title"}},
			},
		},
		{
			name: "browse",
			args: []string{"--sort=feed"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "sort", Value: "feed"}},
			},
		},

		// ---- Filtering flags ----
		{
			name: "browse",
			args: []string{"--limit=10"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "limit", Value: "10"}},
			},
		},
		{
			name: "browse",
			args: []string{`--feed="Boot.dev Blog"`},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "feed", Value: `"Boot.dev Blog"`}},
			},
		},
		{
			name: "browse",
			args: []string{`--since="2024-01-01"`},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "since", Value: `"2024-01-01"`}},
			},
		},
		{
			name: "browse",
			args: []string{"--unread"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Value: "unread"}},
			},
		},

		// ---- Pagination ----
		{
			name: "browse",
			args: []string{"--page=2"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{{Name: "page", Value: "2"}},
			},
		},
		{
			name: "browse",
			args: []string{"--limit=20", "--offset=40"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{
					{Name: "limit", Value: "20"},
					{Name: "offset", Value: "40"},
				},
			},
		},

		// ---- Realistic combos ----
		{
			name: "browse",
			args: []string{"--sort=oldest", "--limit=10"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{
					{Name: "sort", Value: "oldest"},
					{Name: "limit", Value: "10"},
				},
			},
		},
		{
			name: "browse",
			args: []string{`--feed="Boot.dev Blog"`, `--since="2024-01-01"`},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{
					{Name: "feed", Value: `"Boot.dev Blog"`},
					{Name: "since", Value: `"2024-01-01"`},
				},
			},
		},
		{
			name: "browse",
			args: []string{"--unread", "--page=2", "--limit=20"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{
					{Value: "unread"},
					{Name: "page", Value: "2"},
					{Name: "limit", Value: "20"},
				},
			},
		},
		{
			name: "browse",
			args: []string{"--sort=feed", `--feed="Boot.dev Blog"`, "--limit=20", "--offset=40"},
			want: main.Command{
				Name: "browse",
				Args: []main.Argument{
					{Name: "sort", Value: "feed"},
					{Name: "feed", Value: `"Boot.dev Blog"`},
					{Name: "limit", Value: "20"},
					{Name: "offset", Value: "40"},
				},
			},
		}}
	for _, tt := range tests {
		got := main.NewCommand(tt.name, tt.args)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewCommand() = %v, want %v", got, tt.want)
		}
	}
}
