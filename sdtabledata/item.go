package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdregexp"
	"io/fs"
	"path/filepath"
	"regexp"
)

type itemKind string

const (
	fileUnidentified itemKind = "unidentified" // 未识别的文件，文件名匹配不上任何规则
	fileMeta         itemKind = "meta"         // meta file
	fileRow          itemKind = "row"          // 行数据文件
	fileColumn       itemKind = "column"       // 列数据文件
	fileColumnSub    itemKind = "column_sub"   // 列中的字段数据文件
)

type item struct {
	kind     itemKind
	filename string // 相对于root的文件名

	// 解析filename的结果
	row    string
	column string
	sub    string
	ext    string
}

var fileKindPatterns = []struct {
	p    *regexp.Regexp
	kind itemKind
}{
	// meta
	{regexp.MustCompile(`^__meta__.json$`), fileMeta},

	// row
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)\.json$`), fileRow},
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)/index\.json$`), fileRow},

	// column
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)\.(?P<column>\w+)\.(?P<ext>\w+)$`), fileColumn},
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)/(?P<column>\w+)\.(?P<ext>\w+)$`), fileColumn},

	// sub
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)\.(?P<column>\w+)\.(?P<sub>\w+)\.(?P<ext>\w+)$`), fileColumnSub},
	{regexp.MustCompile(`^(?P<row>[^ \f\n\r\t\v./]+)/(?P<column>\w+)\.(?P<sub>\w+)\.(?P<ext>\w+)$`), fileColumnSub},
}

func loadItems(src Source) ([]*item, error) {
	fsys, dir := src.Root, src.Dir
	var filenames []string
	err := fs.WalkDir(fsys, dir, func(fn string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		fn1, err := filepath.Rel(dir, fn)
		if err != nil {
			return sderr.WithStack(err)
		}
		filenames = append(filenames, fn1)
		return nil
	})
	if err != nil {
		return nil, sderr.Wrap(err, "scan file items error")
	}

	var items []*item
	for _, filename := range filenames {
		ok := false
		for _, x := range fileKindPatterns {
			if x.p.MatchString(filename) {
				g := sdregexp.FindStringSubmatchGroup(x.p, filename)
				items = append(items, &item{
					kind:     x.kind,
					filename: filename,
					row:      g["row"],
					column:   g["column"],
					sub:      g["sub"],
					ext:      g["ext"],
				})
				ok = true
				break
			}
		}
		if !ok {
			items = append(items, &item{
				kind:     fileUnidentified,
				filename: filename,
			})
		}
	}
	// sdprint.PrettyL(items)
	return items, nil
}
