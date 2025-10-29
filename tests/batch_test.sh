#!/bin/bash

# DevBox Pack 批量测试脚本
# 用于遍历所有测试用例并生成执行计划

set -e

# 配置变量
DEVBOX_PACK_BIN="/home/sealos/Projects/labring/devbox-pack/impl/go/bin/devbox-pack"
EXAMPLES_DIR="/home/sealos/Projects/labring/devbox-pack/railpack/examples"
OUTPUT_DIR="/home/sealos/Projects/labring/devbox-pack/test_results"
LOG_FILE="$OUTPUT_DIR/batch_test.log"
VALIDATION_REPORT="$OUTPUT_DIR/validation_report.json"
ANALYSIS_REPORT="$OUTPUT_DIR/final_analysis_report.md"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log() {
    echo -e "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log_info() {
    log "${BLUE}[INFO]${NC} $1"
}

log_success() {
    log "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    log "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    log "${RED}[ERROR]${NC} $1"
}

# 创建输出目录
create_output_dirs() {
    log_info "创建输出目录..."
    mkdir -p "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR/plans"
    mkdir -p "$OUTPUT_DIR/logs"
    mkdir -p "$OUTPUT_DIR/errors"
    
    # 清空日志文件
    > "$LOG_FILE"
    log_info "输出目录创建完成: $OUTPUT_DIR"
}

# 获取所有测试用例目录
get_test_cases() {
    find "$EXAMPLES_DIR" -maxdepth 1 -type d -name "*" | \
    grep -v "^$EXAMPLES_DIR$" | \
    sort
}

# 验证 JSON 计划文件
validate_json_plan() {
    local plan_file="$1"
    local test_case_name="$2"
    
    # 检查 JSON 格式
    if ! jq empty "$plan_file" 2>/dev/null; then
        return 1
    fi
    
    # 检查必需字段
    local required_fields=("provider" "base" "runtime" "commands")
    for field in "${required_fields[@]}"; do
        if ! jq -e ".$field" "$plan_file" >/dev/null 2>&1; then
            log_warning "⚠️  $test_case_name - 缺少必需字段: $field"
        fi
    done
    
    # 检查 provider 字段值
    local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null)
    if [[ "$provider" == "unknown" || "$provider" == "null" ]]; then
        log_warning "⚠️  $test_case_name - Provider 字段无效: $provider"
    fi
    
    return 0
}

# 分析错误类型
analyze_error_type() {
    local error_file="$1"
    local error_content=$(cat "$error_file" 2>/dev/null || echo "")
    
    if [[ "$error_content" == *"no supported language or framework detected"* ]]; then
        echo "DETECTION_FAILURE"
    elif [[ "$error_content" == *"FILE_READ_ERROR"* ]]; then
        echo "FILE_READ_ERROR"
    elif [[ "$error_content" == *"timeout"* ]]; then
        echo "TIMEOUT"
    elif [[ "$error_content" == *"Exit code: 1"* ]]; then
        echo "EXECUTION_ERROR"
    else
        echo "UNKNOWN_ERROR"
    fi
}

# 处理单个测试用例
process_test_case() {
    local test_case_path="$1"
    local test_case_name=$(basename "$test_case_path")
    local plan_file="$OUTPUT_DIR/plans/${test_case_name}.json"
    local log_file="$OUTPUT_DIR/logs/${test_case_name}.log"
    local error_file="$OUTPUT_DIR/errors/${test_case_name}.error"
    
    log_info "处理测试用例: $test_case_name"
    
    # 检查测试用例目录是否存在
    if [[ ! -d "$test_case_path" ]]; then
        log_error "测试用例目录不存在: $test_case_path"
        echo "Directory not found: $test_case_path" > "$error_file"
        return 1
    fi
    
    # 运行 devbox-pack 生成执行计划
    # 使用临时文件来分离 stdout 和 stderr
    local temp_output="$OUTPUT_DIR/temp_${test_case_name}.out"
    local temp_error="$OUTPUT_DIR/temp_${test_case_name}.err"
    
    # 增加超时时间到 60 秒，某些复杂项目可能需要更多时间
    if timeout 60s "$DEVBOX_PACK_BIN" "$test_case_path" --offline --format json > "$temp_output" 2> "$temp_error"; then
        # 提取最后一行作为 JSON（假设 JSON 在最后一行）
        tail -n 1 "$temp_output" > "$plan_file"
        
        # 验证生成的 JSON 计划
        if validate_json_plan "$plan_file" "$test_case_name"; then
            log_success "✓ $test_case_name - 执行计划生成成功"
            # 保存完整日志
            cat "$temp_output" > "$log_file"
            rm -f "$temp_output" "$temp_error"
            return 0
        else
            log_error "✗ $test_case_name - 生成的 JSON 格式无效"
            # 保存错误信息
            {
                echo "JSON_VALIDATION_ERROR"
                echo "=== STDOUT ==="
                cat "$temp_output"
                echo "=== STDERR ==="
                cat "$temp_error"
                echo "=== EXTRACTED JSON ==="
                cat "$plan_file"
            } > "$error_file"
            rm -f "$plan_file" "$temp_output" "$temp_error"
            return 1
        fi
    else
        local exit_code=$?
        local error_type=$(analyze_error_type "$temp_error")
        log_error "✗ $test_case_name - 执行失败 (退出码: $exit_code, 错误类型: $error_type)"
        
        # 将错误信息保存到错误文件
        {
            echo "Exit code: $exit_code"
            echo "Error type: $error_type"
            echo "=== STDERR ==="
            cat "$temp_error" 2>/dev/null || echo "No stderr output"
            echo "=== STDOUT ==="
            cat "$temp_output" 2>/dev/null || echo "No stdout output"
        } > "$error_file"
        
        # 清理临时文件
        rm -f "$plan_file" "$temp_output" "$temp_error"
        return 1
    fi
}

# 生成验证报告
generate_validation_report() {
    local validation_data="{\"timestamp\":\"$(date -Iseconds)\",\"results\":[]}"
    
    # 收集成功案例的验证信息
    for plan_file in "$OUTPUT_DIR/plans"/*.json; do
        if [[ -f "$plan_file" ]]; then
            local case_name=$(basename "$plan_file" .json)
            local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
            local has_commands=$(jq -e '.commands' "$plan_file" >/dev/null 2>&1 && echo "true" || echo "false")
            local has_base=$(jq -e '.base' "$plan_file" >/dev/null 2>&1 && echo "true" || echo "false")
            
            validation_data=$(echo "$validation_data" | jq --arg name "$case_name" --arg provider "$provider" --argjson commands "$has_commands" --argjson base "$has_base" \
                '.results += [{"name": $name, "status": "success", "provider": $provider, "has_commands": $commands, "has_base": $base}]')
        fi
    done
    
    # 收集失败案例的验证信息
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local case_name=$(basename "$error_file" .error)
            local error_type=$(analyze_error_type "$error_file")
            
            validation_data=$(echo "$validation_data" | jq --arg name "$case_name" --arg error_type "$error_type" \
                '.results += [{"name": $name, "status": "failed", "error_type": $error_type}]')
        fi
    done
    
    echo "$validation_data" | jq '.' > "$VALIDATION_REPORT"
    log_info "验证报告生成完成: $VALIDATION_REPORT"
}

# 生成详细分析报告
generate_analysis_report() {
    local total_cases=0
    local success_cases=0
    local failed_cases=0
    
    # 统计结果
    success_cases=$(find "$OUTPUT_DIR/plans" -name "*.json" | wc -l)
    failed_cases=$(find "$OUTPUT_DIR/errors" -name "*.error" | wc -l)
    total_cases=$((success_cases + failed_cases))
    
    # 错误类型统计
    declare -A error_types
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local error_type=$(analyze_error_type "$error_file")
            error_types["$error_type"]=$((${error_types["$error_type"]} + 1))
        fi
    done
    
    # Provider 统计
    declare -A provider_count
    declare -A provider_success
    for plan_file in "$OUTPUT_DIR/plans"/*.json; do
        if [[ -f "$plan_file" ]]; then
            local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
            provider_count["$provider"]=$((${provider_count["$provider"]} + 1))
            provider_success["$provider"]=$((${provider_success["$provider"]} + 1))
        fi
    done
    
    # 为失败的案例也统计 provider（如果能从目录名推断）
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local case_name=$(basename "$error_file" .error)
            local inferred_provider="unknown"
            
            # 根据测试用例名称推断 provider
            if [[ "$case_name" == node-* ]]; then
                inferred_provider="node"
            elif [[ "$case_name" == python-* ]]; then
                inferred_provider="python"
            elif [[ "$case_name" == java-* ]]; then
                inferred_provider="java"
            elif [[ "$case_name" == go-* ]]; then
                inferred_provider="go"
            elif [[ "$case_name" == rust-* ]]; then
                inferred_provider="rust"
            elif [[ "$case_name" == ruby-* ]]; then
                inferred_provider="ruby"
            elif [[ "$case_name" == php-* ]]; then
                inferred_provider="php"
            elif [[ "$case_name" == deno-* ]]; then
                inferred_provider="deno"
            elif [[ "$case_name" == elixir-* ]]; then
                inferred_provider="elixir"
            elif [[ "$case_name" == staticfile-* ]]; then
                inferred_provider="staticfile"
            elif [[ "$case_name" == shell-* ]]; then
                inferred_provider="shell"
            fi
            
            provider_count["$inferred_provider"]=$((${provider_count["$inferred_provider"]} + 1))
        fi
    done
    
    # 生成详细分析报告
    cat > "$ANALYSIS_REPORT" << EOF
# DevBox Pack 详细分析报告

## 执行概览

- **测试时间**: $(date '+%Y-%m-%d %H:%M:%S')
- **DevBox Pack 版本**: $($DEVBOX_PACK_BIN --version 2>/dev/null || echo "未知")
- **总测试用例**: $total_cases
- **成功**: $success_cases
- **失败**: $failed_cases
- **成功率**: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%

## Provider 详细统计

| Provider | 总数 | 成功 | 失败 | 成功率 |
|----------|------|------|------|--------|
EOF

    for provider in $(printf '%s\n' "${!provider_count[@]}" | sort); do
        local total=${provider_count[$provider]}
        local success=${provider_success[$provider]:-0}
        local failed=$((total - success))
        local success_rate=$(( total > 0 ? success * 100 / total : 0 ))
        echo "| $provider | $total | $success | $failed | ${success_rate}% |" >> "$ANALYSIS_REPORT"
    done
    
    cat >> "$ANALYSIS_REPORT" << EOF

## 错误类型分析

EOF

    if [[ ${#error_types[@]} -gt 0 ]]; then
        echo "| 错误类型 | 数量 | 占比 |" >> "$ANALYSIS_REPORT"
        echo "|----------|------|------|" >> "$ANALYSIS_REPORT"
        
        for error_type in "${!error_types[@]}"; do
            local count=${error_types[$error_type]}
            local percentage=$(( failed_cases > 0 ? count * 100 / failed_cases : 0 ))
            echo "| $error_type | $count | ${percentage}% |" >> "$ANALYSIS_REPORT"
        done
    else
        echo "无错误发生" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## 成功的测试用例

EOF

    if [[ $success_cases -gt 0 ]]; then
        for plan_file in "$OUTPUT_DIR/plans"/*.json; do
            if [[ -f "$plan_file" ]]; then
                local case_name=$(basename "$plan_file" .json)
                local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
                local base_image=$(jq -r '.base.name // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
                local runtime_version=$(jq -r '.runtime.version // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
                echo "- **$case_name** (Provider: $provider, Base: $base_image, Runtime: $runtime_version)" >> "$ANALYSIS_REPORT"
            fi
        done
    else
        echo "无成功的测试用例" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## 失败的测试用例详情

EOF

    if [[ $failed_cases -gt 0 ]]; then
        for error_file in "$OUTPUT_DIR/errors"/*.error; do
            if [[ -f "$error_file" ]]; then
                local case_name=$(basename "$error_file" .error)
                local error_type=$(analyze_error_type "$error_file")
                local error_preview=$(head -n 5 "$error_file" | tail -n +3 | tr '\n' ' ' | cut -c1-100)
                echo "- **$case_name** (错误类型: $error_type)" >> "$ANALYSIS_REPORT"
                echo "  - 错误详情: $error_preview..." >> "$ANALYSIS_REPORT"
                echo "" >> "$ANALYSIS_REPORT"
            fi
        done
    else
        echo "无失败的测试用例" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## 建议和改进方向

### 高优先级问题
EOF

    # 分析主要问题并给出建议
    local detection_failures=${error_types["DETECTION_FAILURE"]:-0}
    local file_read_errors=${error_types["FILE_READ_ERROR"]:-0}
    
    if [[ $detection_failures -gt 0 ]]; then
        cat >> "$ANALYSIS_REPORT" << EOF
- **语言/框架检测失败** ($detection_failures 个案例): 需要改进检测逻辑，可能是配置文件缺失或检测规则不完整
EOF
    fi
    
    if [[ $file_read_errors -gt 0 ]]; then
        cat >> "$ANALYSIS_REPORT" << EOF
- **文件读取错误** ($file_read_errors 个案例): 检查文件权限和路径配置
EOF
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

### Provider 特定建议
EOF

    for provider in "${!provider_count[@]}"; do
        local total=${provider_count[$provider]}
        local success=${provider_success[$provider]:-0}
        local success_rate=$(( total > 0 ? success * 100 / total : 0 ))
        
        if [[ $success_rate -lt 80 && $total -gt 1 ]]; then
            echo "- **$provider**: 成功率较低 (${success_rate}%)，需要重点关注和改进" >> "$ANALYSIS_REPORT"
        fi
    done
    
    log_success "详细分析报告生成完成: $ANALYSIS_REPORT"
}

# 生成统计报告
generate_summary() {
    local total_cases=0
    local success_cases=0
    local failed_cases=0
    local summary_file="$OUTPUT_DIR/summary.md"
    
    log_info "生成统计报告..."
    
    # 统计结果
    success_cases=$(find "$OUTPUT_DIR/plans" -name "*.json" | wc -l)
    failed_cases=$(find "$OUTPUT_DIR/errors" -name "*.error" | wc -l)
    total_cases=$((success_cases + failed_cases))
    
    # 生成 Markdown 报告
    cat > "$summary_file" << EOF
# DevBox Pack 批量测试报告

## 测试概览

- **测试时间**: $(date '+%Y-%m-%d %H:%M:%S')
- **总测试用例**: $total_cases
- **成功**: $success_cases
- **失败**: $failed_cases
- **成功率**: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%

## 成功的测试用例

EOF

    # 列出成功的测试用例
    if [[ $success_cases -gt 0 ]]; then
        for plan_file in "$OUTPUT_DIR/plans"/*.json; do
            if [[ -f "$plan_file" ]]; then
                local case_name=$(basename "$plan_file" .json)
                local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
                echo "- **$case_name** (Provider: $provider)" >> "$summary_file"
            fi
        done
    else
        echo "无成功的测试用例" >> "$summary_file"
    fi
    
    echo "" >> "$summary_file"
    echo "## 失败的测试用例" >> "$summary_file"
    echo "" >> "$summary_file"
    
    # 列出失败的测试用例
    if [[ $failed_cases -gt 0 ]]; then
        for error_file in "$OUTPUT_DIR/errors"/*.error; do
            if [[ -f "$error_file" ]]; then
                local case_name=$(basename "$error_file" .error)
                local error_preview=$(head -n 3 "$error_file" | tr '\n' ' ')
                echo "- **$case_name**: $error_preview" >> "$summary_file"
            fi
        done
    else
        echo "无失败的测试用例" >> "$summary_file"
    fi
    
    # 按 Provider 分组统计
    echo "" >> "$summary_file"
    echo "## Provider 统计" >> "$summary_file"
    echo "" >> "$summary_file"
    
    declare -A provider_count
    for plan_file in "$OUTPUT_DIR/plans"/*.json; do
        if [[ -f "$plan_file" ]]; then
            local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
            provider_count["$provider"]=$((${provider_count["$provider"]} + 1))
        fi
    done
    
    for provider in "${!provider_count[@]}"; do
        echo "- **$provider**: ${provider_count[$provider]} 个测试用例" >> "$summary_file"
    done
    
    # 生成验证报告和详细分析报告
    generate_validation_report
    generate_analysis_report
    
    log_success "统计报告生成完成: $summary_file"
    
    # 在控制台输出摘要
    echo ""
    log_info "=== 测试结果摘要 ==="
    log_info "总测试用例: $total_cases"
    log_success "成功: $success_cases"
    log_error "失败: $failed_cases"
    log_info "成功率: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%"
    echo ""
    log_info "详细报告文件:"
    log_info "- 基础报告: $summary_file"
    log_info "- 验证报告: $VALIDATION_REPORT"
    log_info "- 分析报告: $ANALYSIS_REPORT"
}

# 主函数
main() {
    log_info "开始 DevBox Pack 批量测试..."
    log_info "DevBox Pack 二进制文件: $DEVBOX_PACK_BIN"
    log_info "测试用例目录: $EXAMPLES_DIR"
    log_info "输出目录: $OUTPUT_DIR"
    
    # 检查 devbox-pack 是否存在
    if [[ ! -x "$DEVBOX_PACK_BIN" ]]; then
        log_error "DevBox Pack 二进制文件不存在或不可执行: $DEVBOX_PACK_BIN"
        exit 1
    fi
    
    # 检查测试用例目录是否存在
    if [[ ! -d "$EXAMPLES_DIR" ]]; then
        log_error "测试用例目录不存在: $EXAMPLES_DIR"
        exit 1
    fi
    
    # 创建输出目录
    create_output_dirs
    
    # 获取所有测试用例
    local test_cases
    test_cases=$(get_test_cases)
    local total_count=$(echo "$test_cases" | wc -l)
    
    log_info "发现 $total_count 个测试用例"
    
    # 处理每个测试用例
    local current=0
    while IFS= read -r test_case_path; do
        current=$((current + 1))
        local test_case_name=$(basename "$test_case_path")
        
        log_info "[$current/$total_count] 处理: $test_case_name"
        
        if ! process_test_case "$test_case_path"; then
            log_warning "测试用例处理失败: $test_case_name"
        fi
        
        # 显示进度
        local progress=$((current * 100 / total_count))
        log_info "进度: $progress% ($current/$total_count)"
        
    done <<< "$test_cases"
    
    # 生成统计报告
    generate_summary
    
    log_success "批量测试完成！"
    log_info "结果保存在: $OUTPUT_DIR"
}

# 运行主函数
main "$@"