// icartravel project main.go
package main

import (
	"github.com/75912001/goz/zprotobuf"
	"github.com/golang/protobuf/proto"
)

func registerHandleFun() (ret int) {
	ret = pbFunMgr.Register(zprotobuf.MessageID(pb_account.CMD_LOGIN_MSG), OnLoginMsg, new(pb_account.LoginMsg))
	if 0 != ret {
		return ret
	}
	ret = pbFunMgr.Register(zprotobuf.MessageID(pb_account.CMD_CREATE_ROLE_MSG), OnCreateRoleMsg, new(pb_account.CreateRoleMsg))
	if 0 != ret {
		return ret
	}
	ret = pbFunMgr.Register(zprotobuf.MessageID(pb_account.CMD_LOAD_USER_MSG), OnLoadUserMsg, new(pb_account.LoadUserMsg))
	if 0 != ret {
		return ret
	}
	ret = pbFunMgr.Register(zprotobuf.MessageID(pb_account.CMD_SYS_TIME_MSG), OnSysTimeMsg, new(pb_account.SysTimeMsg))
	if 0 != ret {
		return ret
	}
	return 0
}

func OnLoginMsg(recvProtoHeadBuf []byte, protoMessage *proto.Message, obj interface{}) (ret int) {

	return
}
func OnCreateRoleMsg(recvProtoHeadBuf []byte, protoMessage *proto.Message, obj interface{}) (ret int) {

	return
}
func OnLoadUserMsg(recvProtoHeadBuf []byte, protoMessage *proto.Message, obj interface{}) (ret int) {

	return
}
func OnSysTimeMsg(recvProtoHeadBuf []byte, protoMessage *proto.Message, obj interface{}) (ret int) {

	return
}
