package sdblueprint

import (
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdfile"
	"github.com/gaorx/stardust5/sdgorm"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

func (bp *Blueprint) RunCli() {
	app := &cli.App{
		Name:  "Blueprint tool",
		Usage: "go run blueprint.go gen|mock-db",
		Commands: []*cli.Command{
			{
				Name:  "gen",
				Usage: "generate source from blueprint",
				Action: func(cc *cli.Context) error {
					return bp.cliGenerate(cc)
				},
			},
			{
				Name:  "mock-db",
				Usage: "create tables / fill dummy data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "groups",
						Aliases: []string{"g"},
						Usage:   "for some group",
					},
				},
				Action: func(cc *cli.Context) error {
					return bp.cliMockDB(cc)
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func (bp *Blueprint) cliGenerate(cc *cli.Context) error {
	moduleIds, err := bp.cliModuleIds(cc, "gen")
	if err != nil {
		return sderr.WithStack(err)
	}
	if len(moduleIds) <= 0 {
		return nil
	}
	buffs, err := bp.Generate(ForModule{Ids: moduleIds})
	if err != nil {
		return sderr.WithStack(err)
	}
	root, ok, err := getProjectRoot()
	if err != nil {
		return sderr.WithStack(err)
	}
	if !ok {
		return sderr.New("not in golang project")
	}
	err = buffs.Save(root, sdcodegen.SimplePrint)
	if err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func (bp *Blueprint) cliMockDB(cc *cli.Context) error {
	getEnv := func(key string, def string) string {
		v := os.Getenv(key)
		if v == "" {
			return def
		}
		return v
	}
	dbDriver := getEnv("SD_DB_DRIVER", "")
	dbDSN := getEnv("SD_DB_DSN", "")
	if dbDriver == "" {
		return sderr.New("no env SD_DB_DRIVER")
	}
	if dbDSN == "" {
		return sderr.New("no env SD_DB_DSN")
	}
	groups := sdstrings.SplitNonempty(cc.String("groups"), ",", true)
	bp1 := getSub(bp, groups)
	if err := bp1.MockDB(sdgorm.Address{
		Driver: dbDriver,
		DSN:    dbDSN,
	}); err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func (bp *Blueprint) cliModuleIds(cc *cli.Context, sub string) ([]string, error) {
	moduleIds := lo.FilterMap(cc.Args().Slice(), func(id string, _ int) (string, bool) {
		id = strings.TrimSpace(id)
		return id, id != ""
	})
	if len(moduleIds) <= 0 {
		return nil, sderr.NewWith("./sdcli %s <module>", sub)
	}
	if lo.Contains(moduleIds, "all") {
		moduleIds = bp.ModuleIds()
	}
	return moduleIds, nil
}

func getProjectRoot() (string, bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false, sderr.Wrap(err, "get current work directory error(1)")
	}
	d, err := filepath.Abs(wd)
	if err != nil {
		return "", false, sderr.Wrap(err, "get current work directory error(2)")
	}
	for {
		goMod := filepath.Join(d, "go.mod")
		if sdfile.Exists(goMod) {
			return d, true, nil
		} else {
			d = filepath.Dir(d)
			if filepath.ToSlash(d) == "/" {
				break
			}
		}
	}
	return "", false, nil
}
