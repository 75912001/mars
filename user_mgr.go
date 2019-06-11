package main

import (
	"github.com/75912001/goz/ztcp"
)

//UserIDMap 用户map
type UserIDMap map[UserID]*User

type userMgr struct {
	ztcp.PeerConnMgr

	userIDMap UserIDMap
}

//GuserMgr 用户管理器
var GuserMgr userMgr

func init() {
	ztcp.PeerConnMgr.Init()
	GuserMgr.Init()
}

//Init 初始化
func (p *userMgr) Init() {
	p.userIDMap = make(UserIDMap)
}

//AddUserID 加
func (p *userMgr) Add(uid UserID, user *User) {
	p.userIDMap[uid] = user
}

//DelUserID 删
func (p *userMgr) Del(uid UserID) {
	delete(p.userIDMap, uid)
}

//FindID 查
func (p *userMgr) Find(uid UserID) (user *User) {
	user, _ = p.userIDMap[uid]
	return user
}
