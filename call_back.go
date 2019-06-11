package main

import (
	"github.com/75912001/goz/zprotobuf"
	"github.com/75912001/goz/ztcp"
)

func onInit() (ret int) {
	log.Trace("onInit")
	pbFunMgr.Init(&log)
	registerHandleFun()
	return 0
}

func onFini() (ret int) {
	log.Trace("onFini")
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//客户端相关的回调函数
func onCliConn(realPeerConn *ztcp.PeerConn) (ret int) {
	log.Trace("onCliConn")
	user := GuserMgr.Add(realPeerConn)
	user.PeerConn = realPeerConn
	return 0
}

func onCliConnClosed(realPeerConn *ztcp.PeerConn) (ret int) {
	log.Trace("onCliConnClosed")

	user := GuserMgr.Find(realPeerConn)
	GuserMgr.Del(user.UID)
	GuserMgr.Del(realPeerConn)

	return 0
}

func onParseProtoHead(peerConn *ztcp.PeerConn, length int) (ret int) {

	if length < GProtoHeadLength { //长度不足一个包头的长度大小
		return 0
	}

	packetLength := int(parseProtoHeadPacketLength(peerConn.Buf))

	if int(packetLength) < GProtoHeadLength {
		log.Error("PacketLength:", packetLength)
		return -1
	}
	if tcpServer.PacketLengthMax <= uint32(length) {
		log.Error("PacketLengthMax:", tcpServer.PacketLengthMax, length)
		return -1
	}

	if length < int(packetLength) {
		return 0
	}

	return packetLength
}

func onCliPacket(peerConn *ztcp.PeerConn, recvBuf []byte) (ret int) {
	packetLength, sessionID, messageID, resultID, userID := parseProtoHead(recvBuf)
	log.Trace(packetLength, messageID, sessionID, userID, resultID)

	user := GuserMgr.Find(peerConn)

	if nil == user {
		log.Crit("")
		return
	}

	return pbFunMgr.OnRecv(zprotobuf.MessageID(messageID), recvBuf[:GProtoHeadLength], recvBuf[GProtoHeadLength:packetLength], user)

}
