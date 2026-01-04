#!/bin/bash

# 更新第三方资源脚本
# 用于下载和更新本地化的第三方库

set -e

VENDOR_JS_DIR="web/static/js/vendor"
VENDOR_CSS_DIR="web/static/css/vendor"

# 创建目录
mkdir -p "$VENDOR_JS_DIR" "$VENDOR_CSS_DIR"

echo "🔄 开始更新第三方资源..."

# Bootstrap 5.1.3
echo "📦 下载 Bootstrap..."
curl -L -o "$VENDOR_CSS_DIR/bootstrap.min.css" \
    "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css"
curl -L -o "$VENDOR_JS_DIR/bootstrap.bundle.min.js" \
    "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"

# SheetJS (xlsx) 0.18.5
echo "📊 下载 SheetJS..."
curl -L -o "$VENDOR_JS_DIR/xlsx.full.min.js" \
    "https://cdnjs.cloudflare.com/ajax/libs/xlsx/0.18.5/xlsx.full.min.js"

# 验证文件
echo "✅ 验证下载的文件..."
for file in "$VENDOR_CSS_DIR/bootstrap.min.css" \
           "$VENDOR_JS_DIR/bootstrap.bundle.min.js" \
           "$VENDOR_JS_DIR/xlsx.full.min.js"; do
    if [ -f "$file" ] && [ -s "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "  ✓ $file ($size)"
    else
        echo "  ❌ $file 下载失败或为空"
        exit 1
    fi
done

echo "🎉 所有第三方资源更新完成！"
echo ""
echo "📋 资源清单："
echo "  - Bootstrap 5.1.3 (CSS + JS)"
echo "  - SheetJS 0.18.5 (Excel导出)"
echo ""
echo "💡 提示：这些资源现在存储在本地，不再依赖外部CDN"