package main

import (
	"mytest/arraytest"
	"mytest/cgotest"
	"mytest/channeltest"
	"mytest/closuretest"
	"mytest/contexttest"
	"mytest/goroutionstest"
	"mytest/interfacetest"
	"mytest/panictest"
	producerconsumertest "mytest/producer_consumer_test"
	"mytest/serverqpslimittest"
	"mytest/servertest"
	"mytest/variable_parameter_test"
)

func ArrayTests() {
	arraytest.ArrayTest()
	arraytest.TestTopSort()
	arraytest.StringTest()
	arraytest.TwoDimensionalByteTest()
	arraytest.MakeTest()
	arraytest.QueueStackTest()
	arraytest.QueueTest()
	arraytest.MapAsSetTest()
}

func ChannelTests1() {
	channeltest.ChannelTest()
	channeltest.PrintCatFishDog()
	channeltest.PrintWordAndNums()
}

func ChannelTests2() {
	channeltest.TestCounter()
	channeltest.TestCounterV2()
	channeltest.SelectTest()
}

func ClosureTests() {
	closuretest.Test1()
	closuretest.Test2()
	closuretest.TestPage()
}

func GoRoutionsTests() {
	//goroutionstest.TestPageV2()
	goroutionstest.DuTest()
	goroutionstest.DuTestV2()
	goroutionstest.DuTestV3()
}

func GoRoutionsCacheTests() {
	//goroutionstest.CacheTest()
	goroutionstest.CacheTestV2()
}

func PanicTests() {
	panictest.PanicTest()
}

func SvrTests() {
	//servertest.TestHttpSvrV1()
	//servertest.TestTcpSvrV1()
	servertest.TestTcpSvrV2()
}

func CgoTests() {
	cgotest.GoSumTest(1, 2)
	cgotest.GoCppTest()
}

func PublisherTests() {
	producerconsumertest.PublisherTest()
}

func ContextTests() {
	contexttest.ContextTest()
	contexttest.ContextTest2()
}

func SvrQpsLimiteTests() {
	//serverqpslimittest.QpsLimiteTest()
	serverqpslimittest.SvrQpsLimiteTest()
}

func VariableParameterTestS() {
	variable_parameter_test.VariableParameterTest()
	variable_parameter_test.OptionTest()
}

func InterfaceConvertTests() {
	interfacetest.InterfaceConvertTest()
}

func main() {
	ArrayTests()
	//ChannelTests1()
	//ChannelTests2()
	//ClosureTests()
	//PanicTests()
	//GoRoutionsTests()
	//SvrTests()
	//GoRoutionsCacheTests()
	//CgoTests()
	//PublisherTests()
	//ContextTests()
	//SvrQpsLimiteTests()
	//VariableParameterTestS()
	//InterfaceConvertTests()
}
