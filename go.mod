module worker

go 1.17

replace (
	github.com/TD-Hackathon-2022/DCoB-Scheduler => ../DCoB-Scheduler
)

require (
	github.com/smartystreets/goconvey v1.7.2
	google.golang.org/protobuf v1.27.1
	github.com/TD-Hackathon-2022/DCoB-Scheduler v0.1.0
)

require (
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
)
