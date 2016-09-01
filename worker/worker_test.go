package worker_test

import(
	"github.com/supersid/iris2/worker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/supersid/iris2/constants"
)

var _ = Describe("Worker", func() {
	Context("When a new worker is initialized", func(){
		It("Should setup the logger", func(){
			_, err := worker.NewWorker("tcp://127.0.0.1:5555", "echo", constants.TEST_ENV)
			Expect(err).To(BeNil())
			Expect(worker.GetLogger()).NotTo(BeNil())
		})
	})

	Context("When ", func(){
		It("should set debugMode on in test and development environments", func(){
			w, _ := worker.NewWorker("tcp://127.0.0.1:5555", "echo", constants.TEST_ENV)
			Expect(w.DebugMode).To(Equal(true))
		})
	})
})


