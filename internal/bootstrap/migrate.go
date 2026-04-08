package bootstrap

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

func splitMigrationStatements(sql string) []string {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil
	}
	var out []string
	var b strings.Builder
	inSingleQuote := false
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inSingleQuote {
			b.WriteByte(c)
			if c == '\'' {
				if i+1 < len(sql) && sql[i+1] == '\'' {
					b.WriteByte(sql[i+1])
					i++
					continue
				}
				inSingleQuote = false
			}
			continue
		}
		if c == '\'' {
			inSingleQuote = true
			b.WriteByte(c)
			continue
		}
		if c == ';' {
			s := strings.TrimSpace(b.String())
			if s != "" {
				out = append(out, s)
			}
			b.Reset()
			continue
		}
		b.WriteByte(c)
	}
	s := strings.TrimSpace(b.String())
	if s != "" {
		out = append(out, s)
	}
	return out
}

type schemaMigration struct {
	Version   string    `gorm:"column:version;primaryKey;type:text"`
	AppliedAt time.Time `gorm:"column:applied_at;autoCreateTime"`
}

func (schemaMigration) TableName() string {
	return "schema_migrations"
}

func RunMigrations(ctx context.Context, db *gorm.DB, fs embed.FS) error {
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return fmt.Errorf("инициализация таблицы миграций: %w", err)
	}

	migrationsDir := "migrations"
	entries, err := fs.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("чтение каталога миграций: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		version := name
		path := migrationsDir + "/" + name

		applied, err := isMigrationApplied(ctx, db, version)
		if err != nil {
			return fmt.Errorf("проверка миграции %s: %w", version, err)
		}
		if applied {
			continue
		}

		content, err := fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("чтение миграции %s: %w", version, err)
		}
		sql := strings.TrimSpace(string(content))
		if sql == "" {
			if err := markMigrationApplied(ctx, db, version); err != nil {
				return fmt.Errorf("запись версии %s: %w", version, err)
			}
			continue
		}

		stmts := splitMigrationStatements(sql)
		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, stmt := range stmts {
				if err := tx.Exec(stmt).Error; err != nil {
					return err
				}
			}
			return markMigrationAppliedWithDB(tx, version)
		})
		if err != nil {
			return fmt.Errorf("выполнение миграции %s: %w", version, err)
		}
	}

	return nil
}

func ensureSchemaMigrations(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(&schemaMigration{})
}

func isMigrationApplied(ctx context.Context, db *gorm.DB, version string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&schemaMigration{}).
		Where("version = ?", version).
		Count(&count).Error
	return count > 0, err
}

func markMigrationApplied(ctx context.Context, db *gorm.DB, version string) error {
	return markMigrationAppliedWithDB(db.WithContext(ctx), version)
}

func markMigrationAppliedWithDB(db *gorm.DB, version string) error {
	return db.Create(&schemaMigration{
		Version: version,
	}).Error
}
