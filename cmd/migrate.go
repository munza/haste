package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(help())
		os.Exit(0)
	}

	switch os.Args[1] {
	case "up":
		up()
	case "down":
		down()
	case "redo":
		redo()
	case "status":
		status()
	case "new":
		new()
	default:
		fmt.Printf("Unknows command '%s' for cmd/migrate.go\n", os.Args[1])
	}
}

func up() {
	err := applyMigrations(migrate.Up, false, 0)
	if err != nil {
		panic(err)
	}
}

func down() {
	err := applyMigrations(migrate.Down, false, 1)
	if err != nil {
		panic(err)
	}
}

func redo() {
	db, dialect, err := getDBConnection()
	if err != nil {
		panic(err)
	}

	source, _ := getSources()

	migrations, _, _ := migrate.PlanMigration(db, dialect, source, migrate.Down, 1)
	if len(migrations) == 0 {
		fmt.Println("Nothing to do!")
	}

	_, err = migrate.ExecMax(db, dialect, source, migrate.Down, 1)
	if err != nil {
		fmt.Printf("Migration (down) failed: %s\n", err)
	}

	_, err = migrate.ExecMax(db, dialect, source, migrate.Up, 1)
	if err != nil {
		fmt.Printf("Migration (up) failed: %s\n", err)
	}

	fmt.Printf("Reapplied migration %s.\n", migrations[0].Id)
}

func status() {

	db, dialect, _ := getDBConnection()

	records, err := migrate.GetMigrationRecords(db, dialect)
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]*statusRow)

	_, migrations := getSources()
	for _, m := range migrations {
		rows[m.Id] = &statusRow{
			ID:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range migrations {
		if rows[m.Id].Migrated {
			table.Append([]string{
				m.Id,
				rows[m.Id].AppliedAt.String(),
			})
		} else {
			table.Append([]string{
				m.Id,
				"no",
			})
		}
	}

	table.Render()

	os.Exit(0)
}

func new() {
	if len(os.Args) < 3 {
		fmt.Println("A name for the migration is needed")
		os.Exit(0)
	}

	if err := createMigration(os.Args[2]); err != nil {
		panic(err)
	}
}

func applyMigrations(dir migrate.MigrationDirection, dryrun bool, limit int) error {
	db, dialect, err := getDBConnection()
	if err != nil {
		return err
	}

	source, _ := getSources()

	if dryrun {
		migrations, _, err := migrate.PlanMigration(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Cannot plan migration: %s", err)
		}

		for _, m := range migrations {
			printMigration(m, dir)
		}
	} else {
		n, err := migrate.ExecMax(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			fmt.Println("Applied 1 migration")
		} else {
			fmt.Printf("Applied %d migrations\n", n)
		}
	}

	return nil
}

func printMigration(m *migrate.PlannedMigration, dir migrate.MigrationDirection) {
	if dir == migrate.Up {
		fmt.Printf("==> Would apply migration %s (up)\n", m.Id)
		for _, q := range m.Up {
			fmt.Println(q)
		}
	} else if dir == migrate.Down {
		fmt.Printf("==> Would apply migration %s (down)\n", m.Id)
		for _, q := range m.Down {
			fmt.Println(q)
		}
	} else {
		panic("Not reached")
	}
}

func getSources() (migrate.FileMigrationSource, []*migrate.Migration) {
	source := migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	migrations, err := source.FindMigrations()
	if err != nil {
		panic(err)
	}

	return source, migrations
}

func getDBConnection() (*sql.DB, string, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:33061)/hastedb?parseTime=true")
	if err != nil {
		return nil, "", fmt.Errorf("Cannot connect to database: %s", err)
	}

	return db, "mysql", nil
}

func createMigration(name string) error {
	if _, err := os.Stat("db/migrations"); os.IsNotExist(err) {
		panic(err)
	}

	var templateContent = `-- +migrate Up

-- +migrate Down
`
	var tpl *template.Template

	tpl = template.Must(template.New("new_migration").Parse(templateContent))

	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join("db/migrations", fileName)

	f, err := os.Create(pathName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tpl.Execute(f, nil); err != nil {
		panic(err)
	}

	fmt.Printf("Created migration %s\n", pathName)
	return nil
}

func help() string {
	return `Usage: go run cmd/migrate.go COMMAND

Available Commands:
    up          Migrates the database to the most recent version available.
    down        Undo a database migration.
    redo        Reapply the last migration.
    status      Show migration status.
    new         Create a new a database migration.
`
}

type statusRow struct {
	ID        string
	Migrated  bool
	AppliedAt time.Time
}
