package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CancelManagerTestSuite struct {
	suite.Suite
}

func TestCancelManagerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(CancelManagerTestSuite))
}

func (s *CancelManagerTestSuite) Test_Add_IDEqual_CancelsContext_Returns_Context() {
	cm := NewCancelManager()
	ctx := cm.Add("id1", 1)
	cm.Add("id1", 2)

L:
	for {
		select {
		case <-time.After(time.Second * 5):
			s.Fail("Timeout")
			return
		case <-ctx.Done():
			break L
		}
	}
}

func (s *CancelManagerTestSuite) Test_Add_IDNotExist_Returns_Context() {

	cm := NewCancelManager()
	firstCtx := cm.Add("id1", 1)
	s.NotNil(firstCtx)

	s.Require().Contains(cm.v, "id1")
	s.Equal(cm.v["id1"].ReqID, int64(1))
}

func (s *CancelManagerTestSuite) Test_Delete_IDEqual_ReqIDNotEqual_DoesNothing() {
	cm := NewCancelManager()
	cm.Add("id1", 1)

	s.Require().Len(cm.v, 1)

	cm.Delete("id1", 2)
	s.Require().Len(cm.v, 1)
	s.Require().Contains(cm.v, "id1")
	s.Equal(cm.v["id1"].ReqID, int64(1))

}

func (s *CancelManagerTestSuite) Test_Delete_IDEqual_ReqIDEqual_CallsCancel_RemovesFromMemory() {
	cm := NewCancelManager()
	ctx := cm.Add("id1", 1)

	s.Require().Len(cm.v, 1)

	cm.Delete("id1", 1)
	s.Require().Len(cm.v, 0)

L:
	for {
		select {
		case <-time.After(time.Second * 5):
			s.Fail("Timeout")
			return
		case <-ctx.Done():
			break L
		}
	}
}
