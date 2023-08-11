package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdgorm"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateTableOptions struct {
	TableIds  []string
	Drop      bool
	DisableFK bool
}

type FillDummyDataOptions struct {
	TableIds []string
}

func (bp *Blueprint) CreateTableForMysql(db *gorm.DB, opts CreateTableOptions) error {
	const sqlFilename = "all.gen.sql"
	buffs, err := bp.Generate(MysqlDDL{
		TableIds:      opts.TableIds,
		DisableFK:     opts.DisableFK,
		WithDrop:      opts.Drop,
		WithoutCreate: false,
		FileForCreate: sqlFilename,
	})
	if err != nil {
		return sderr.WithStack(err)
	}
	code := buffs.Get(sqlFilename).String()
	if code != "" {
		err = db.Transaction(func(tx *gorm.DB) error {
			var r any
			dbr := tx.Raw(code).Scan(&r)
			if dbr.Error != nil {
				return dbr.Error
			}
			return nil
		})
		if err != nil {
			return sderr.Wrap(err, "execute create table sql error")
		}
	}
	return nil
}

func (bp *Blueprint) FillDummyData(db *gorm.DB, opts FillDummyDataOptions) error {
	tableIds := matchIds(bp.TableIds(), opts.TableIds)
	for _, tableId := range tableIds {
		t := bp.Table(tableId)
		if t == nil {
			return sderr.NewWith("not found table", t.Id())
		}
		if t == nil {
			return sderr.NewWith("not found table for fill dummy data", tableId)
		}
		records := t.DummyData()
		if len(records) <= 0 {
			continue
		}
		records = lo.Map(records, func(record DummyRecord, _ int) DummyRecord { return record.ToDB(t) })
		for _, record := range records {
			dbr := db.Table(t.NameForDB()).Create(map[string]any(record))
			if dbr.Error != nil {
				return sderr.WrapWith(dbr.Error, "fill dummy data error", t.Id())
			}
		}
	}
	return nil
}

func (bp *Blueprint) MockDB(addr sdgorm.Address) error {
	db, err := sdgorm.Dial(addr, nil)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = bp.CreateTableForMysql(db, CreateTableOptions{Drop: true, DisableFK: false})
	if err != nil {
		return sderr.WithStack(err)
	}
	err = bp.FillDummyData(db, FillDummyDataOptions{})
	if err != nil {
		return sderr.WithStack(err)
	}
	return nil
}
