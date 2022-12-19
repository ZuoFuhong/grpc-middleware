package dc

import (
	"context"
	pb "github.com/ZuoFuhong/grpc-standard-pb/go_datacollector_svr"
	"testing"
	"time"
)

func Test_DcReport(t *testing.T) {
	start := time.Now()
	Report(context.Background(), "/go_datacollector_svr.go_datacollector_svr/ReportTrace", new(pb.ReportTraceReq), start, nil)
	time.Sleep(time.Minute)
}
