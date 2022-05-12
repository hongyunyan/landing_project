/*
本文件用于测试 parse.go 的正确性
*/
package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func compare_parse_ans(input_log string, ans string) bool {
	// fmt.Printf("basic log is %v", input_log)
	// fmt.Println()

	parse_result := ParseLog(input_log)
	parse_json, err := json.Marshal(parse_result)
	//fmt.Printf("%+v\n", string(jsons))
	if err != nil {
		fmt.Printf("Json Marshal Failed %v\n", err)
	}

	fmt.Printf("log length is %v\n", len(input_log))

	if ans != string(parse_json) {
		fmt.Printf("Error:parse incorrect! Parse result is %v, ans is %v", string(parse_json), ans)
		return false
	}
	return true
}

// func TestSimpleCase(t *testing.T) {
// 	// [2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] ["TiKV Started"] [ddl_job_id=1]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"INFO",
// 				"SourceFile":"tikv-server.rs",
// 				"LineNumber":13
// 			},
// 			"LogMessageSection":"TiKV Started",
// 			"LogFieldsSection":[
// 				{
// 				"ddl_job_id":"1"
// 				}
// 			]
// 		}
// 	*/
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"INFO\",\"SourceFile\":\"tikv-server.rs\",\"LineNumber\":13},\"LogMessageSection\":\"TiKV Started\",\"LogFieldsSection\":[{\"ddl_job_id\":\"1\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkSimpleCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func BenchmarkSimpleCaseOnlyReadOnce(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
// 	for n := 0; n < b.N; n++ {
// 		ReadStringOnce(log)
// 	}
// }

func BenchmarkNoOtherChar(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>]"
	for n := 0; n < b.N; n++ {
		ParseLog(log)
	}
}

func BenchmarkNoOtherCharOnlyReadOnce(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>]"
	for n := 0; n < b.N; n++ {
		ReadStringOnce(log)
	}
}
func BenchmarkNoOtherCharFindStr(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>]"
	for n := 0; n < b.N; n++ {
		FindStr(log)
	}
}

// func TestBasicCase(t *testing.T) {
// 	/*
// 			log is [2018/12/15 14:20:11.015 +08:00] [FATAL] [panic_hook.rs:45] ["TiKV panic"] [stack="   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace
// 		             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49
// 				1: std::sys_common::backtrace::_print
// 							at /checkout/src/libstd/sys_common/backtrace.rs:71
// 				2: std::panicking::default_hook::{{closure}}
// 							at /checkout/src/libstd/sys_common/backtrace.rs:60
// 							at /checkout/src/libstd/panicking.rs:381"] [error="thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99"]
// 	*/

// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"FATAL",
// 				"SourceFile":"panic_hook.rs",
// 				"LineNumber":45
// 			},
// 			"LogMessageSection":"TiKV panic",
// 			"LogFieldsSection":[
// 				{
// 				"stack":" 0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\n at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\n 1: std::sys_common::backtrace::_print\n at /checkout/src/libstd/sys_common/backtrace.rs:71\n 2: std::panicking::default_hook::{{closure}}\n at /checkout/src/libstd/sys_common/backtrace.rs:60\n at /checkout/src/libstd/panicking.rs:381"
// 				},
// 				{
// 				"error":"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99"
// 				}
// 			]
// 		}
// 	*/
// 	log := "[2018/12/15 14:20:11.015 +08:00] [FATAL] [panic_hook.rs:45] [\"TiKV panic\"] [stack=\"   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\n             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\n   1: std::sys_common::backtrace::_print\n             at /checkout/src/libstd/sys_common/backtrace.rs:71\n   2: std::panicking::default_hook::{{closure}}\n             at /checkout/src/libstd/sys_common/backtrace.rs:60\n             at /checkout/src/libstd/panicking.rs:381\"] [error=\"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99\"]"

// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"FATAL\",\"SourceFile\":\"panic_hook.rs\",\"LineNumber\":45},\"LogMessageSection\":\"TiKV panic\",\"LogFieldsSection\":[{\"stack\":\"   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\\n             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\\n   1: std::sys_common::backtrace::_print\\n             at /checkout/src/libstd/sys_common/backtrace.rs:71\\n   2: std::panicking::default_hook::{{closure}}\\n             at /checkout/src/libstd/sys_common/backtrace.rs:60\\n             at /checkout/src/libstd/panicking.rs:381\"},{\"error\":\"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkBasicCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [FATAL] [panic_hook.rs:45] [\"TiKV panic\"] [stack=\"   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\n             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\n   1: std::sys_common::backtrace::_print\n             at /checkout/src/libstd/sys_common/backtrace.rs:71\n   2: std::panicking::default_hook::{{closure}}\n             at /checkout/src/libstd/sys_common/backtrace.rs:60\n             at /checkout/src/libstd/panicking.rs:381\"] [error=\"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99\"]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func TestComplexCase(t *testing.T) {
// 	// [2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [my_custom_message] [sql="insert into t values (\"]This should not break log parsing!\")"]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"INFO",
// 				"SourceFile":"unknown",
// 				"LineNumber":-1
// 			},
// 			"LogMessageSection":"my_custom_message",
// 			"LogFieldsSection":[
// 				{
// 				"sql":"insert into t values (\"]This should not break log parsing!\")"
// 				}
// 			]
// 		}
// 	*/

// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [my_custom_message] [sql=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"INFO\",\"SourceFile\":\"unknown\",\"LineNumber\":-1},\"LogMessageSection\":\"my_custom_message\",\"LogFieldsSection\":[{\"sql\":\"insert into t values (\\\"]This should not break log parsing!\\\")\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkComplexCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [my_custom_message] [sql=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func TestSlashCase(t *testing.T) {
// 	// [2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] ["abc\"]abc[def]"] ["sql="="insert into t values (\"]This should not break log parsing!\")"]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"INFO",
// 				"SourceFile":"unknown",
// 				"LineNumber":-1
// 			},
// 			"LogMessageSection":"abc\"]abc[def]",
// 			"LogFieldsSection":[
// 				{
// 				"sql=":"insert into t values (\"]This should not break log parsing!\")"
// 				}
// 			]
// 		}
// 	*/
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [\"abc\\\"]abc[def]\"] [\"sql=\"=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"INFO\",\"SourceFile\":\"unknown\",\"LineNumber\":-1},\"LogMessageSection\":\"abc\\\"]abc[def]\",\"LogFieldsSection\":[{\"sql=\":\"insert into t values (\\\"]This should not break log parsing!\\\")\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkSlashCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [\"abc\\\"]abc[def]\"] [\"sql=\"=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func TestFieldsEmptyCase(t *testing.T) {
// 	// [2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] ["TiKV Started"]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"INFO",
// 				"SourceFile":"tikv-server.rs",
// 				"LineNumber":13
// 			},
// 			"LogMessageSection":"TiKV Started",
// 			"LogFieldsSection":[
// 			]
// 		}
// 	*/
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"INFO\",\"SourceFile\":\"tikv-server.rs\",\"LineNumber\":13},\"LogMessageSection\":\"TiKV Started\",\"LogFieldsSection\":[]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkFieldsEmptyCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func TestFieldsMultiCase(t *testing.T) {
// 	//[2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] ["TiKV Started"] ["user name"=foo] [duration=1.345s] [client=192.168.0.123:12345]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"ERROR",
// 				"SourceFile":"tikv-server.rs",
// 				"LineNumber":13
// 			},
// 			"LogMessageSection":"TiKV Started",
// 			"LogFieldsSection":[
// 				{
// 				"user name":"foo"
// 				},
// 				{
// 				"duration":"1.345s"
// 				},
// 				{
// 				"client":"192.168.0.123:12345"
// 				}
// 			]
// 		}
// 	*/
// 	log := "[2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] [\"TiKV Started\"] [\"user name\"=foo] [duration=1.345s] [client=192.168.0.123:12345]"

// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"ERROR\",\"SourceFile\":\"tikv-server.rs\",\"LineNumber\":13},\"LogMessageSection\":\"TiKV Started\",\"LogFieldsSection\":[{\"user name\":\"foo\"},{\"duration\":\"1.345s\"},{\"client\":\"192.168.0.123:12345\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func BenchmarkFieldsMultiCase(b *testing.B) {
// 	log := "[2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] [\"TiKV Started\"] [\"user name\"=foo] [duration=1.345s] [client=192.168.0.123:12345]"
// 	for n := 0; n < b.N; n++ {
// 		ParseLog(log)
// 	}
// }

// func TestMessageWithFinalSlash(t *testing.T) {
// 	//log is [2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] ["TiKV Started\\"] ["user name"=foo]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"ERROR",
// 				"SourceFile":"tikv-server.rs",
// 				"LineNumber":13
// 			},
// 			"LogMessageSection":"TiKV Started\\",
// 			"LogFieldsSection":[
// 				{
// 				"user name":"foo"
// 				}
// 			]
// 		}
// 	*/

// 	log := "[2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] [\"TiKV Started\\\\\"] [\"user name\"=foo]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"ERROR\",\"SourceFile\":\"tikv-server.rs\",\"LineNumber\":13},\"LogMessageSection\":\"TiKV Started\\\\\",\"LogFieldsSection\":[{\"user name\":\"foo\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }

// func TestMessageWithMultiSlash(t *testing.T) {
// 	// log : [2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] ["TiKV\"Started\\"] ["user name\"]\""=foo]
// 	/*
// 		{
// 			"LogHeaderSection":{
// 				"DataTime":"2018/12/15 14:20:11.015 +08:00",
// 				"Level":"ERROR",
// 				"SourceFile":"tikv-server.rs",
// 				"LineNumber":13
// 			},
// 			"LogMessageSection":"TiKV\"Started\\",
// 			"LogFieldsSection":[
// 				{
// 				"user name\"]\"":"foo"
// 				}
// 			]
// 		}
// 	*/

// 	log := "[2018/12/15 14:20:11.015 +08:00] [ERROR] [tikv-server.rs:13] [\"TiKV\\\"Started\\\\\"] [\"user name\\\"]\\\"\"=foo]"
// 	ans := "{\"LogHeaderSection\":{\"DataTime\":\"2018/12/15 14:20:11.015 +08:00\",\"Level\":\"ERROR\",\"SourceFile\":\"tikv-server.rs\",\"LineNumber\":13},\"LogMessageSection\":\"TiKV\\\"Started\\\\\",\"LogFieldsSection\":[{\"user name\\\"]\\\"\":\"foo\"}]}"

// 	if !compare_parse_ans(log, ans) {
// 		t.Error("Compare Failed")
// 	} else {
// 		fmt.Println(" -- Passed --")
// 	}
// }
