package broker_test

import(
	"github.com/supersid/iris2/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Broker", func() {
	Context("When the client is new", func(){
		It("Should add the client request to the corressponding service", func(){

		})
	})

	Context("When the client already exists", func(){
		It("Should create a new service and add the client request to that service", func(){

		})
	})

	Context("When a new worker comes online", func(){
		BeforeEach(func(){
			//broker := broker.Start("tcp://*:5555")
		})

		AfterEach(func(){

		})
		It("Should receive a message with command WORKER_READY", func(){

		})
	})
})
