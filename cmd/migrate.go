package main

import (
	"database/sql"
	"fmt"
	"haste/config"
	"html/template"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "up":
		up()
	case "down":
		down()
	case "redo":
		redo()
	case "reset":
		reset()
	case "refresh":
		reset()
		up()
	case "status":
		status()
	case "new":
		new()
	case "help":
		help()
	default:
		fmt.Printf("Unknows command '%s' for cmd/migrate.go\n", os.Args[1])
	}

	os.Exit(0)
}

func up() {
	err := applyMigrations(migrate.Up, isDryrun(), getLimit(0))
	if err != nil {
		panic(err)
	}
}

func down() {
	err := applyMigrations(migrate.Down, isDryrun(), getLimit(1))
	if err != nil {
		panic(err)
	}
}

func reset() {
	err := applyMigrations(migrate.Down, isDryrun(), 0)
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

	if isDryrun() {
		printMigration(migrations[0], migrate.Down)
		printMigration(migrations[0], migrate.Up)
	} else {
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

		action := "Applied"
		if dir > 0 {
			action = "Rollback"
		}

		if n == 1 {
			fmt.Printf("%s 1 migration\n", action)
		} else {
			fmt.Printf("%s %d migrations\n", action, n)
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
		Dir: config.Database().MigrationPath,
	}

	migrations, err := source.FindMigrations()
	if err != nil {
		panic(err)
	}

	return source, migrations
}

func getDBConnection() (*sql.DB, string, error) {
	connSrc := config.Database().Username + ":" + config.Database().Password + "@tcp(" + config.Database().Host + ":" + strconv.Itoa(config.Database().Port) + ")/" + config.Database().Name + "?parseTime=true"

	db, err := sql.Open(config.Database().Driver, connSrc)
	if err != nil {
		return nil, "", fmt.Errorf("Cannot connect to database: %s", err)
	}

	return db, "mysql", nil
}

func createMigration(name string) error {
	if _, err := os.Stat(config.Database().MigrationPath); os.IsNotExist(err) {
		panic(err)
	}

	var templateContent = `-- +migrate Up

-- +migrate Down
`
	var tpl *template.Template

	tpl = template.Must(template.New("new_migration").Parse(templateContent))

	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join(config.Database().MigrationPath, fileName)

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

func isDryrun() bool {
	var dryrun bool

	for _, value := range os.Args {
		if value == "-dryrun" || value == "--dryrun" {
			dryrun = true
		}
	}

	return dryrun
}

func getLimit(defaultval int) int {
	limit := defaultval

	for i, value := range os.Args {
		if value == "-limit" || value == "--limit" {
			if l, err := strconv.Atoi(os.Args[i+1]); err == nil {
				limit = l
			}
		}
	}

	return limit
}

func help() {
	fmt.Println(`Usage: go run cmd/migrate.go COMMAND [Options]

Available Commands:
    up [N=0]        Migrates the database to the most recent version available.
    down            Undo a database migration.
    redo            Reapply the last migration.
    refresh         Reapply all migrations.
    status          Show migration status.
    new             Create a new a database migration.
    help            Show the usage instruction.

Avaliable Options:

    -dryrun         Don't apply migrations, just print them.
    -limit [N=0]    Limit the number of migrations (0 = unlimited).
`)
}

type statusRow struct {
	ID        string
	Migrated  bool
	AppliedAt time.Time
}
