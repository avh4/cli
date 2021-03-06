package api_test

import (
	. "cf/api"
	"cf/models"
	"cf/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	testapi "testhelpers/api"
	testconfig "testhelpers/configuration"
	testnet "testhelpers/net"
	"time"
)

var _ = Describe("AppInstancesRepo", func() {
	It("TestAppInstancesGetInstances", func() {
		ts, handler, repo := createAppInstancesRepo([]testnet.TestRequest{
			appInstancesRequest,
			appStatsRequest,
		})
		defer ts.Close()
		appGuid := "my-cool-app-guid"

		instances, err := repo.GetInstances(appGuid)
		Expect(handler.AllRequestsCalled()).To(BeTrue())
		Expect(err.IsNotSuccessful()).To(BeFalse())

		Expect(len(instances)).To(Equal(2))

		Expect(instances[0].State).To(Equal(models.InstanceRunning))
		Expect(instances[1].State).To(Equal(models.InstanceStarting))

		instance0 := instances[0]
		Expect(instance0.Since).To(Equal(time.Unix(1379522342, 0)))
		Expect(instance0.DiskQuota).To(Equal(uint64(1073741824)))
		Expect(instance0.DiskUsage).To(Equal(uint64(56037376)))
		Expect(instance0.MemQuota).To(Equal(uint64(67108864)))
		Expect(instance0.MemUsage).To(Equal(uint64(19218432)))
		Expect(instance0.CpuUsage).To(Equal(3.659571249238058e-05))
	})
})

var appStatsRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method: "GET",
	Path:   "/v2/apps/my-cool-app-guid/stats",
	Response: testnet.TestResponse{Status: http.StatusOK, Body: `
{
  "1":{
    "stats": {
        "disk_quota": 10000,
        "mem_quota": 1024,
        "usage": {
            "cpu": 0.3,
            "disk": 10000,
            "mem": 1024
        }
    }
  },
  "0":{
    "stats": {
        "disk_quota": 1073741824,
        "mem_quota": 67108864,
        "usage": {
            "cpu": 3.659571249238058e-05,
            "disk": 56037376,
            "mem": 19218432
        }
    }
  }
}`}})

var appInstancesRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method: "GET",
	Path:   "/v2/apps/my-cool-app-guid/instances",
	Response: testnet.TestResponse{Status: http.StatusOK, Body: `
{
  "1": {
    "state": "STARTING",
    "since": 1379522342.6783738
  },
  "0": {
    "state": "RUNNING",
    "since": 1379522342.6783738
  }
}`}})

func createAppInstancesRepo(requests []testnet.TestRequest) (ts *httptest.Server, handler *testnet.TestHandler, repo AppInstancesRepository) {
	ts, handler = testnet.NewTLSServer(requests)
	space := models.SpaceFields{}
	space.Guid = "my-space-guid"
	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetApiEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway()
	repo = NewCloudControllerAppInstancesRepository(configRepo, gateway)
	return
}
