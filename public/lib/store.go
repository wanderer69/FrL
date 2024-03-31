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

func (oc *Store) SaveFrameRecord(frameId *Value, slotName string, slotProperty string, slotValue *Value, debug int) error {
	if frameId.typev != VtString {
		return fmt.Errorf("frame id must be string")
	}
	frameIdS := frameId.String()
	//slotNameS := slotName
	//slotPropertyS := slotProperty
	slotValueS := ""
	if slotValue.typev == VtNil {

	} else {
		d, err := SaveValueStore(slotValue)
		if err != nil {
			return err
		}
		slotValueS = string(d)
	}

	h1 := sha1.New()
	h1.Write([]byte(frameIdS))
	frameIdSlotNameHash := hex.EncodeToString(h1.Sum([]byte(slotName)))

	h2 := sha1.New()
	h2.Write([]byte(slotName))
	slotNameSlotValueHash := hex.EncodeToString(h2.Sum([]byte(slotValueS)))

	_, _, err := oc.Sdb.StoreRecord(frameIdS, slotName, slotProperty, slotValueS, frameIdSlotNameHash, slotNameSlotValueHash)
	return err
}

func (oc *Store) LoadFrameRecord(debug int) (func() (frameId string, slotName string, slotProperty string, slotValue *Value, err error), error) {
	lazy, err := oc.Sdb.LoadLazyRecords(0)
	if err != nil {
		return nil, err
	}
	loader := func() (frameId string, slotName string, slotProperty string, slotValue *Value, err error) {
		rec, _, err := lazy()
		if err != nil {
			return "", "", "", nil, err
		}
		frameId = rec.FieldsValue[0]
		slotName = rec.FieldsValue[1]
		slotProperty = rec.FieldsValue[2]
		valueS := rec.FieldsValue[3]
		d, _, err := LoadValueStore([]byte(valueS))
		if err != nil {
			return "", "", "", nil, err
		}
		slotValue = d

		return
	}
	return loader, nil
}

func (oc *Store) Find(qri *QueryRelationItem) (func() (frameId *Value, name *Value, property *Value, value *Value, err error), error) {
	name := ""
	val := ""

	slotValue := ""
	if qri.Value.typev == VtNil {

	} else {
		d, err := SaveValueStore(qri.Value)
		if err != nil {
			return nil, err
		}
		slotValue = string(d)
	}

	switch qri.ObjectType {
	case "frame":
		name = "frame_id"
		/*
			h1 := sha1.New()
			h1.Write([]byte(slotValue))
			hex.EncodeToString(h1.Sum([]byte(slot_name_s)))
		*/
		val = slotValue
	case "relation":
		if len(slotValue) > 0 {
			name = "slot_name_slot_value"
			h2 := sha1.New()
			h2.Write([]byte(qri.Object))
			val = hex.EncodeToString(h2.Sum([]byte(slotValue)))
		} else {
			name = "slot_name"
			val = qri.Object
		}
	}
	ds, err, err1 := oc.Sdb.FindRecordIndexString([]string{name}, []string{val})
	if err1 != nil {
		fmt.Printf("Error %v %v\r\n", err, err1)
		return nil, fmt.Errorf("FindRecordIndexString error %v: %w", err, err1)
	}
	//Print("-> err %v\r\n", err)
	/*
		for k := range ds {
			fmt.Printf("ds[%v] %#v\r\n", k, ds[k])
		}
	*/
	pos := 0
	finder := func() (frame_id *Value, name *Value, property *Value, value *Value, err error) {
		//
		if len(ds) <= pos {
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
		//value = CreateValue(value_s)
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
