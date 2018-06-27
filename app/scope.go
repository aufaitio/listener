package app

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
)

// RequestScope contains the application-specific information that are carried around in a request.
type RequestScope interface {
	Logger
	// RequestID returns the ID of the current request
	RequestID() string
	// Now returns the timestamp representing the time when the request is being processed
	Now() time.Time
	DB() *mongo.Database
	Context() context.Context
}

type requestScope struct {
	Logger                    // the logger tagged with the current request information
	now       time.Time       // the time when the request is being processed
	requestID string          // an ID identifying one or multiple correlated HTTP requests
	db        *mongo.Database // the mongo db client
	request   *http.Request
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}

func (rs *requestScope) DB() *mongo.Database {
	return rs.db
}

func (rs *requestScope) Context() context.Context {
	return rs.request.Context()
}

// newRequestScope creates a new RequestScope with the current request information.
func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request, db *mongo.Database) RequestScope {
	l := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		l.SetField("RequestID", requestID)
	}

	return &requestScope{
		Logger:    l,
		now:       now,
		requestID: requestID,
		db:        db,
		request:   request,
	}
}
