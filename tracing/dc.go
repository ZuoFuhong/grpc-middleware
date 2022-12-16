package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func dcReport(ctx context.Context, fullMethod string, req interface{}, st time.Time, err error) {
	var errcode int
	var errmsg string
	if err != nil {
		if rpcErr, ok := status.FromError(err); ok {
			fmt.Println(rpcErr)
			errcode = int(rpcErr.Code())
			errmsg = rpcErr.Message()
		}
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		reqBytes = []byte("request json.Marshal failed")
	}
	// todo: 日志上报字段
	// timestamp
	// server_ip
	// source  调用方
	// project 被调方
	// /go_wallet_manage_svr.go_wallet_manage_svr/ImportWallet
	log.Printf("fullMethod: %s, req: %s, timecost: %v, errcode: %d, errmsg: %s\n", fullMethod, string(reqBytes), time.Since(st), errcode, errmsg)
}
