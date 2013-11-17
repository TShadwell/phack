//Package persistent provides persistent storage.
package persistent

import (
	"bytes"
	lvl"github.com/TShadwell/level"
	levigo "github.com/TShadwell/level/levigo"
	"github.com/golang/glog"
	"encoding/gob"
	"errors"
)

var Inexistant = errors.New("persistent: Record does not exist");

var (
	level = levigo.Level
	database = &lvl.Database {
		Cache: level.NewCache(500 * lvl.Megabyte),
		Options: level.NewOptions().SetCreateIfMissing(
			true,
		),
	}
)

func Store(i interface{}, k []byte) (err error){
	bt, err := gb(i)
	if err != nil {
		return
	}
	return database.Put(k, bt)
}

func Retreive(i interface{}, k []byte) (err error) {
	bt, err := database.Get(k)
	if err != nil {
		return
	}
	return ungob(i, bt)
}

type AtomItem struct {
	I interface{}
	K []byte
}

type Atom []AtomItem

func (an Atom) Store() (err error) {
	a := new(lvl.Atom)
	for _, it := range an {
		var bt []byte
		bt, err = gb(it.I)
		if err != nil {
			return
		}
		a.Put(bt, it.K)

	}
	return database.Commit(a)
}

func gb(i interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(i)
	return b.Bytes(), err
}

func ungob(i interface{}, b []byte) error {
	if b == nil {
		return Inexistant
	}
	return gob.NewDecoder(bytes.NewReader(b)).Decode(i)
}

func init() {
	if err := level.OpenDatabase(database, "LevelDatabase"); err != nil{
		glog.Fatalf("Could not load database due to error %+q", err)
	}
}
