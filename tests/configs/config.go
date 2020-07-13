package configs

import (
	"os"
	"path"
	"path/filepath"

	"github.com/xiaolingzi/lingorm"
)

func GetDB() lingorm.IQuery {
	dir, _ := os.Getwd()
	dir, _ = filepath.Abs(filepath.Dir(dir))
	os.Setenv("LINGORM_CONFIG", path.Join(dir, "configs/database.json"))
	db := lingorm.DB("test")
	return db
}
