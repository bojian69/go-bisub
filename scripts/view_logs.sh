#!/bin/bash

# 日志查看工具脚本
# 用法: ./scripts/view_logs.sh [选项]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认日志目录
LOG_DIR="./logs"

# 显示帮助信息
show_help() {
    cat << EOF
日志查看工具

用法: $0 [选项]

选项:
    -h, --help              显示帮助信息
    -d, --date DATE         指定日期 (格式: YYMMDD, 默认: 今天)
    -t, --type TYPE         日志类型 (api|sql, 默认: api)
    -f, --follow            实时跟踪日志
    -e, --errors            只显示错误日志
    -s, --slow              只显示慢查询 (仅SQL日志)
    -r, --request-id ID     按request_id过滤
    -p, --path PATH         按API路径过滤
    --stats                 显示统计信息
    --tail N                显示最后N行 (默认: 100)

示例:
    # 查看今天的API日志
    $0

    # 查看今天的SQL日志
    $0 -t sql

    # 实时跟踪API日志
    $0 -f

    # 查看错误日志
    $0 -e

    # 查看慢查询
    $0 -t sql -s

    # 按request_id查询
    $0 -r 550e8400-e29b-41d4-a716-446655440000

    # 显示统计信息
    $0 --stats

EOF
}

# 检查jq是否安装
check_jq() {
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}错误: 需要安装jq工具${NC}"
        echo "安装方法:"
        echo "  macOS: brew install jq"
        echo "  Ubuntu/Debian: sudo apt-get install jq"
        echo "  CentOS/RHEL: sudo yum install jq"
        exit 1
    fi
}

# 获取日志文件路径
get_log_file() {
    local date=$1
    local type=$2
    
    if [ "$type" = "sql" ]; then
        echo "${LOG_DIR}/${date}_sql.log"
    else
        echo "${LOG_DIR}/${date}.log"
    fi
}

# 显示API统计
show_api_stats() {
    local log_file=$1
    
    if [ ! -f "$log_file" ]; then
        echo -e "${RED}日志文件不存在: $log_file${NC}"
        return 1
    fi
    
    echo -e "${BLUE}=== API统计信息 ===${NC}"
    echo ""
    
    # 总请求数
    local total=$(cat "$log_file" | wc -l)
    echo -e "${GREEN}总请求数:${NC} $total"
    
    # 错误请求数
    local errors=$(cat "$log_file" | jq -r 'select(.status_code >= 400)' | wc -l)
    echo -e "${RED}错误请求数:${NC} $errors"
    
    # 平均响应时间
    local avg_duration=$(cat "$log_file" | jq -s 'map(.duration_ms) | add / length')
    echo -e "${YELLOW}平均响应时间:${NC} ${avg_duration}ms"
    
    # 最慢请求
    echo ""
    echo -e "${BLUE}=== 最慢的5个请求 ===${NC}"
    cat "$log_file" | jq -s 'sort_by(.duration_ms) | reverse | .[0:5] | .[] | {path, duration_ms, status_code}'
    
    # API路径统计
    echo ""
    echo -e "${BLUE}=== API调用统计 ===${NC}"
    cat "$log_file" | jq -r '.path' | sort | uniq -c | sort -rn | head -10
    
    # 状态码分布
    echo ""
    echo -e "${BLUE}=== 状态码分布 ===${NC}"
    cat "$log_file" | jq -r '.status_code' | sort | uniq -c | sort -rn
}

# 显示SQL统计
show_sql_stats() {
    local log_file=$1
    
    if [ ! -f "$log_file" ]; then
        echo -e "${RED}日志文件不存在: $log_file${NC}"
        return 1
    fi
    
    echo -e "${BLUE}=== SQL统计信息 ===${NC}"
    echo ""
    
    # 总查询数
    local total=$(cat "$log_file" | wc -l)
    echo -e "${GREEN}总查询数:${NC} $total"
    
    # 慢查询数
    local slow=$(cat "$log_file" | jq -r 'select(.sql | contains("[SLOW QUERY]"))' | wc -l)
    echo -e "${RED}慢查询数:${NC} $slow"
    
    # 平均执行时间
    local avg_duration=$(cat "$log_file" | jq -s 'map(.duration_ms) | add / length')
    echo -e "${YELLOW}平均执行时间:${NC} ${avg_duration}ms"
    
    # 执行时间分布
    echo ""
    echo -e "${BLUE}=== 执行时间分布 ===${NC}"
    cat "$log_file" | jq -r '.duration_ms' | awk '{
        if ($1 < 10) fast++
        else if ($1 < 100) normal++
        else if ($1 < 1000) slow++
        else very_slow++
    }
    END {
        printf "Fast (<10ms):       %d\n", fast
        printf "Normal (10-100ms):  %d\n", normal
        printf "Slow (100-1000ms):  %d\n", slow
        printf "Very Slow (>1000ms): %d\n", very_slow
    }'
    
    # 最慢的查询
    echo ""
    echo -e "${BLUE}=== 最慢的5个查询 ===${NC}"
    cat "$log_file" | jq -s 'sort_by(.duration_ms) | reverse | .[0:5] | .[] | {sql: (.sql | .[0:100]), duration_ms}'
}

# 主函数
main() {
    check_jq
    
    # 默认参数
    local date=$(date +%y%m%d)
    local log_type="api"
    local follow=false
    local errors_only=false
    local slow_only=false
    local request_id=""
    local path_filter=""
    local show_stats=false
    local tail_lines=100
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -d|--date)
                date="$2"
                shift 2
                ;;
            -t|--type)
                log_type="$2"
                shift 2
                ;;
            -f|--follow)
                follow=true
                shift
                ;;
            -e|--errors)
                errors_only=true
                shift
                ;;
            -s|--slow)
                slow_only=true
                shift
                ;;
            -r|--request-id)
                request_id="$2"
                shift 2
                ;;
            -p|--path)
                path_filter="$2"
                shift 2
                ;;
            --stats)
                show_stats=true
                shift
                ;;
            --tail)
                tail_lines="$2"
                shift 2
                ;;
            *)
                echo -e "${RED}未知选项: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 获取日志文件
    local log_file=$(get_log_file "$date" "$log_type")
    
    # 显示统计信息
    if [ "$show_stats" = true ]; then
        if [ "$log_type" = "sql" ]; then
            show_sql_stats "$log_file"
        else
            show_api_stats "$log_file"
        fi
        exit 0
    fi
    
    # 检查文件是否存在
    if [ ! -f "$log_file" ]; then
        echo -e "${RED}日志文件不存在: $log_file${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}查看日志: $log_file${NC}"
    echo ""
    
    # 构建jq过滤器
    local jq_filter="."
    
    if [ -n "$request_id" ]; then
        jq_filter="select(.request_id == \"$request_id\")"
    fi
    
    if [ "$errors_only" = true ]; then
        jq_filter="$jq_filter | select(.status_code >= 400)"
    fi
    
    if [ "$slow_only" = true ] && [ "$log_type" = "sql" ]; then
        jq_filter="$jq_filter | select(.sql | contains(\"[SLOW QUERY]\"))"
    fi
    
    if [ -n "$path_filter" ]; then
        jq_filter="$jq_filter | select(.path | contains(\"$path_filter\"))"
    fi
    
    # 显示日志
    if [ "$follow" = true ]; then
        tail -f "$log_file" | jq "$jq_filter"
    else
        tail -n "$tail_lines" "$log_file" | jq "$jq_filter"
    fi
}

main "$@"
