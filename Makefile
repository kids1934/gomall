.PHONY: gen-demo-thrift
gen-demo-thrift:
	@cd demo/demo_thrift && cwgo server --type RPC --module github.com/kids1934/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift