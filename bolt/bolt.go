package bolt

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/utilgo/stringx"
	"go.etcd.io/bbolt"
	"math"
	"os"
	"path/filepath"
	"regexp"
)

// Select 数据检索，支持时间范围和前缀扫描
func Select(db *bbolt.DB, min, max int64, prefix string, limit int) ([]Element, error) {
	if limit <= 0 {
		limit = math.MaxInt64
	}
	minS := fmt.Sprintf("%d", min)
	maxS := fmt.Sprintf("%d", max)
	seek := stringx.StringToBytes(&minS)
	compare := stringx.StringToBytes(&maxS)
	data := make([]Element, 0, 100)
	view := func(tx *bbolt.Tx) (_ error) {
		cur := tx.Bucket(bucket).Cursor()
		reg := regexp.MustCompile(prefix)
		for k, v := cur.Seek(seek); k != nil && bytes.Compare(k, compare) <= 0; k, v = cur.Next() {
			if prefix != "" {
				if !reg.Match(k) {
					continue
				}
			}
			if len(data) >= limit {
				break
			}
			data = append(data, Element{stringx.BytesToString(k), v})
		}
		return
	}
	if err := db.View(view); err != nil {
		return nil, errors.New(err.Error())
	}
	return data, nil
}

// SetAll 批量写入数据，value != nil ? insert : delete
func SetAll(db *bbolt.DB, elements []Element) (_ error) {
	update := func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket(bucket)
		for i := range elements {
			if elements[i].Value != nil {
				err = b.Put(stringx.StringToBytes(&elements[i].Index), elements[i].Value)
			} else {
				err = b.Delete(stringx.StringToBytes(&elements[i].Index))
			}
			if err != nil {
				return
			}
		}
		return
	}
	if err := db.Update(update); err != nil {
		return errors.New(err.Error())
	}
	return
}

// New 打开数据库
func New(dbname string) *bbolt.DB {
	_, readErr := os.Stat(dbname)
	_ = os.MkdirAll(filepath.Dir(dbname), 0777)
	db, err := bbolt.Open(dbname, 0666, nil)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	if readErr != nil {
		update := func(tx *bbolt.Tx) error {
			b, createErr := tx.CreateBucketIfNotExists(bucket)
			if createErr != nil {
				return createErr
			}
			return b.Put([]byte("test"), []byte("test"))
		}
		if err = db.Update(update); err != nil {
			panic(errors.New(err.Error()))
		}
		_ = SetAll(db, ValuesByMap(map[string][]byte{"test": nil}))
	}
	return db
}

type Element struct {
	Index string
	Value []byte
}

var bucket = []byte("bucket")

func ValuesByMap(data map[string][]byte) []Element {
	elements := make([]Element, 0, len(data))
	for index, value := range data {
		elements = append(elements, Element{index, value})
	}
	return elements
}

func Get(db *bbolt.DB, index string) (value []byte, _ error) {
	view := func(tx *bbolt.Tx) error {
		value = tx.Bucket(bucket).Get(stringx.StringToBytes(&index))
		return nil
	}
	if err := db.View(view); err != nil {
		return nil, errors.New(err.Error())
	}
	return
}

func GetKeyAll(db *bbolt.DB) ([]string, error) {
	keys := make([]string, 0, 100)
	view := func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, stringx.BytesToString(k))
			return nil
		})
	}
	if err := db.View(view); err != nil {
		return nil, errors.New(err.Error())
	}
	return keys, nil
}

func DeleteAll(db *bbolt.DB, keys []string) error {
	update := func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		for _, v := range keys {
			err := b.Delete(stringx.StringToBytes(&v))
			if err != nil {
				return err
			}
		}
		return nil
	}
	if err := db.Update(update); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
