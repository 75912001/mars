package main

import (
	"github.com/75912001/goz/zhttp"
	"github.com/75912001/goz/zprotobuf"
	"github.com/75912001/goz/ztcp"
	"github.com/75912001/goz/zutility"
)

var ini zutility.Ini
var log zutility.Log

var tcpServer ztcp.Server
var httpServer zhttp.Server

var pbFunMgr zprotobuf.PbFunMgr
