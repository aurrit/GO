package database

import (
    "database/sql"
    "os"
    "path/filepath"
    "sort"
)

func Migrate(db *sql.DB, dir string) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY)`)
    if err != nil {
        return err
    }

    files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
    if err != nil {
        return err
    }
    sort.Strings(files)

    for _, f := range files {
        name := filepath.Base(f)

        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE name = ?", name).Scan(&count)
        if err != nil {
            return err
        }
        if count > 0 {
            continue
        }

        content, err := os.ReadFile(f)
        if err != nil {
            return err
        }
        if _, err := db.Exec(string(content)); err != nil {
            return err
        }
        if _, err := db.Exec("INSERT INTO schema_migrations (name) VALUES (?)", name); err != nil {
            return err
        }
    }
    return nil
}