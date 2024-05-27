package flashkv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"unsafe"
)

const (
	batchSize = 1000
)

type DB struct {
	name    string
	changes bool

	sync.Mutex
	data map[string][]byte
}

func New() *DB {
	return &DB{
		data: make(map[string][]byte),
	}
}

func (d *DB) Load(name string) error {
	d.Lock()
	defer d.Unlock()

	d.name = name

	// --- blocks ---
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	fileSize, _ := file.Stat()

	buf := make([]byte, 1024*batchSize)
	bufPos := int64(0)
	i := 0
	for int64(bufPos) < fileSize.Size() {
		if i%batchSize == 0 {
			_, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			bufPos = 0
		}

		keySize := *(*int32)(unsafe.Pointer(&buf[bufPos]))
		key := string(buf[bufPos+4 : bufPos+4+int64(keySize)])

		var val []byte
		valSize := *(*int32)(unsafe.Pointer(&buf[bufPos+4+int64(keySize)]))
		if valSize > 0 {
			val = bytes.Clone(buf[bufPos+4+int64(keySize)+4 : bufPos+4+int64(keySize)+4+int64(valSize)])
		}
		d.data[key] = val
		bufPos += 4 + int64(keySize) + 4 + int64(valSize)
		i++
	}

	d.changes = true
	return err
}

func (d *DB) Save() error {
	if d.name == "" {
		return fmt.Errorf("no name to save database")
	}
	return d.SaveAs(d.name)
}

func (d *DB) SaveAs(name string) error {
	d.Lock()
	defer d.Unlock()

	//no changes
	if !d.changes {
		return nil
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf4 [4]byte
	buf := make([]byte, 0, 1024*batchSize)
	num := int64(0)
	for key, val := range d.data {
		keyBuf := []byte(key)
		keySize := (*int32)(unsafe.Pointer(&buf4[0]))
		*keySize = int32(len(keyBuf))
		buf = append(buf, buf4[:]...)
		buf = append(buf, keyBuf...)

		valBuf := []byte(val)
		valSize := (*int32)(unsafe.Pointer(&buf4[0]))
		*valSize = int32(len(valBuf))
		buf = append(buf, buf4[:]...)
		buf = append(buf, valBuf...)
		if (num+1)%batchSize == 0 {
			_, err := file.Write(buf)
			if err != nil {
				return err
			}
			buf = buf[:0]
		}
		num++
	}

	// Write remaining elements
	if len(buf) > 0 {
		_, err := file.Write(buf)
		if err != nil {
			return err
		}
	}

	d.changes = false
	return err
}

func (d *DB) Add(key string, val []byte) {
	d.Lock()
	defer d.Unlock()

	d.data[key] = val
	d.changes = true
}

// found stored value, if value found
func (d *DB) Get(key string) ([]byte, bool) {
	d.Lock()
	defer d.Unlock()

	val, ok := d.data[key]
	return val, ok
}

// dont modify database in iterate!
func (d *DB) Iterate(fn func(key string, val []byte) bool) {
	d.Lock()
	defer d.Unlock()

	for k, v := range d.data {
		if !fn(k, v) {
			return
		}
	}
}

func (d *DB) Exist(key string) bool {
	d.Lock()
	defer d.Unlock()

	_, ok := d.data[key]
	return ok
}

func (d *DB) Delete(key string) {
	d.Lock()
	defer d.Unlock()

	_, ok := d.data[key]
	if ok {
		delete(d.data, key)
		d.changes = true
	}
}

func (d *DB) Count() int {
	return len(d.data)
}
