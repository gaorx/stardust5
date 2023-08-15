package sdgorm

import (
	"github.com/gaorx/stardust5/sdfile"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"path/filepath"
	"testing"
)

type user struct {
	Id   int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
}

func TestTransaction(t *testing.T) {
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		db, err := Dial(Address{
			Driver: "sqlite",
			DSN:    filepath.Join(dirname, "test.db"),
		}, nil)
		assert.NoError(t, err)
		err = db.AutoMigrate(&user{})
		assert.NoError(t, err)
		ul, err := Transaction(db, func(tx *gorm.DB) ([]*user, error) {
			u0 := user{Name: "aaa", Age: 22}
			u1 := user{Name: "bbb", Age: 33}
			dbr := tx.Create(&u0)
			assert.NoError(t, dbr.Error)
			dbr = tx.Create(&u1)
			assert.NoError(t, dbr.Error)
			ul0, err0 := Find[*user](tx.Where("id == ? OR id == ?", u0.Id, u1.Id).Order("id ASC"))
			assert.NoError(t, err0)
			return ul0, nil
		})
		assert.NoError(t, err)
		assert.True(t, ul[0].Id == 1 && ul[0].Name == "aaa" && ul[0].Age == 22)
		assert.True(t, ul[1].Id == 2 && ul[1].Name == "bbb" && ul[1].Age == 33)
	})
}
