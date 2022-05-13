/*
parse.go 用于解析标准的单条 Unified Log Format 成 标准的json格式，目前只支持解析合法的单条 log 条目，不进行错误处理
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

type LogHeaderSection struct {
	DataTime   string
	Level      string
	SourceFile string
	LineNumber int
}

type UnifiedLog struct {
	LogHeaderSection  LogHeaderSection
	LogMessageSection string
	LogFieldsSection  []map[string]string
}

// 单纯从头遍历 log 字符串
func ReadStringOnce(log string) {
	var pre_index int
	var pre_char rune
	for index, char := range log {
		pre_index = index
		pre_char = char
	}
	_ = pre_index
	_ = pre_char

}

// 解析log的 LogHeaderSection 部分，返回后续部分的 begin_index
func ParseLogHeaderSection(log string, unified_log *UnifiedLog) int {
	index := strings.Index(log, "]")
	unified_log.LogHeaderSection.DataTime = log[1:index]
	begin_index := index + 2

	index = strings.Index(log[begin_index:], "]")
	unified_log.LogHeaderSection.Level = log[begin_index+1 : begin_index+index]
	begin_index += index + 2

	index = strings.Index(log[begin_index:], "]")
	value := log[begin_index+1 : begin_index+index]

	begin_index += index + 2

	if value == "<unknown>" {
		unified_log.LogHeaderSection.SourceFile = "unknown"
	} else {
		index = strings.Index(value, ":")
		unified_log.LogHeaderSection.SourceFile = value[:index]
		number, err := strconv.Atoi(value[index+1:])
		if err != nil {
			fmt.Printf("unmatched number is %v, content is %v", number, value)
			fmt.Println()
		}
		unified_log.LogHeaderSection.LineNumber = number
	}

	return begin_index
}

// 解析带有引号的字符串，必要时候需要进行转义处理，返回结束的 index 和 最终的 字符串
func ParseQuotedJsonString(log string) (int, string) {
	var content strings.Builder
	is_slash := false
	begin_index := 0

	for index, item := range log {
		if item == '"' && !is_slash {
			// 如果是前面没有有效的\的"， 就说明到了结尾
			if content.Len() == 0 { // 如果content是空的，就说明中途没有 slash，所以也没有重写过，直接拿对应的字符串就可以
				return index, log[:index]
			} else {
				content.WriteString(log[begin_index:index])
				return index, content.String()
			}
		} else if item == '\\' && !is_slash {
			// 如果前面的不是 slash， 那就说明这是个转义的，所以要把这个转义符号去除
			content.WriteString(log[begin_index:index])
			is_slash = true
			begin_index = index + 1
		} else {
			// 其他情况都只需要把is_slash 设为 false 就可以
			is_slash = false
		}
	}

	return -1, ""
}

// 解析log的 LogMessageSection 部分，返回后续部分的 begin_index
func ParseLogMessageSection(log string, unified_log *UnifiedLog) int {
	//先判断是否是带双引号的 JsonString 字符串
	if log[1] == '"' {
		end_index, content := ParseQuotedJsonString(log[2:])
		unified_log.LogMessageSection = content
		return end_index + 5
	} else {
		// 直接查找 ] 结尾标志符
		index := strings.Index(log, "]")
		unified_log.LogMessageSection = log[1:index]
		return index + 2
	}
}

// 解析log的 LogMessageField 的单个Field, 返回后续部分的 begin_index
func ParseLogField(log string, unified_log *UnifiedLog) int {
	// 先解析 key， 然后解析 value
	var key string
	var value string

	var value_part_begin_index int
	var next_begin_index int
	var index int

	// Search for Key
	if log[1] == '"' {
		index, key = ParseQuotedJsonString(log[2:])
		value_part_begin_index = index + 4
	} else {
		// 直接找 =
		index := strings.Index(log, "=")
		key = log[1:index]
		value_part_begin_index = index + 1
	}

	// Search for Value
	if log[value_part_begin_index] == '"' {
		index, value = ParseQuotedJsonString(log[value_part_begin_index+1:])
		next_begin_index = index + value_part_begin_index + 4
	} else {
		// 直接找 ]
		index := strings.Index(log[value_part_begin_index:], "]")
		value = log[value_part_begin_index : index+value_part_begin_index]
		next_begin_index = index + value_part_begin_index + 2
	}

	key_value_pair := make(map[string]string)
	key_value_pair[key] = value
	unified_log.LogFieldsSection = append(unified_log.LogFieldsSection, key_value_pair)

	return next_begin_index
}

// 这边用一个从头遍历的方式来进行处理Parse，如果遇到转义，就进行判断，只有 \ 和 " 是需要 额外的转义符号
func ParseLog(log string) UnifiedLog {
	var unified_log UnifiedLog = UnifiedLog{
		LogHeaderSection:  LogHeaderSection{DataTime: "", Level: "", SourceFile: "", LineNumber: -1},
		LogMessageSection: "",
		LogFieldsSection:  []map[string]string{},
	}

	begin_index := ParseLogHeaderSection(log, &unified_log)
	log = log[begin_index:]

	begin_index = ParseLogMessageSection(log, &unified_log)

	// 开始解析 N fields
	for begin_index < len(log) {
		log = log[begin_index:]

		begin_index = ParseLogField(log, &unified_log)
	}

	return unified_log
}

func main() {
	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	//采样cpu运行状态
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	i := 1
	log_test := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
	for i < 1000000 {
		i += 1
		ParseLog(log_test)
	}

}
