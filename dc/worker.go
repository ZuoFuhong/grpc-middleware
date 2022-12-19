package dc

import (
	"context"
	"errors"
	pb "github.com/ZuoFuhong/grpc-standard-pb/go_datacollector_svr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	NetAddress    = "monica://Production/go_datacollector_svr"
	QueueCapacity = 10000 // 默认工作队列大小
)

var emptyCtx = context.Background()

func newDcStub() pb.GoDatacollectorSvrClient {
	conn, err := grpc.Dial(NetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return pb.NewGoDatacollectorSvrClient(conn)
}

type Worker struct {
	stub   pb.GoDatacollectorSvrClient
	dataq  chan *pb.ReportTraceReq
	closed bool
}

func NewWorker() *Worker {
	w := &Worker{
		stub:  newDcStub(),
		dataq: make(chan *pb.ReportTraceReq, QueueCapacity),
	}
	go w.start()
	return w
}

func (w *Worker) Report(item *pb.ReportTraceReq) error {
	if w.closed {
		return errors.New("worker already close")
	}
	select {
	case w.dataq <- item:
		// 监控 DC 上报量
	default:
		return errors.New("worker queue is full")
	}
	return nil
}

func (w *Worker) start() {
	for !w.closed {
		item := <-w.dataq
		if _, err := w.stub.ReportTrace(emptyCtx, item); err != nil {
			log.Printf("ERROR: [grpc-middleware] report trace failed, %v\n", err)
		}
	}
}

func (w *Worker) Close() {
	w.closed = true
}
