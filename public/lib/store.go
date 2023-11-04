package frl

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	smalldb "github.com/wanderer69/SmallDB/public/index"
	print "github.com/wanderer69/tools/parser/print"
)

type Store struct {
	Sdb   *smalldb.SmallDB
	Debug int
}

func NewStore(path string, o *print.Output) (*Store, error) {
	sdb := smalldb.InitSmallDB(path)
	//sdb := &psdb
	sdb.Debug = 0
	if !sdb.Inited {
		fl := []string{"frame_id", "slot_name", "slot_property", "slot_value", "frame_id_slot_name", "slot_name_slot_value"}
		err := sdb.CreateDB(fl, path)
		if err != nil {
			o.Print("Error creating DB %v\r\n", err)
			return nil, err
		}
		sdb.CreateIndex([]string{"frame_id"})
		sdb.CreateIndex([]string{"slot_name"})
		sdb.CreateIndex([]string{"frame_id_slot_name"})
		sdb.CreateIndex([]string{"slot_name_slot_value"})
		o.Print("CreateDB end\r\n")
		sdb = smalldb.InitSmallDB(path)
		//sdb = &psdb
	}
	sdb.OpenDB()
	oc := &Store{}
	oc.Sdb = sdb
	return oc, nil
}

func (oc *Store) SaveFrameRecord(frame_id *Value, slot_name string, slot_property string, slot_value *Value, debug int) error {
	if frame_id.typev != VtString {
		return fmt.Errorf("frame id must be string")
	}
	frame_id_s := frame_id.String()
	slot_name_s := slot_name
	slot_property_s := slot_property
	slot_value_s := ""
	if slot_value.typev == VtNil {

	} else {
		d, err := SaveValueStore(slot_value)
		if err != nil {
			return err
		}
		slot_value_s = string(d)
	}

	h1 := sha1.New()
	h1.Write([]byte(frame_id_s))
	frame_id_slot_name_hash := hex.EncodeToString(h1.Sum([]byte(slot_name_s)))

	h2 := sha1.New()
	h2.Write([]byte(slot_name_s))
	slot_name_slot_value_hash := hex.EncodeToString(h1.Sum([]byte(slot_value_s)))

	_, _, err := oc.Sdb.StoreRecord(frame_id_s, slot_name_s, slot_property_s, slot_value_s, frame_id_slot_name_hash, slot_name_slot_value_hash)
	return err
}

func (oc *Store) LoadFrameRecord(debug int) (func() (frame_id string, slot_name string, slot_property string, slot_value *Value, err error), error) {
	lazy, err := oc.Sdb.LoadLazyRecords(0)
	if err != nil {
		return nil, err
	}
	loader := func() (frame_id string, slot_name string, slot_property string, slot_value *Value, err error) {
		rec, _, err := lazy()
		if err != nil {
			return "", "", "", nil, err
		}
		frame_id = rec.FieldsValue[0]
		slot_name = rec.FieldsValue[1]
		slot_property = rec.FieldsValue[2]
		value_s := rec.FieldsValue[3]
		d, _, err := LoadValueStore([]byte(value_s))
		if err != nil {
			return "", "", "", nil, err
		}
		slot_value = d

		return
	}
	return loader, nil
}

func (oc *Store) Find(qri *QueryRelationItem, debug int, output *print.Output) (func() (frame_id *Value, name *Value, property *Value, value *Value, err error), error) {
	name := ""
	val := ""
	switch qri.ObjectType {
	case "frame":

	case "relation":

	}
	ds, err, err1 := oc.Sdb.FindRecordIndexString([]string{name}, []string{val})
	if err1 != nil {
		output.Print("Error %v %v\r\n", err, err1)
		return nil, err1
	}
	output.Print("-> err %v\r\n", err)
	for k := range ds {
		output.Print("ds[%v] %v\r\n", k, ds[k])
	}
	pos := 0
	finder := func() (frame_id *Value, name *Value, property *Value, value *Value, err error) {
		//
		if len(ds) > pos {
			return nil, nil, nil, nil, fmt.Errorf("empty")
		}

		rec := ds[pos]
		pos = pos + 1
		frame_id_s := rec.FieldsValue[0]
		name_s := rec.FieldsValue[1]
		property_s := rec.FieldsValue[2]
		value_s := rec.FieldsValue[3]
		frame_id = CreateValue(frame_id_s)
		name = CreateValue(name_s)
		property = CreateValue(property_s)
		value = CreateValue(value_s)
		d, _, err := LoadValueStore([]byte(value_s))
		if err != nil {
			return nil, nil, nil, nil, err
		}
		value = d

		return
	}
	return finder, nil
}

func (oc *Store) CloseSmallDB() {
	oc.Sdb.CloseData()
}
