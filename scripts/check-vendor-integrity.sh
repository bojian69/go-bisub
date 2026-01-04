#!/bin/bash

# 检查第三方资源完整性
# 验证本地资源文件是否存在且有效

set -e

VENDOR_JS_DIR="web/static/js/vendor"
VENDOR_CSS_DIR="web/static/css/vendor"

echo "🔍 检查第三方资源完整性..."

# 定义预期的文件和最小大小（字节）
declare -A EXPECTED_FILES=(
    ["$VENDOR_CSS_DIR/bootstrap.min.css"]=150000
    ["$VENDOR_CSS_DIR/bootstrap-icons.css"]=70000
    ["$VENDOR_JS_DIR/bootstrap.bundle.min.js"]=70000
    ["$VENDOR_JS_DIR/xlsx.full.min.js"]=800000
)

# 检查函数
check_file() {
    local file="$1"
    local min_size="$2"
    
    if [ ! -f "$file" ]; then
        echo "  ❌ 文件不存在: $file"
        return 1
    fi
    
    local actual_size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    
    if [ "$actual_size" -lt "$min_size" ]; then
        echo "  ❌ 文件大小异常: $file (实际: ${actual_size}B, 预期: >${min_size}B)"
        return 1
    fi
    
    # 检查文件是否为空或损坏
    if [ ! -s "$file" ]; then
        echo "  ❌ 文件为空: $file"
        return 1
    fi
    
    local size_mb=$(echo "scale=2; $actual_size/1024/1024" | bc 2>/dev/null || echo "N/A")
    echo "  ✅ $file (${size_mb}MB)"
    return 0
}

# 执行检查
all_good=true
for file in "${!EXPECTED_FILES[@]}"; do
    if ! check_file "$file" "${EXPECTED_FILES[$file]}"; then
        all_good=false
    fi
done

echo ""

if [ "$all_good" = true ]; then
    echo "🎉 所有第三方资源完整性检查通过！"
    echo ""
    echo "📊 资源统计："
    echo "  - Bootstrap CSS: $(du -h "$VENDOR_CSS_DIR/bootstrap.min.css" | cut -f1)"
    echo "  - Bootstrap JS:  $(du -h "$VENDOR_JS_DIR/bootstrap.bundle.min.js" | cut -f1)"
    echo "  - SheetJS:       $(du -h "$VENDOR_JS_DIR/xlsx.full.min.js" | cut -f1)"
    exit 0
else
    echo "❌ 发现问题！建议运行以下命令修复："
    echo "  ./scripts/update-vendor-assets.sh"
    exit 1
fi