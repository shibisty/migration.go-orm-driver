package ormstore_test

import (
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
	"time"

	ormstore "migration-orm"
	"migration/cli"
)

// The template makes valid Go for the three kinds of names.
func TestTemplate(t *testing.T) {
	dir := t.TempDir()
	opts := cli.Options{Dir: dir, Package: "migrations", Template: ormstore.Template, Now: func() time.Time { return time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC) }}
	for name, want := range map[string]string{
		"create_posts_table":        `db.Schema.Create(ctx, "posts"`,
		"add_avatar_to_users_table": `db.Schema.Table(ctx, "users"`,
		"backfill_slugs":            "return nil",
	} {
		path, err := cli.Make(opts, name)
		if err != nil {
			t.Fatal(name, err)
		}
		src, _ := os.ReadFile(path)
		if _, err := parser.ParseFile(token.NewFileSet(), path, src, parser.AllErrors); err != nil {
			t.Fatalf("%s: %v\n%s", name, err, src)
		}
		if !strings.Contains(string(src), want) || !strings.Contains(string(src), `Name: "2026_11_09_000000_`+name+`"`) {
			t.Errorf("%s:\n%s", name, src)
		}
		if name == "backfill_slugs" && strings.Contains(string(src), `"orm/schema"`) {
			t.Errorf("unused import in %s", name)
		}
	}
}
