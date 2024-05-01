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
		err = sdb.CreateIndex([]string{"frame_id"})
		if err != nil {
			o.Print("Error creating index frame_id %v\r\n", err)
			return nil, err
		}
		err = sdb.CreateIndex([]string{"slot_name"})
		if err != nil {
			o.Print("Error creating index slot_name %v\r\n", err)
			return nil, err
		}
		err = sdb.CreateIndex([]string{"frame_id_slot_name"})
		if err != nil {
			o.Print("Error creating index frame_id_slot_name %v\r\n", err)
			return nil, err
		}
		err = sdb.CreateIndex([]string{"slot_name_slot_value"})
		if err != nil {
			o.Print("Error creating index slot_name_slot_value %v\r\n", err)
			return nil, err
		}
		o.Print("CreateDB end\r\n")
		sdb = smalldb.InitSmallDB(path)
	}
	err := sdb.OpenDB()
	if err != nil {
		o.Print("Error open DB %v\r\n", err)
		return nil, err
	}
	oc := &Store{
		Sdb: sdb,
	}
	return oc, nil
}

func (oc *Store) SaveFrameRecord(frameId *Value, slotName string, slotProperty string, slotValue *Value, debug int) error {
	if frameId.typev != VtString {
		return fmt.Errorf("frame id must be string")
	}
	frameIdS := frameId.String()
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
	pos := 0
	finder := func() (frame_id *Value, name *Value, property *Value, value *Value, err error) {
		if len(ds) <= pos {
			return nil, nil, nil, nil, fmt.Errorf("empty")
		}

		rec := ds[pos]
		pos = pos + 1
		frameIDs := rec.FieldsValue[0]
		names := rec.FieldsValue[1]
		properties := rec.FieldsValue[2]
		values := rec.FieldsValue[3]
		frame_id = CreateValue(frameIDs)
		name = CreateValue(names)
		property = CreateValue(properties)
		d, _, err := LoadValueStore([]byte(values))
		if err != nil {
			return nil, nil, nil, nil, err
		}
		value = d

		return
	}
	return finder, nil
}

func (oc *Store) FindShort(qri *QueryRelationItem) (func() (frameId *Value, name string, property string, value *Value, err error), error) {
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
		val = qri.Value.String() //slotValue
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
	pos := 0
	finder := func() (frameID *Value, name string, property string, value *Value, err error) {
		if len(ds) <= pos {
			return nil, "", "", nil, fmt.Errorf("empty")
		}

		rec := ds[pos]
		pos = pos + 1
		frameIDs := rec.FieldsValue[0]
		name = rec.FieldsValue[1]
		property = rec.FieldsValue[2]
		values := rec.FieldsValue[3]
		frameID = CreateValue(frameIDs)
		//name = CreateValue(names)
		//property = CreateValue(properties)
		d, _, err := LoadValueStore([]byte(values))
		if err != nil {
			return nil, "", "", nil, err
		}
		value = d

		return
	}
	return finder, nil
}

func (oc *Store) FindByTemplate(frame *Frame) (func() (frameId *Value, name string, property string, value *Value, err error), error) {
	names := []string{}
	values := []string{}

	ff := frame.Iterate()
	for {
		s, ok, err := ff()
		if err != nil {
			break
		}
		ssl := s.GetSlotValue()
		slotName := s.GetSlotName()
		//slotProperty := s.GetSlotProperty()
		switch slotName {
		case "ID":
			name := "frame_id"
			names = append(names, name)
			/*
				d, err := SaveValueStore(ssl[0])
				if err != nil {
					return nil, err
				}
			*/
			values = append(values, ssl[0].String())
		default:
			if len(ssl) > 0 {
				for j := range ssl {
					name := "slot_name_slot_value"
					h2 := sha1.New()
					h2.Write([]byte(slotName))
					d, err := SaveValueStore(ssl[j])
					if err != nil {
						return nil, err
					}
					val := hex.EncodeToString(h2.Sum([]byte(d)))

					names = append(names, name)

					values = append(values, string(val))
				}
			} else {
				name := "slot_name"
				names = append(names, name)
				values = append(values, string(slotName))
			}

		}
		if ok {
			break
		}
	}
	ds, err, err1 := oc.Sdb.FindRecordIndexString(names, values)
	if err1 != nil {
		fmt.Printf("Error %v %v\r\n", err, err1)
		return nil, fmt.Errorf("FindRecordIndexString error %v: %w", err, err1)
	}
	pos := 0
	finder := func() (frameID *Value, name string, property string, value *Value, err error) {
		if len(ds) <= pos {
			return nil, "", "", nil, fmt.Errorf("empty")
		}

		rec := ds[pos]
		pos = pos + 1
		frameIDs := rec.FieldsValue[0]
		name = rec.FieldsValue[1]
		property = rec.FieldsValue[2]
		values := rec.FieldsValue[3]
		frameID = CreateValue(frameIDs)
		//name = CreateValue(names)
		//property = CreateValue(properties)
		d, _, err := LoadValueStore([]byte(values))
		if err != nil {
			return nil, "", "", nil, err
		}
		value = d

		return
	}
	return finder, nil
}

func (oc *Store) CloseSmallDB() {
	oc.Sdb.CloseData()
}
