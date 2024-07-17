package config

//客户端注册消息
type RegisterMessage struct {
	Token string
	Event string
}

//服务端返回消息
type ResMessage struct {
	Error int         `json:"error"`
	Msg   string      `json:"msg"`
	Data  MessageData `json:"data"`
	Event string      `json:"event"`
}
//更新已读的参数
type UpdateMessage struct {
	ID int64
	MsgType    int    //消息类型 1发送在线用户即时消息 2登录后必达消息
	BusMsgType int
	IsRead int
	UserIds    string //用户id以,号分隔
	MsgKey    string //MsgKey 业务消息key
	SenderId   int64  //发送者id
}
//推送数据结构
type PushMessage struct {
	Token      string //api token
	MsgType    int    //消息类型 1发送在线用户即时消息 2登录后必达消息
	BusMsgType int
	MsgKey string
	UserIds    string //用户id以,号分隔
	SenderId   int64  //发送者id
	SenderName string //发送者姓名
	Title      string //消息标题
	Content    string //消息内容
	Options    string //消息弹窗配置
	Timestamp  string //时间戳
}

//发送给客端data消息数据结构
type MessageData struct {
	SenderId   int64  `json:"senderId"`
	SenderName string `json:"senderName"`
	MsgTime    string `json:"msgTime"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Options    string `json:"options"`
	MsgId      int64  `json:"msgId,omitempty"`    //推送消息数据库记录id
	MsgType    int    `json:"msgType,omitempty"`  //消息类型
	BusMsgType int    `json:"busMsgType,omitempty"`  //业务消息类型
	NoReadTotal int  `json:"noReadTotal,omitempty"`  //未读消息总数量
	MsgKey    string `json:"msgKey"`
	MsgLogId   int64  `json:"msgLogId,omitempty"` //用户消息数据库记录id
	MessList   []Message `json:"messList,omitempty"` //未读消息列表
}
func (Message) TableName() string {
	return "xhx_push_message"
}
type Message struct{
	ID      int64  `json:"id,omitempty"`    //推送消息数据库记录id
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime    string `json:"createTime"`
}