/*
本文件用于来做 benchmark 测试
*/

package main

import (
	"testing"
)

func BenchmarkComplexCase(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV\\\"Started\\\\\"] [ddl_job_id=1] [stack=\"   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\n             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\n   1: std::sys_common::backtrace::_print\n             at /checkout/src/libstd/sys_common/backtrace.rs:71\n   2: std::panicking::default_hook::{{closure}}\n             at /checkout/src/libstd/sys_common/backtrace.rs:60\n             at /checkout/src/libstd/panicking.rs:381\"] [error=\"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99\"] [\"sql=\"=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
	for n := 0; n < b.N; n++ {
		ParseLog(log)
	}
}

func BenchmarkComplexCaseOnlyReadOnce(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV\\\"Started\\\\\"] [ddl_job_id=1] [stack=\"   0: std::sys::imp::backtrace::tracing::imp::unwind_backtrace\n             at /checkout/src/libstd/sys/unix/backtrace/tracing/gcc_s.rs:49\n   1: std::sys_common::backtrace::_print\n             at /checkout/src/libstd/sys_common/backtrace.rs:71\n   2: std::panicking::default_hook::{{closure}}\n             at /checkout/src/libstd/sys_common/backtrace.rs:60\n             at /checkout/src/libstd/panicking.rs:381\"] [error=\"thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 99\"] [\"sql=\"=\"insert into t values (\\\"]This should not break log parsing!\\\")\"]"
	for n := 0; n < b.N; n++ {
		ReadStringOnce(log)
	}
}

func BenchmarkSimpleCase(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
	for n := 0; n < b.N; n++ {
		ParseLog(log)
	}
}

func BenchmarkSimpleCaseOnlyReadOnce(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
	for n := 0; n < b.N; n++ {
		ReadStringOnce(log)
	}
}

func BenchmarkNoOtherChar(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [TiKV Started]"
	for n := 0; n < b.N; n++ {
		ParseLog(log)
	}
}

func BenchmarkNoOtherCharOnlyReadOnce(b *testing.B) {
	log := "[2018/12/15 14:20:11.015 +08:00] [INFO] [<unknown>] [TiKV Started]"
	for n := 0; n < b.N; n++ {
		ReadStringOnce(log)
	}
}
