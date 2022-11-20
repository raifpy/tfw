package src

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

type DatabaseOptions struct {
	Path string
}

type Database struct {
	DB *gorm.DB
}

func (d *Database) New(s GeneralDatabaseSchema) error {
	return d.DB.Create(&s).Error
}

func (d *Database) CountAll() (i int64) {
	var o GeneralDatabaseSchema
	d.DB.Find(&o).Count(&i)
	return
}

func (d *Database) CountCustom(l []string) (i int64, err error) {
	var o GeneralDatabaseSchema
	db := d.DB.Find(&o)
	for i, v := range l {
		if i%2 == 1 {
			continue
		}
		db.Where(v+" = ?", l[i+1])
	}

	err = db.Count(&i).Error
	return
}

func (d *Database) FindFirst(key string, value any) (g GeneralDatabaseSchema, err error) {

	err = d.DB.First(&g, key+" = ?", value).Error
	g.d = d
	return
}

func (d *Database) FindFirstMulti(field ...string) (g GeneralDatabaseSchema, err error) {
	var gdb *gorm.DB
	for i, v := range field {
		if i%2 == 1 {
			continue
		}

		if gdb == nil {
			gdb = d.DB.Where(v+" = ?", field[i+1])
		} else {
			gdb.Where(v+" = ?", field[i+1])
		}

	}

	if err = gdb.First(&g).Error; err == nil {
		g.d = d
	}

	return
}

func (d *Database) FindAllMulti(field ...string) (g []GeneralDatabaseSchema, err error) {
	var gdb *gorm.DB
	for i, v := range field {
		if i%2 == 1 {
			continue
		}

		if gdb == nil {
			gdb = d.DB.Where(v+" = ?", field[i+1])
		} else {
			gdb.Where(v+" = ?", field[i+1])
		}

	}

	if err = gdb.Find(&g).Error; err == nil {
		for i, _ := range g {
			g[i].d = d
		}
	}

	return
}

func (d *Database) FindAll(key string, value any) (g []GeneralDatabaseSchema, err error) {

	err = d.DB.Find(&g, key+" = ?", value).Error
	for i, _ := range g {
		g[i].d = d
	}
	return
}

func (d *Database) GetAll() (g []GeneralDatabaseSchema, err error) {

	err = d.DB.Find(&g).Error
	for i, _ := range g {
		g[i].d = d
	}
	return
}

func NewDatabase(o DatabaseOptions) (db *Database, err error) {
	db = &Database{}

	if db.DB, err = gorm.Open(sqlite.Open(o.Path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}); err != nil {
		return
	}
	err = db.init()
	return
}
func (d *Database) init() error {
	return d.DB.AutoMigrate(&GeneralDatabaseSchema{})
}
