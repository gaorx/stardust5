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

func (bp *Blueprint) CreateTableForMysql(tx *gorm.DB, opts *CreateTableOptions) error {
	opts1 := lo.FromPtr(opts)
	const sqlFn = "all.gen.sql"
	buffs, err := bp.Generate(MysqlDDL{
		TableIds:      opts1.TableIds,
		DisableFK:     opts1.DisableFK,
		WithDrop:      opts1.Drop,
		WithoutCreate: false,
		FileForCreate: sqlFn,
	})
	if err != nil {
		return sderr.WithStack(err)
	}
	sqlSrc := buffs.Get(sqlFn).String()
	if sqlSrc != "" {
		_, err := sdgorm.Exec(tx, sqlSrc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bp *Blueprint) FillDummyData(tx *gorm.DB, opts *FillDummyDataOptions) error {
	opts1 := lo.FromPtr(opts)
	tableIds := matchIds(bp.TableIds(), opts1.TableIds)
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
			dbr := tx.Table(t.NameForDB()).Create(map[string]any(record))
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
	return db.Transaction(func(tx *gorm.DB) error {
		err := bp.CreateTableForMysql(tx, &CreateTableOptions{Drop: true, DisableFK: false})
		if err != nil {
			return sderr.WithStack(err)
		}
		err = bp.FillDummyData(tx, &FillDummyDataOptions{})
		if err != nil {
			return sderr.WithStack(err)
		}
		return nil
	})
}
