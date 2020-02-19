package models
import (
	"time"

	"github.com/shiyongabc/go-websocket-broadcast/src/config"
)

type PushMessageModel struct {
	ID         int64
	Title      string
	Content    string
	Options    string
	MsgType    int
	BusMsgType int
	IsRead int
	UserIds    string
	SenderId   int64
	SenderName string
	CreateTime string
}

func (PushMessageModel) TableName() string {
	return "xhx_push_message"
}

func (PushMessageModel) Create(m PushMessageModel) int64 {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return 0
	}

	m.CreateTime = time.Now().Format(config.TIMESTAMP_FORMAT)

	db.Create(&m)
	defer db.Close()
	return m.ID
}
func (PushMessageModel) Update(m PushMessageModel) int64 {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return 0
	}

	m.CreateTime = time.Now().Format(config.TIMESTAMP_FORMAT)
    if m.ID!=0{
		//db=db.Update(&m).Where("id=?",m.ID)
		db=db.Model(PushMessageModel{}).Where("id=?",m.ID).Update("is_read",1)
	}else{
		//db=db.Update(&m).Where("bus_msg_type=? and user_ids=?",m.BusMsgType,m.UserIds)
		db=db.Model(PushMessageModel{}).Where("bus_msg_type=? and user_ids=?",m.BusMsgType,m.UserIds).Update("is_read",1)
	}

	defer db.Close()
	return db.RowsAffected
}