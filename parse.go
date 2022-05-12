/* parse.go 用于解析标准的单条 Unified Log Format 成 标准的json格式，目前只支持解析合法的单条 log 条目，不进行错误处理
目前现设定为string 输入， json 输出
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

func FindStr(log string) {
	index := strings.Index(log, "]")
	if index != -1 && index != len(log)-1 {
		value := log[:index]
		index = strings.Index(log[index+1:], "]")
		_ = index
		_ = value
	}
}

func ChangeToRuneArray(log string) {
	rune_array := []rune(log)
	size := len(rune_array)
	i := 0
	var char rune
	for i < size {
		char = rune_array[i]
		i += 1
	}
	_ = char
}

func ParseLogHeaderSection(log string, unified_log *UnifiedLog) int {
	// Choice 1 :先将前面三个part切分出来
	index := strings.Index(log, "]")
	unified_log.LogHeaderSection.DataTime = log[1:index]
	begin_index := index + 2

	index = strings.Index(log[begin_index:], "]")
	unified_log.LogHeaderSection.Level = log[begin_index+1 : begin_index+index]
	begin_index += index + 2

	index = strings.Index(log[begin_index:], "]")
	value := log[begin_index+1 : begin_index+index]

	begin_index += index + 2
	//log = log[begin_index:]

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

/*
// 尝试先把 log 转成 rune 的 array 模式，然后用 index 遍历的方式来做 parse，benchmark 测下来还不如原来的方法
func ParseLogWithRuneArray(log string) UnifiedLog {
	var unified_log UnifiedLog = UnifiedLog{
		LogHeaderSection:  LogHeaderSection{DataTime: "", Level: "", SourceFile: "", LineNumber: -1},
		LogMessageSection: "",
		LogFieldsSection:  []map[string]string{},
	}

	// 用直接找 index 的方式处理前面的 Header 部分
	new_index := ParseLogHeaderSection(log, &unified_log)

	if new_index > len(log) {
		return unified_log
	}
	log = log[new_index:]

	// 用 part 来标记目前 search 的是哪个部分的字段， 1/2/3/4/5 分别对应了 DataTime/Level/SourceFile/LineNumber/Message/Fields, 其中 5 可以在最后反复出现
	part := 3 // 从0 开始的话就是不用前面的直接parse 前三个部分

	if log[1] != '"' {
		//message 部分，如果不是"开头的，那就直接找后续的]
		message_end_index := strings.Index(log, "]")
		unified_log.LogMessageSection = log[1:message_end_index]
		if message_end_index+2 < len(log) {
			log = log[message_end_index+2:]
			part = 4
		} else { //如果没有 field 部分，就直接对 log 置空
			return unified_log
		}
	}

	log_rune_array := []rune(log)
	size := len(log_rune_array)

	var content strings.Builder
	var key_content strings.Builder
	var value_content strings.Builder

	var key string

	is_slash := false // 用于 part 4/5 中区分 \

	is_search_for_key := true

	//直接从 cur_index = 2 开始，因为0是[,1是”
	cur_index := 0
	begin_index := -1
	search_state := -1

	if part == 3 {
		cur_index = 2
		begin_index = 2
		search_state = 1
		part = 4
	}

	for cur_index < size {
		if search_state == -1 {
			// 进入识别状态
			if log_rune_array[cur_index] == '[' {
				if part < 5 {
					part += 1
				}
				begin_index = cur_index + 1
				search_state = 0
				is_search_for_key = true
			}
			cur_index += 1
		} else {
			if part == 4 {
				//因为前面已经处理过了，所以这边的message肯定是带双引号的，所以不用额外判断了
				//如果里面有 slash 出现过，就必须要拿出来重新写过，否则就不需要
				if log_rune_array[cur_index] == '"' && !is_slash {
					// 如果是前面没有有效的\的"， 就说明到了结尾
					if content.Len() == 0 { // 如果content是空的，就说明中途没有 slash，所以也没有重写过，直接拿对应的字符串就可以
						unified_log.LogMessageSection = log[begin_index:cur_index]
					} else {
						content.WriteString(log[begin_index:cur_index])
						unified_log.LogMessageSection = content.String()
					}
					search_state = -1
				} else if log_rune_array[cur_index] == '\\' && !is_slash {
					// 如果前面的不是 slash， 那就说明这是个转义的，要把前面这段先写入content
					content.WriteString(log[begin_index:cur_index])
					is_slash = true
					begin_index = cur_index + 1
				} else {
					// 其他情况都只需要把is_slash 设为 false 就可以
					is_slash = false
				}
			} else if part == 5 {
				if search_state == 0 {
					if log_rune_array[cur_index] == '=' { // 这边特殊处理一下 =
						cur_index += 1
						continue
					} else if log_rune_array[cur_index] == '"' {
						begin_index = cur_index + 1
						is_slash = false
						search_state = 1
					} else {
						// 如果没有引号，那就直接搜索就可
						//["TiKV Started"] [ddl_job_id=1]
						rest_log := log[cur_index:]
						if is_search_for_key {
							end_index := strings.Index(rest_log, "=")
							key = rest_log[:end_index]

							cur_index += end_index
							search_state = 0
							is_search_for_key = false
						} else {
							end_index := strings.Index(rest_log, "]")

							key_value_pair := make(map[string]string)
							key_value_pair[key] = rest_log[:end_index]
							unified_log.LogFieldsSection = append(unified_log.LogFieldsSection, key_value_pair)

							search_state = -1
							cur_index += end_index
						}
					}
				} else {
					if log_rune_array[cur_index] == '"' && !is_slash {
						// 如果遇到 前面没有转义符号的 "
						if is_search_for_key {
							key_content.WriteString(log[begin_index:cur_index])
							key = key_content.String()
							is_search_for_key = false
							search_state = 0
						} else {
							key_value_pair := make(map[string]string)
							if value_content.Len() == 0 {
								key_value_pair[key] = log[begin_index:cur_index]
							} else {
								value_content.WriteString(log[begin_index:cur_index])
								key_value_pair[key] = value_content.String()
							}
							unified_log.LogFieldsSection = append(unified_log.LogFieldsSection, key_value_pair)
							search_state = -1
							key_content.Reset()
							value_content.Reset()
						}
					} else if log_rune_array[cur_index] == '\\' && !is_slash {
						// 如果前面的不是 slash， 那就说明这是个转义的，要把前面这段先写入content
						if is_search_for_key {
							key_content.WriteString(log[begin_index:cur_index])
						} else {
							value_content.WriteString(log[begin_index:cur_index])
						}

						is_slash = true
						begin_index = cur_index + 1
					} else {
						// 其他情况都只需要把is_slash 设为 false 就可以
						is_slash = false
					}
				}
			}
			cur_index += 1
		}
	}
	return unified_log
}
*/

// 这边用一个从头遍历的方式来进行处理Parse，如果遇到转义，就进行判断，只有 \ 和 " 是需要 额外的转义符号
func ParseLog(log string) UnifiedLog {
	var unified_log UnifiedLog = UnifiedLog{
		LogHeaderSection:  LogHeaderSection{DataTime: "", Level: "", SourceFile: "", LineNumber: -1},
		LogMessageSection: "",
		LogFieldsSection:  []map[string]string{},
	}

	// 用直接找 index 的方式处理前面的 Header 部分
	new_index := ParseLogHeaderSection(log, &unified_log)

	if new_index > len(log) {
		return unified_log
	}
	log = log[new_index:]

	// 用 part 来标记目前 search 的是哪个部分的字段， 1/2/3/4/5 分别对应了 DataTime/Level/SourceFile/LineNumber/Message/Fields, 其中 5 可以在最后反复出现
	part := 3 // 从0 开始的话就是不用前面的直接parse 前三个部分

	if log[1] != '"' {
		//message 部分，如果不是"开头的，那就直接找后续的]
		message_end_index := strings.Index(log, "]")
		unified_log.LogMessageSection = log[1:message_end_index]
		if message_end_index+2 < len(log) {
			log = log[message_end_index+2:]
			part = 4
		} else { //如果没有 field 部分，就直接对 log 置空
			return unified_log
		}
	}

	// search_state 用来表示识别的状态，-1 表示没有在识别中，0 表示刚识别到 [ , 后面需要判断有没有双引号，1 表示在识别具体内容中了
	search_state := -1

	var content strings.Builder
	var key_content strings.Builder
	var value_content strings.Builder

	var key string

	begin_index := -1 //不保存string的每个item，而是保存开头和结尾

	is_double_quoted_string := false // 用来在 part 4 或者 part 5 中表示是否是双引号中的字符串
	is_slash := false                // 用于 part 4/5 中区分 \

	is_search_for_key := true

	for char_index, item := range log {
		// 如果目前在识别的间隙
		if search_state == -1 {
			// 进入识别状态
			if item == '[' {
				if part < 5 {
					part += 1
				}
				begin_index = char_index + 1
				search_state = 0
				is_search_for_key = true
			}
		} else {
			/*
				// 通过遍历的方式会慢
				// 如果是前三种 part,等到识别到 ']' 再结束
				if part == 1 || part == 2 || part == 3 {
					if item == ']' {
						switch {
						case part == 1:
							unified_log.LogHeaderSection.DataTime = log[begin_index:char_index]
						case part == 2:
							unified_log.LogHeaderSection.Level = log[begin_index:char_index]
						case part == 3:
							value := log[begin_index:char_index]
							if value == "<unknown>" {
								unified_log.LogHeaderSection.SourceFile = "unknown"
							} else {
								index := strings.Index(value, ":")
								unified_log.LogHeaderSection.SourceFile = value[:index]
								number, err := strconv.Atoi(value[index+1:])
								if err != nil {
									fmt.Printf("unmatched number is %v, content is %v", number, value)
									fmt.Println()
								}
								unified_log.LogHeaderSection.LineNumber = number
							}
						}
						search_state = -1
					}
				}
			*/
			if part == 4 {
				//因为前面已经处理过了，所以这边的message肯定是带双引号的，所以不用额外判断了
				if search_state == 0 {
					is_double_quoted_string = true
					begin_index = char_index + 1
					search_state = 1
				} else {
					//如果里面有 slash 出现过，就必须要拿出来重新写过，否则就不需要
					if item == '"' && !is_slash {
						// 如果是前面没有有效的\的"， 就说明到了结尾
						if content.Len() == 0 { // 如果content是空的，就说明中途没有 slash，所以也没有重写过，直接拿对应的字符串就可以
							unified_log.LogMessageSection = log[begin_index:char_index]
						} else {
							content.WriteString(log[begin_index:char_index])
							unified_log.LogMessageSection = content.String()
							//content.Reset()
						}
						search_state = -1

					} else if item == '\\' && !is_slash {
						// 如果前面的不是 slash， 那就说明这是个转义的，要把前面这段先写入content
						content.WriteString(log[begin_index:char_index])
						is_slash = true
						begin_index = char_index + 1
					} else {
						// 其他情况都只需要把is_slash 设为 false 就可以
						is_slash = false
					}
				}
			} else if part == 5 {
				if search_state == 0 {
					if item == '=' { // 这边特殊处理一下 =
						continue
					} else if item == '"' {
						is_double_quoted_string = true
						begin_index = char_index + 1
						is_slash = false
					} else {
						is_double_quoted_string = false
						begin_index = char_index
						is_slash = false
					}
					search_state = 1
				} else {
					if is_double_quoted_string {
						if item == '"' && !is_slash {
							// 如果遇到 前面没有转义符号的 "
							if is_search_for_key {
								if key_content.Len() == 0 {
									key = log[begin_index:char_index]
								} else {
									key_content.WriteString(log[begin_index:char_index])
									key = key_content.String()
									key_content.Reset()
								}
								is_search_for_key = false
								search_state = 0
							} else {
								key_value_pair := make(map[string]string)
								if value_content.Len() == 0 {
									key_value_pair[key] = log[begin_index:char_index]
								} else {
									value_content.WriteString(log[begin_index:char_index])
									key_value_pair[key] = value_content.String()
									value_content.Reset()
								}
								unified_log.LogFieldsSection = append(unified_log.LogFieldsSection, key_value_pair)
								search_state = -1

							}
						} else if item == '\\' && !is_slash {
							// 如果前面的不是 slash， 那就说明这是个转义的，要把前面这段先写入content
							if is_search_for_key {
								key_content.WriteString(log[begin_index:char_index])
							} else {
								value_content.WriteString(log[begin_index:char_index])
							}

							is_slash = true
							begin_index = char_index + 1
						} else {
							// 其他情况都只需要把is_slash 设为 false 就可以
							is_slash = false
						}
					} else {
						// 如果是没有特殊字符的字符串，那只需要识别 = 或者 ] 就可以了
						if item == '=' {
							is_search_for_key = false
							search_state = 0
							key = log[begin_index:char_index]
						} else if item == ']' {
							key_value_pair := make(map[string]string)
							key_value_pair[key] = log[begin_index:char_index]
							unified_log.LogFieldsSection = append(unified_log.LogFieldsSection, key_value_pair)
							search_state = -1
						}
					}
				}
			}
		}

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

	var i = 1
	for i < 1000000 {
		i += 1
		log_test := "[2018/12/15 14:20:11.015 +08:00] [INFO] [tikv-server.rs:13] [\"TiKV Started\"] [ddl_job_id=1]"
		ParseLog(log_test)
	}

}
