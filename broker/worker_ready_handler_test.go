package broker_test

//import(
//	"github.com/supersid/iris2/broker"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/supersid/iris2/worker"
//	"github.com/supersid/iris2/constants"
//	"os"
//
//	"time"
//)

//var _ = Describe("WorkerReadyHandler", func() {
//	var b *broker.Broker
//	const brokerUrl = "tcp://127.0.0.1:5555"
//
//	BeforeEach(func(){
//
//		_b, err := broker.NewBroker("tcp://*:5555", constants.TEST_ENV)
//		b = _b
//		b.Socket.Bind("tcp://*:5555")
//		defer b.Socket.Close()
//		Expect(err).NotTo(HaveOccurred())
//		os.Setenv("IRIS_ENV", constants.TEST_ENV)
//	})
//
//	Context("When a new worker is added", func(){
//		It("should create a new service", func(){
//			b.Process()
//			w, e := worker.NewWorker(brokerUrl, "echo", constants.TEST_ENV)
//			Expect(e).NotTo(HaveOccurred())
//			e = w.Socket.Connect(brokerUrl)
//			Expect(e).NotTo(HaveOccurred())
//			w.Socket.Connect(brokerUrl)
//			messageChannel := make(chan []string)
//			go w.Process(messageChannel)
//			time.Sleep(1 * time.Second)
//			Expect(len(b.Services)).To(Equal(1))
//		})
//
//		It("should add a service worker with the same identity to the service", func(){
//
//		})
//
//		It("should add the service to the broker's list of services", func(){
//
//		})
//	})
//
//	Context("When the client already exists", func(){
//		It("Should create a new service and add the client request to that service", func(){
//
//		})
//	})
//
//	Context("When a new worker comes online", func(){
//		BeforeEach(func(){
//			//broker := broker.Start("tcp://*:5555")
//		})
//
//		AfterEach(func(){
//
//		})
//		It("Should receive a message with command WORKER_READY", func(){
//
//		})
//	})
//})
