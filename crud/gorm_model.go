package crud

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Model interface {
	Identity() (fieldName string, value any)
}

type GormModel struct {
	ID        uint           `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type ModelBase GormModel

func (m ModelBase) Identity() (fieldName string, value any) {
	fmt.Printf("Init Identity %+v", m)
	return "ID", m.ID
}
func (m ModelBase) UpdateForID(db, id any) {
	// GLOB DB???
}

func LoadGormDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, _ := gorm.Open(sqlite.Open("file:autodb.db?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	return db
}
