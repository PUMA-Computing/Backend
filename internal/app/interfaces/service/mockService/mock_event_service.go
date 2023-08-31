// mocks/mock_event_service.go
package mocks

import (
	"Backend/internal/app/domain/event"
	"github.com/golang/mock/gomock"
	"reflect"
)

type MockEventService struct {
	ctrl     *gomock.Controller
	recorder *MockEventServiceMockRecorder
}

func NewMockEventService(ctrl *gomock.Controller) *MockEventService {
	mock := &MockEventService{ctrl: ctrl}
	mock.recorder = &MockEventServiceMockRecorder{mock}
	return mock
}

func (m *MockEventService) EXPECT() *MockEventServiceMockRecorder {
	return m.recorder
}

type MockEventServiceMockRecorder struct {
	mock *MockEventService
}

// Implement methods from the EventService interface

func (r *MockEventServiceMockRecorder) CreateEvent(event *event.Event) *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "CreateEvent", reflect.TypeOf((*MockEventService)(nil)), event)
}

func (r *MockEventServiceMockRecorder) UpdateEvent(eventID string, updatedEvent *event.Event) *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "UpdateEvent", reflect.TypeOf((*MockEventService)(nil)), eventID, updatedEvent)
}

func (r *MockEventServiceMockRecorder) DeleteEvent(eventID string) *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "DeleteEvent", reflect.TypeOf((*MockEventService)(nil)), eventID)
}

func (r *MockEventServiceMockRecorder) GetEvent() *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetEvent", reflect.TypeOf((*MockEventService)(nil)))
}

func (r *MockEventServiceMockRecorder) GetEventUsers(eventID string) *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetEventUsers", reflect.TypeOf((*MockEventService)(nil)), eventID)
}

func (r *MockEventServiceMockRecorder) RegisterUserForEvent(UserID, eventID string) *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "RegisterUserForEvent", reflect.TypeOf((*MockEventService)(nil)), UserID, eventID)
}
