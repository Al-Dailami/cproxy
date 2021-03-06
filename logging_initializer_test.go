package cproxy

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

func TestLoggingInitializerFixture(t *testing.T) {
	gunit.Run(new(LoggingInitializerFixture), t)
}

type LoggingInitializerFixture struct {
	*gunit.Fixture

	client      *TestSocket
	server      *TestSocket
	fakeInner   *TestInitializer
	initializer *LoggingInitializer
}

func (this *LoggingInitializerFixture) Setup() {
	this.client = NewTestSocket()
	this.server = NewTestSocket()
	this.client.address = "1.2.3.4"
	this.client.port = 4321
	this.server.address = "5.6.7.8"
	this.server.port = 8765

	this.fakeInner = NewTestInitializer(true)
	this.initializer = NewLoggingInitializer(this.fakeInner)
	this.initializer.logger = logging.Capture()
}

func (this *LoggingInitializerFixture) TestInnerInitializerCalled() {
	result := this.initializer.Initialize(this.client, this.server)

	this.So(result, should.BeTrue)
	this.So(this.fakeInner.client, should.Equal, this.client)
	this.So(this.fakeInner.server, should.Equal, this.server)
}

func (this *LoggingInitializerFixture) TestLoggingOnFailure() {
	this.fakeInner.success = false

	this.initializer.Initialize(this.client, this.server)

	this.So(this.initializer.logger.Calls, should.Equal, 1)
	this.So(this.initializer.logger.Log.String(), should.EndWith,
		"[INFO] Connection failed [1.2.3.4:4321] -> [5.6.7.8:8765]\n")
}

func (this *LoggingInitializerFixture) TestLoggingOnSuccess() {
	this.initializer.Initialize(this.client, this.server)

	this.So(this.initializer.logger.Calls, should.Equal, 1)
	this.So(this.initializer.logger.Log.String(), should.EndWith,
		"[INFO] Established connection [1.2.3.4:4321] -> [5.6.7.8:8765]\n")
}
