package frl

import (
	"fmt"

	print "github.com/wanderer69/tools/parser/print"
)

type DataBaseType string

const (
	DataBaseTypeSimple DataBaseType = "simple"
)

type DataBaseState string

const (
	DataBaseStateNotConnected DataBaseState = "not_connected"
	DataBaseStateConnected    DataBaseState = "connected"
)

type DataBase struct {
	dbType DataBaseType
	state  DataBaseState
	oc     *Store
	//debug  int
}

func NewDataBase() *DataBase {
	return &DataBase{
		state: DataBaseStateNotConnected,
	}
}

func (db *DataBase) Connect(typeDB DataBaseType, pathToDB string, output *print.Output) error {
	switch typeDB {
	case DataBaseTypeSimple:
		ns, err := NewStore(pathToDB, output)
		if err != nil {
			return err
		}
		db.oc = ns
	}
	db.dbType = typeDB
	db.state = DataBaseStateConnected
	return nil
}

func (db *DataBase) Close() error {
	db.oc.CloseSmallDB()
	return nil
}

func (db *DataBase) StoreFrame(frame *Frame) error {
	ff := frame.Iterate()
	frameIDs, err := frame.GetValue("ID")
	if err != nil {
		return err
	}

	frameID := frameIDs[0]
	for {
		s, ok, err := ff()
		if err != nil {
			break
		}
		ssl := s.GetSlotValue()
		slotName := s.GetSlotName()
		slotProperty := s.GetSlotProperty()
		if slotName != "ID" {
			for j := range ssl {
				err := db.oc.SaveFrameRecord(frameID, slotName, slotProperty, ssl[j], 0)
				if err != nil {
					return err
				}
			}
		}
		if ok {
			break
		}
	}
	return nil
}

func (db *DataBase) LoadFrameByID(frameID *Value) (*Frame, error) {
	fn, err := db.oc.FindShort(&QueryRelationItem{ObjectType: "frame_id", Value: frameID})
	if err != nil {
		return nil, err
	}
	var f *Frame
	for {
		frameId, slotName, slotProperty, slotValue, err := fn()
		if err != nil {
			break
		}
		if f == nil {
			f = NewFrame()
			// добавляем поле уникального идентификатора
			id := "ID"
			err = f.AddSlot(id)
			if err != nil {
				return nil, err
			}
			_, err := f.SetValue(id, frameId)
			if err != nil {
				return nil, err
			}
			/*
				fe.AddRelations(f, AddRelationItem{ObjectType: "frame", Object: "", Value: v})
				fe.Frames = append(fe.Frames, f)
			*/
			continue
		}
		err = f.AddSlot(slotName)
		if err != nil {
			return nil, err
		}
		err = f.SetSlotProperty(slotName, slotProperty)
		if err != nil {
			return nil, err
		}
		_, err = f.SetValue(slotName, slotValue)
		if err != nil {
			return nil, err
		}
		//fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: slot_name, Value: slot_value})
	}
	return f, nil
}

func (db *DataBase) FindFrames(frameTemplate *Value) ([]*Frame, error) {
	if frameTemplate.GetType() != VtFrame {
		return nil, fmt.Errorf("must be frame, has %v", frameTemplate.GetType())
	}
	frame := frameTemplate.Frame()
	fn, err := db.oc.FindByTemplate(frame)
	if err != nil {
		return nil, err
	}
	fs := []*Frame{}
	var f *Frame
	currentFrameID := ""
	for {
		frameId, slotName, slotProperty, slotValue, err := fn()
		if err != nil {
			break
		}
		if currentFrameID != frameId.String() && f != nil {
			fs = append(fs, f)
			f = nil
		}
		if f == nil {
			f = NewFrame()
			// добавляем поле уникального идентификатора
			id := "ID"
			err = f.AddSlot(id)
			if err != nil {
				return nil, err
			}
			_, err := f.SetValue(id, frameId)
			if err != nil {
				return nil, err
			}
			currentFrameID = frameId.String()
			//continue
		}
		err = f.AddSlot(slotName)
		if err != nil {
			return nil, err
		}
		err = f.SetSlotProperty(slotName, slotProperty)
		if err != nil {
			return nil, err
		}
		_, err = f.SetValue(slotName, slotValue)
		if err != nil {
			return nil, err
		}
		//fe.AddRelations(f, frl.AddRelationItem{ObjectType: "relation", Object: slot_name, Value: slot_value})
	}
	return fs, nil
}
