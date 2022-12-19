package dc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	meta "github.com/ZuoFuhong/grpc-middleware/metadata"
	_ "github.com/ZuoFuhong/grpc-naming-monica"
	pb "github.com/ZuoFuhong/grpc-standard-pb/go_datacollector_svr"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var workerIns *Worker

func getWorker() *Worker {
	once.Do(func() {
		workerIns = NewWorker()
	})
	return workerIns
}

// Report 上报链路调用
func Report(ctx context.Context, fullMethod string, req interface{}, st time.Time, err error) {
	var errcode int
	var errmsg string
	if err != nil {
		if rpcErr, ok := status.FromError(err); ok {
			errcode = int(rpcErr.Code())
			errmsg = rpcErr.Message()
		}
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		reqBytes = []byte("request json.Marshal failed")
	}
	var traceId string
	if ctxTids := metadata.ValueFromIncomingContext(ctx, meta.TraceId); len(ctxTids) > 0 {
		traceId = ctxTids[0]
	}

	project, cmd := extractTarget(fullMethod)
	binName := getBinName()
	serverIP := getLocalIp()
	if err := getWorker().Report(&pb.ReportTraceReq{
		TraceId:   traceId,
		Cmd:       cmd,
		Project:   project,
		Source:    binName,
		ServerIp:  serverIP,
		Errcode:   int64(errcode),
		Errmsg:    errmsg,
		Timestamp: st.Unix(),
		Timecost:  time.Since(st).Milliseconds(),
		Reqbody:   string(reqBytes),
	}); err != nil {
		log.Printf("ERROR: [grpc-middleware] worker report failed, %v\n", err)
	}
}

// extractTarget 提取目标命令字
func extractTarget(fullMethod string) (string, string) {
	arr := strings.Split(strings.TrimLeft(fullMethod, "/"), "/")
	if len(arr) == 2 {
		sarr := strings.Split(arr[0], ".")
		if len(sarr) == 2 {
			return sarr[1], arr[1]
		}
	}
	return "", ""
}

// getBinName 获取进程名字
func getBinName() string {
	binName := os.Args[0]
	binNames := strings.Split(os.Args[0], "/")
	if len(binNames) > 0 {
		binName = binNames[len(binNames)-1]
	}
	return binName
}

// getLocalIp 获取本地机器 IP
func getLocalIp() string {
	for i := 0; i < 10; i++ {
		ip, err := getInterIP(fmt.Sprintf("eth%d", i))
		if err == nil && ip != "" {
			return ip
		}
	}
	return "127.0.0.1"
}

func getInterIP(name string) (string, error) {
	eths, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	regex := regexp.MustCompile("[0-9.]+")
	for _, eth := range eths {
		if eth.Name == name {
			addrs, err := eth.Addrs()
			if err != nil {
				return "", err
			}
			if 0 == len(addrs) {
				return "", errors.New("empty interface")
			}
			fileds := regex.FindAllString(addrs[0].String(), -1)
			return fileds[0], nil
		}
	}
	return "", errors.New("invalid interface")
}
