package server
import (
	"encoding/json"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
	"go-websocket-broadcast/src/config"
	"go-websocket-broadcast/src/models"
	"go-websocket-broadcast/src/utils"
	log "github.com/sirupsen/logrus"
)
type Client struct {
	ID           string
	Socket       *websocket.Conn
	Send         chan []byte
	UserId       string
	Token        string
	RegisterTime int64
}

//客户端连接后激活这里读取并注册client
func (c *Client) Read() {
	defer c.Close()

	for {
		_, message, err := c.Socket.ReadMessage()

		//如果读取不到token数据关闭连接
		if err != nil {
			c.Close()
			log.Debug("读取socket数据失败:%s",err.Error())
			break
		}

		var rm config.RegisterMessage
		err1 := json.Unmarshal(message, &rm)
		if err1 != nil {
			c.CloseAndRes(101, "解析数据失败，请检查数据格式", "register")
			break
		}

		if rm.Event != "register" {
			continue
		}

		if rm.Token == "" {
			c.CloseAndRes(102, "token 必传", "register")
			break
		}


		if !c.CheckToken(rm.Token) {
			c.CloseAndRes(100, "token过期", "register")
			break
		}

		c.Token = rm.Token
		userId:=utils.ObtainUserByToken(c.Token,"userId")
		//log.Printf("userId=",userId)
		c.UserId=userId
		jsonMessage, _ := json.Marshal(&config.ResMessage{Error: 0, Msg: "ok", Event: "register"})
		c.Send <- jsonMessage

		//发送必读消息
		go c.SendMustReadMsg()
	}
}

//客户端消息写入程序
func (c *Client) Write() {
	defer c.Close()

	for {
		select {
		case message, ok := <-c.Send: //发送数据
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var res config.ResMessage
			json.Unmarshal(message, &res)

			if !c.CheckToken("") {
				c.CloseAndRes(200, "token过期", "messaage")
				go c.SaveUserMsgLog(res.Data, 2)
				break
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				go c.SaveUserMsgLog(res.Data, 2)
				c.Close()
				break
			}

			go c.SaveUserMsgLog(res.Data, 1)
		}
	}
}

//保存用户消息日志
func (c *Client) SaveUserMsgLog(data config.MessageData, status int) {
	var pushMsgLogModel models.PushMessageLogModel
	if data.MsgLogId == 0 {
		pushMsgLogModel.Create(data.MsgId, data.MsgType, c.UserId, c.ID, status)
	} else {
		pushMsgLogModel.Save(data.MsgLogId, c.ID, status)
	}
}

//发送给用户最近未读的必达消息
func (c *Client) SendMustReadMsg() {
	var pushMsgLogModel models.PushMessageLogModel
	log.Printf("c.userId=",c.UserId)
	msg_limit:=time.Now().Unix()-config.LAST_MSG_TIME_LIMIT
	log.Printf("msg_limit=",msg_limit)
	msgData := pushMsgLogModel.GetMustReadMsgByUserId(c.UserId, msg_limit)
	log.Printf("msgData=",len(msgData))
	for _, row := range msgData {
		message, _ := json.Marshal(&config.ResMessage{Error: 0, Msg: "ok", Event: "message", Data: row})
		c.Send <- message
	}
}

func (c *Client) SendRes(res config.ResMessage) {
	resJson, _ := json.Marshal(res)
	err := c.Socket.WriteMessage(websocket.TextMessage, resJson)
	if err != nil {
		log.Error("Send Res Error:" + c.ID)
		c.Close()
	}
}

func (c *Client) Close() {
	Manager.Unregister <- c
	c.Socket.Close()
}

func (c *Client) CloseAndRes(errCode int, msg string, event string) {
	c.SendRes(config.ResMessage{Error: errCode, Msg: msg, Event: event})
	c.Close()
}

//check user login token
func (c *Client) CheckToken(reqToken string) bool {
	if reqToken == "" {
		reqToken = c.Token
	}

	if reqToken == "" {
		return false
	}

	var expSecond string
	expSecond=utils.ObtainUserByToken(reqToken,"exp")
	 exp, err:= strconv.ParseInt(expSecond, 10, 64)
     //log.Printf("exp=",exp)
     if err!=nil{
     	log.Printf("err=",err.Error())
	 }
	if (time.Now().After(time.Unix(exp,0))){
		return false
	}
	return true
}
