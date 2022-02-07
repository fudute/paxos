package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fudute/paxos"
	"github.com/fudute/paxos/config"
	pb "github.com/fudute/paxos/protoc"
	"github.com/fudute/paxos/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UUID struct {
	id int64
}

func NewUUIDService() paxos.IDService {
	return &UUID{id: 1}
}

func (u *UUID) Next() int64 {
	return atomic.AddInt64(&u.id, 1)
}

var N = flag.Int("n", 10, "times")
var wg sync.WaitGroup

func main() {
	log.SetFlags(0)
	flag.Parse()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	log.Printf("load config: %#v", conf)
	var nAcceptor = len(conf.Cluster.Nodes)

	var quorum = (nAcceptor + 1) / 2
	log.Println("quorum is ", quorum)

	idService := NewUUIDService()

	for _, node := range conf.Cluster.Nodes {
		go func(node *config.Node) {
			s := grpc.NewServer()

			path := fmt.Sprint(os.TempDir()+"/test_", node.Name)
			store, err := store.NewLevelDBLogStore(path)
			if err != nil {
				log.Fatal(err)
			}
			service := paxos.NewPaxosService(node.Name, node.Addr, quorum, conf, idService, store)
			s.RegisterService(&pb.Proposer_ServiceDesc, service)
			s.RegisterService(&pb.Learner_ServiceDesc, service)
			s.RegisterService(&pb.Acceptor_ServiceDesc, service)
			lis, err := net.Listen("tcp", node.Addr)
			if err != nil {
				log.Fatal("listen failed", err)
			}
			log.Println("proposer listen at ", node.Addr)
			log.Fatal(s.Serve(lis))
		}(node)
	}

	// client
	<-time.After(time.Second)

	start := time.Now()
	cc1, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cc2, err := grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)

	go func() {
		for i := 0; i < *N; i++ {
			ch <- fmt.Sprint(i + 1)
		}
		close(ch)
	}()

	wg.Add(2)
	go send(pb.NewProposerClient(cc1), ch)
	go send(pb.NewProposerClient(cc2), ch)

	wg.Wait()

	fmt.Println("qps: ", float64(*N)/time.Since(start).Seconds())
}

func send(cli pb.ProposerClient, in <-chan string) {
	defer wg.Done()
	var err error
	for {
		value, ok := <-in
		if !ok {
			return
		}

		_, err = cli.Propose(context.Background(), &pb.ProposeRequest{
			ProposalValue: value,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
