#!/bin/bash

# DevBox Pack batch testing script
# Used to traverse all test cases and generate execution plans

set -e

# Configuration variables
DEVBOX_PACK_BIN="/home/sealos/Projects/labring/devbox-pack/impl/go/bin/devbox-pack"
EXAMPLES_DIR="/home/sealos/Projects/labring/devbox-pack/railpack/examples"
OUTPUT_DIR="/home/sealos/Projects/labring/devbox-pack/test_results"
LOG_FILE="$OUTPUT_DIR/batch_test.log"
VALIDATION_REPORT="$OUTPUT_DIR/validation_report.json"
ANALYSIS_REPORT="$OUTPUT_DIR/final_analysis_report.md"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

# Create output directories
create_output_dirs() {
    log_info "Creating output directories..."
    mkdir -p "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR/plans"
    mkdir -p "$OUTPUT_DIR/logs"
    mkdir -p "$OUTPUT_DIR/errors"
    
    # Clear log file
    > "$LOG_FILE"
    log_info "Output directories created: $OUTPUT_DIR"
}

# Get all test case directories
get_test_cases() {
    find "$EXAMPLES_DIR" -maxdepth 1 -type d -name "*" | \
    grep -v "^$EXAMPLES_DIR$" | \
    sort
}

# Validate JSON plan file
validate_json_plan() {
    local plan_file="$1"
    local test_case_name="$2"
    
    # Check JSON format
    if ! jq empty "$plan_file" 2>/dev/null; then
        return 1
    fi
    
    # Check required fields
    local required_fields=("provider" "base" "runtime" "commands")
    for field in "${required_fields[@]}"; do
        if ! jq -e ".$field" "$plan_file" >/dev/null 2>&1; then
            log_warning "⚠️  $test_case_name - Missing required field: $field"
        fi
    done
    
    # Check provider field value
    local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null)
    if [[ "$provider" == "unknown" || "$provider" == "null" ]]; then
        log_warning "⚠️  $test_case_name - Invalid provider field: $provider"
    fi
    
    return 0
}

# Analyze error type
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

# Process single test case
process_test_case() {
    local test_case_path="$1"
    local test_case_name=$(basename "$test_case_path")
    local plan_file="$OUTPUT_DIR/plans/${test_case_name}.json"
    local log_file="$OUTPUT_DIR/logs/${test_case_name}.log"
    local error_file="$OUTPUT_DIR/errors/${test_case_name}.error"
    
    log_info "Processing test case: $test_case_name"
    
    # Check if test case directory exists
    if [[ ! -d "$test_case_path" ]]; then
        log_error "Test case directory does not exist: $test_case_path"
        echo "Directory not found: $test_case_path" > "$error_file"
        return 1
    fi
    
    # Run devbox-pack to generate execution plan
    # Use temporary files to separate stdout and stderr
    local temp_output="$OUTPUT_DIR/temp_${test_case_name}.out"
    local temp_error="$OUTPUT_DIR/temp_${test_case_name}.err"
    
    # Increase timeout to 60 seconds, some complex projects may need more time
    if timeout 60s "$DEVBOX_PACK_BIN" "$test_case_path" --offline --format json > "$temp_output" 2> "$temp_error"; then
        # Extract last line as JSON (assuming JSON is on the last line)
        tail -n 1 "$temp_output" > "$plan_file"
        
        # Validate generated JSON plan
        if validate_json_plan "$plan_file" "$test_case_name"; then
            log_success "✓ $test_case_name - Execution plan generated successfully"
            # Save complete log
            cat "$temp_output" > "$log_file"
            rm -f "$temp_output" "$temp_error"
            return 0
        else
            log_error "✗ $test_case_name - Generated JSON format is invalid"
            # Save error information
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
        log_error "✗ $test_case_name - Execution failed (exit code: $exit_code, error type: $error_type)"
        
        # Save error information to error file
        {
            echo "Exit code: $exit_code"
            echo "Error type: $error_type"
            echo "=== STDERR ==="
            cat "$temp_error" 2>/dev/null || echo "No stderr output"
            echo "=== STDOUT ==="
            cat "$temp_output" 2>/dev/null || echo "No stdout output"
        } > "$error_file"
        
        # Clean up temporary files
        rm -f "$plan_file" "$temp_output" "$temp_error"
        return 1
    fi
}

# Generate validation report
generate_validation_report() {
    local validation_data="{\"timestamp\":\"$(date -Iseconds)\",\"results\":[]}"
    
    # Collect validation information for successful cases
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
    
    # Collect validation information for failed cases
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local case_name=$(basename "$error_file" .error)
            local error_type=$(analyze_error_type "$error_file")
            
            validation_data=$(echo "$validation_data" | jq --arg name "$case_name" --arg error_type "$error_type" \
                '.results += [{"name": $name, "status": "failed", "error_type": $error_type}]')
        fi
    done
    
    echo "$validation_data" | jq '.' > "$VALIDATION_REPORT"
    log_info "Validation report generated: $VALIDATION_REPORT"
}

# Generate detailed analysis report
generate_analysis_report() {
    local total_cases=0
    local success_cases=0
    local failed_cases=0
    
    # Statistical results
    success_cases=$(find "$OUTPUT_DIR/plans" -name "*.json" | wc -l)
    failed_cases=$(find "$OUTPUT_DIR/errors" -name "*.error" | wc -l)
    total_cases=$((success_cases + failed_cases))
    
    # Error type statistics
    declare -A error_types
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local error_type=$(analyze_error_type "$error_file")
            error_types["$error_type"]=$((${error_types["$error_type"]} + 1))
        fi
    done
    
    # Provider statistics
    declare -A provider_count
    declare -A provider_success
    for plan_file in "$OUTPUT_DIR/plans"/*.json; do
        if [[ -f "$plan_file" ]]; then
            local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
            provider_count["$provider"]=$((${provider_count["$provider"]} + 1))
            provider_success["$provider"]=$((${provider_success["$provider"]} + 1))
        fi
    done
    
    # Also count provider for failed cases (if can be inferred from directory name)
    for error_file in "$OUTPUT_DIR/errors"/*.error; do
        if [[ -f "$error_file" ]]; then
            local case_name=$(basename "$error_file" .error)
            local inferred_provider="unknown"
            
            # Infer provider based on test case name
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
    
    # Generate detailed analysis report
    cat > "$ANALYSIS_REPORT" << EOF
# DevBox Pack Detailed Analysis Report

## Execution Overview

- **Test Time**: $(date '+%Y-%m-%d %H:%M:%S')
- **DevBox Pack Version**: $($DEVBOX_PACK_BIN --version 2>/dev/null || echo "Unknown")
- **Total Test Cases**: $total_cases
- **Successful**: $success_cases
- **Failed**: $failed_cases
- **Success Rate**: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%

## Provider Detailed Statistics

| Provider | Total | Success | Failed | Success Rate |
|----------|-------|---------|--------|--------------|
EOF

    for provider in $(printf '%s\n' "${!provider_count[@]}" | sort); do
        local total=${provider_count[$provider]}
        local success=${provider_success[$provider]:-0}
        local failed=$((total - success))
        local success_rate=$(( total > 0 ? success * 100 / total : 0 ))
        echo "| $provider | $total | $success | $failed | ${success_rate}% |" >> "$ANALYSIS_REPORT"
    done
    
    cat >> "$ANALYSIS_REPORT" << EOF

## Error Type Analysis

EOF

    if [[ ${#error_types[@]} -gt 0 ]]; then
        echo "| Error Type | Count | Percentage |" >> "$ANALYSIS_REPORT"
        echo "|------------|-------|------------|" >> "$ANALYSIS_REPORT"
        
        for error_type in "${!error_types[@]}"; do
            local count=${error_types[$error_type]}
            local percentage=$(( failed_cases > 0 ? count * 100 / failed_cases : 0 ))
            echo "| $error_type | $count | ${percentage}% |" >> "$ANALYSIS_REPORT"
        done
    else
        echo "No errors occurred" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## Successful Test Cases

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
        echo "No successful test cases" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## Failed Test Cases Details

EOF

    if [[ $failed_cases -gt 0 ]]; then
        for error_file in "$OUTPUT_DIR/errors"/*.error; do
            if [[ -f "$error_file" ]]; then
                local case_name=$(basename "$error_file" .error)
                local error_type=$(analyze_error_type "$error_file")
                local error_preview=$(head -n 5 "$error_file" | tail -n +3 | tr '\n' ' ' | cut -c1-100)
                echo "- **$case_name** (Error Type: $error_type)" >> "$ANALYSIS_REPORT"
                echo "  - Error Details: $error_preview..." >> "$ANALYSIS_REPORT"
                echo "" >> "$ANALYSIS_REPORT"
            fi
        done
    else
        echo "No failed test cases" >> "$ANALYSIS_REPORT"
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

## Recommendations and Improvement Directions

### High Priority Issues
EOF

    # Analyze main issues and provide recommendations
    local detection_failures=${error_types["DETECTION_FAILURE"]:-0}
    local file_read_errors=${error_types["FILE_READ_ERROR"]:-0}
    
    if [[ $detection_failures -gt 0 ]]; then
        cat >> "$ANALYSIS_REPORT" << EOF
- **Language/Framework Detection Failure** ($detection_failures cases): Need to improve detection logic, possibly due to missing configuration files or incomplete detection rules
EOF
    fi
    
    if [[ $file_read_errors -gt 0 ]]; then
        cat >> "$ANALYSIS_REPORT" << EOF
- **File Read Errors** ($file_read_errors cases): Check file permissions and path configuration
EOF
    fi
    
    cat >> "$ANALYSIS_REPORT" << EOF

### Provider-Specific Recommendations
EOF

    for provider in "${!provider_count[@]}"; do
        local total=${provider_count[$provider]}
        local success=${provider_success[$provider]:-0}
        local success_rate=$(( total > 0 ? success * 100 / total : 0 ))
        
        if [[ $success_rate -lt 80 && $total -gt 1 ]]; then
            echo "- **$provider**: Low success rate (${success_rate}%), needs focused attention and improvement" >> "$ANALYSIS_REPORT"
        fi
    done
    
    log_success "Detailed analysis report generated: $ANALYSIS_REPORT"
}

# Generate summary report
generate_summary() {
    local total_cases=0
    local success_cases=0
    local failed_cases=0
    local summary_file="$OUTPUT_DIR/summary.md"
    
    log_info "Generating summary report..."
    
    # Statistical results
    success_cases=$(find "$OUTPUT_DIR/plans" -name "*.json" | wc -l)
    failed_cases=$(find "$OUTPUT_DIR/errors" -name "*.error" | wc -l)
    total_cases=$((success_cases + failed_cases))
    
    # Generate Markdown report
    cat > "$summary_file" << EOF
# DevBox Pack Batch Test Report

## Test Overview

- **Test Time**: $(date '+%Y-%m-%d %H:%M:%S')
- **Total Test Cases**: $total_cases
- **Successful**: $success_cases
- **Failed**: $failed_cases
- **Success Rate**: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%

## Successful Test Cases

EOF

    # List successful test cases
    if [[ $success_cases -gt 0 ]]; then
        for plan_file in "$OUTPUT_DIR/plans"/*.json; do
            if [[ -f "$plan_file" ]]; then
                local case_name=$(basename "$plan_file" .json)
                local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
                echo "- **$case_name** (Provider: $provider)" >> "$summary_file"
            fi
        done
    else
        echo "No successful test cases" >> "$summary_file"
    fi
    
    echo "" >> "$summary_file"
    echo "## Failed Test Cases" >> "$summary_file"
    echo "" >> "$summary_file"
    
    # List failed test cases
    if [[ $failed_cases -gt 0 ]]; then
        for error_file in "$OUTPUT_DIR/errors"/*.error; do
            if [[ -f "$error_file" ]]; then
                local case_name=$(basename "$error_file" .error)
                local error_preview=$(head -n 3 "$error_file" | tr '\n' ' ')
                echo "- **$case_name**: $error_preview" >> "$summary_file"
            fi
        done
    else
        echo "No failed test cases" >> "$summary_file"
    fi
    
    # Group statistics by Provider
    echo "" >> "$summary_file"
    echo "## Provider Statistics" >> "$summary_file"
    echo "" >> "$summary_file"
    
    declare -A provider_count
    for plan_file in "$OUTPUT_DIR/plans"/*.json; do
        if [[ -f "$plan_file" ]]; then
            local provider=$(jq -r '.provider // "unknown"' "$plan_file" 2>/dev/null || echo "unknown")
            provider_count["$provider"]=$((${provider_count["$provider"]} + 1))
        fi
    done
    
    for provider in "${!provider_count[@]}"; do
        echo "- **$provider**: ${provider_count[$provider]} test cases" >> "$summary_file"
    done
    
    # Generate validation report and detailed analysis report
    generate_validation_report
    generate_analysis_report
    
    log_success "Summary report generated: $summary_file"
    
    # Output summary to console
    echo ""
    log_info "=== Test Results Summary ==="
    log_info "Total test cases: $total_cases"
    log_success "Successful: $success_cases"
    log_error "Failed: $failed_cases"
    log_info "Success rate: $(( total_cases > 0 ? success_cases * 100 / total_cases : 0 ))%"
    echo ""
    log_info "Detailed report files:"
    log_info "- Basic report: $summary_file"
    log_info "- Validation report: $VALIDATION_REPORT"
    log_info "- Analysis report: $ANALYSIS_REPORT"
}

# Main function
main() {
    log_info "Starting DevBox Pack batch testing..."
    log_info "DevBox Pack binary: $DEVBOX_PACK_BIN"
    log_info "Test cases directory: $EXAMPLES_DIR"
    log_info "Output directory: $OUTPUT_DIR"
    
    # Check if devbox-pack exists
    if [[ ! -x "$DEVBOX_PACK_BIN" ]]; then
        log_error "DevBox Pack binary does not exist or is not executable: $DEVBOX_PACK_BIN"
        exit 1
    fi
    
    # Check if test cases directory exists
    if [[ ! -d "$EXAMPLES_DIR" ]]; then
        log_error "Test cases directory does not exist: $EXAMPLES_DIR"
        exit 1
    fi
    
    # Create output directories
    create_output_dirs
    
    # Get all test cases
    local test_cases
    test_cases=$(get_test_cases)
    local total_count=$(echo "$test_cases" | wc -l)
    
    log_info "Found $total_count test cases"
    
    # Process each test case
    local current=0
    while IFS= read -r test_case_path; do
        current=$((current + 1))
        local test_case_name=$(basename "$test_case_path")
        
        log_info "[$current/$total_count] Processing: $test_case_name"
        
        if ! process_test_case "$test_case_path"; then
            log_warning "Test case processing failed: $test_case_name"
        fi
        
        # Show progress
        local progress=$((current * 100 / total_count))
        log_info "Progress: $progress% ($current/$total_count)"
        
    done <<< "$test_cases"
    
    # Generate summary report
    generate_summary
    
    log_success "Batch testing completed!"
    log_info "Results saved in: $OUTPUT_DIR"
}

# Run main function
main "$@"