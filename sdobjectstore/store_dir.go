package sdobjectstore

import (
	"github.com/gaorx/stardust5/sderr"
	"io"
	"os"
	"path/filepath"
)

type Dir struct {
	Root string
}

func (d Dir) Store(src Source, objectName string) (*Target, error) {
	if src == nil {
		return nil, sderr.New("nil source")
	}
	// 展开文件名称
	expandedObjectName, err := expandObjectName(src, objectName)
	if err != nil {
		return nil, err
	}

	absRoot, err := filepath.Abs(d.Root)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	absFn := filepath.Join(absRoot, expandedObjectName)
	err = os.MkdirAll(filepath.Dir(absFn), 0755)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	in, err := src.Open()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() {
		if c, ok := in.(io.Closer); ok {
			_ = c.Close()
		}
	}()
	out, err := os.OpenFile(absFn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() { _ = out.Close() }()
	_, err = io.Copy(out, in)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &Target{
		Typ:            FileTarget,
		Prefix:         absRoot,
		InternalPrefix: absRoot,
		Path:           expandedObjectName,
	}, nil
}
