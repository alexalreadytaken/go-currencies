package services

import (
	"testing"

	"github.com/alexalreadytaken/go-currencies/internal/utils"
	"github.com/stretchr/testify/suite"
)

type DaemonWorkerTestSuite struct {
	suite.Suite
	worker *CurrencyDaemonWorker
}

func TestCurrencyDaemonWorker(t *testing.T) {
	suite.Run(t, &CurrenciesClientTestSuite{})
}

func (s *DaemonWorkerTestSuite) SetupSuite() {
	s.worker = &CurrencyDaemonWorker{
		repo:     &utils.MockCurrenciesRepo{},
		producer: &utils.MockCurrenciesCourseProducer{},
	}
}

func (s *DaemonWorkerTestSuite) TestGetBtcUsdtCourseAndStore() {
	s.worker.producer.(*utils.MockCurrenciesCourseProducer).
		On("GetBtcUsdCourse").Return(100.10, nil)

	s.worker.GetBtcUsdtCourseAndStore()
	s.Equal(100.10, s.worker.latestBtcUsdCourse)
}
