// BI 订阅管理系统 - 通用 JavaScript 工具

// API 配置
const API_CONFIG = {
    baseURL: '/api',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
};

// Toast 通知（使用 Bootstrap Toast）
class Toast {
    static show(message, type = 'info') {
        const toastContainer = document.getElementById('toastContainer');
        if (!toastContainer) {
            this.createToastContainer();
        }

        const toastId = 'toast-' + Date.now();
        const bgClass = this.getBackgroundClass(type);

        const toastHTML = `
            <div id="${toastId}" class="toast align-items-center text-white ${bgClass} border-0" role="alert">
                <div class="d-flex">
                    <div class="toast-body">
                        ${this.getIcon(type)} ${message}
                    </div>
                    <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
                </div>
            </div>
        `;

        document.getElementById('toastContainer').insertAdjacentHTML('beforeend', toastHTML);
        const toastElement = document.getElementById(toastId);
        const toast = new bootstrap.Toast(toastElement, { delay: 3000 });
        toast.show();

        toastElement.addEventListener('hidden.bs.toast', () => {
            toastElement.remove();
        });
    }

    static createToastContainer() {
        const container = document.createElement('div');
        container.id = 'toastContainer';
        container.className = 'toast-container position-fixed top-0 end-0 p-3';
        container.style.zIndex = '9999';
        document.body.appendChild(container);
    }

    static getBackgroundClass(type) {
        const classes = {
            success: 'bg-success',
            error: 'bg-danger',
            warning: 'bg-warning',
            info: 'bg-info'
        };
        return classes[type] || 'bg-info';
    }

    static getIcon(type) {
        const icons = {
            success: '✅',
            error: '❌',
            warning: '⚠️',
            info: 'ℹ️'
        };
        return icons[type] || 'ℹ️';
    }

    static success(message) {
        this.show(message, 'success');
    }

    static error(message) {
        this.show(message, 'error');
    }

    static warning(message) {
        this.show(message, 'warning');
    }

    static info(message) {
        this.show(message, 'info');
    }
}

// API 请求封装
class API {
    static async request(url, options = {}) {
        const config = {
            ...API_CONFIG,
            ...options,
            headers: {
                ...API_CONFIG.headers,
                ...options.headers
            }
        };

        try {
            const response = await fetch(API_CONFIG.baseURL + url, config);
            const result = await response.json();

            if (result.code === 'OK') {
                return { success: true, data: result.data, message: result.message };
            } else {
                return { success: false, error: result.message, code: result.code };
            }
        } catch (error) {
            console.error('API request error:', error);
            return { success: false, error: error.message };
        }
    }

    static async get(url, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const fullUrl = queryString ? `${url}?${queryString}` : url;
        return this.request(fullUrl, { method: 'GET' });
    }

    static async post(url, data = {}) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    static async put(url, data = {}) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    static async delete(url) {
        return this.request(url, { method: 'DELETE' });
    }
}

// 日期时间格式化
class DateFormatter {
    static format(date, format = 'YYYY-MM-DD HH:mm:ss') {
        if (!date) return '-';

        const d = new Date(date);
        if (isNaN(d.getTime())) return '-';

        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');

        return format
            .replace('YYYY', year)
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    }

    static formatRelative(date) {
        if (!date) return '-';

        const d = new Date(date);
        const now = new Date();
        const diff = now - d;

        const seconds = Math.floor(diff / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);

        if (days > 7) return this.format(date, 'YYYY-MM-DD');
        if (days > 0) return `${days}天前`;
        if (hours > 0) return `${hours}小时前`;
        if (minutes > 0) return `${minutes}分钟前`;
        return '刚刚';
    }
}

// 数据验证
class Validator {
    static isJSON(str) {
        try {
            JSON.parse(str);
            return true;
        } catch (e) {
            return false;
        }
    }

    static isEmpty(value) {
        return value === null || value === undefined || value === '';
    }

    static isValidKey(key) {
        return /^[a-zA-Z0-9_-]+$/.test(key);
    }

    static isValidVersion(version) {
        return Number.isInteger(version) && version > 0 && version <= 255;
    }
}

// 确认对话框
class Confirm {
    static async show(message, title = '确认操作') {
        return new Promise((resolve) => {
            const result = confirm(`${title}\n\n${message}`);
            resolve(result);
        });
    }

    static async delete(itemName) {
        return this.show(`确定要删除 ${itemName} 吗？此操作不可恢复。`, '确认删除');
    }
}

// 加载状态管理
class Loading {
    static show(element) {
        if (typeof element === 'string') {
            element = document.querySelector(element);
        }
        if (element) {
            element.disabled = true;
            element.dataset.originalText = element.innerHTML;
            element.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>加载中...';
        }
    }

    static hide(element) {
        if (typeof element === 'string') {
            element = document.querySelector(element);
        }
        if (element && element.dataset.originalText) {
            element.disabled = false;
            element.innerHTML = element.dataset.originalText;
            delete element.dataset.originalText;
        }
    }
}

// 表格工具
class TableHelper {
    static renderEmptyRow(colspan, message = '暂无数据') {
        return `<tr><td colspan="${colspan}" class="text-center text-muted py-4">${message}</td></tr>`;
    }

    static renderLoadingRow(colspan) {
        return `
            <tr>
                <td colspan="${colspan}" class="text-center py-4">
                    <div class="spinner-border text-primary" role="status">
                        <span class="visually-hidden">加载中...</span>
                    </div>
                </td>
            </tr>
        `;
    }
}

// 分页工具
class Pagination {
    static render(container, currentPage, totalPages, onPageChange) {
        if (typeof container === 'string') {
            container = document.querySelector(container);
        }

        if (!container || totalPages <= 1) {
            container.innerHTML = '';
            return;
        }

        let html = '';

        // 上一页
        if (currentPage > 1) {
            html += `<li class="page-item"><a class="page-link" href="#" data-page="${currentPage - 1}">上一页</a></li>`;
        }

        // 页码
        for (let i = 1; i <= totalPages; i++) {
            if (i === 1 || i === totalPages || (i >= currentPage - 2 && i <= currentPage + 2)) {
                const active = i === currentPage ? 'active' : '';
                html += `<li class="page-item ${active}"><a class="page-link" href="#" data-page="${i}">${i}</a></li>`;
            } else if (i === currentPage - 3 || i === currentPage + 3) {
                html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
            }
        }

        // 下一页
        if (currentPage < totalPages) {
            html += `<li class="page-item"><a class="page-link" href="#" data-page="${currentPage + 1}">下一页</a></li>`;
        }

        container.innerHTML = html;

        // 绑定事件
        container.querySelectorAll('a.page-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const page = parseInt(e.target.dataset.page);
                if (page && onPageChange) {
                    onPageChange(page);
                }
            });
        });
    }
}

// 导出工具类
window.Toast = Toast;
window.API = API;
window.DateFormatter = DateFormatter;
window.Validator = Validator;
window.Confirm = Confirm;
window.Loading = Loading;
window.TableHelper = TableHelper;
window.Pagination = Pagination;

// 移动端检测
class DeviceDetector {
    static isMobile() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    }

    static isTablet() {
        return /iPad|Android/i.test(navigator.userAgent) && window.innerWidth >= 768;
    }

    static isDesktop() {
        return !this.isMobile() && !this.isTablet();
    }

    static getDeviceType() {
        if (this.isMobile()) return 'mobile';
        if (this.isTablet()) return 'tablet';
        return 'desktop';
    }

    static getScreenSize() {
        return {
            width: window.innerWidth,
            height: window.innerHeight
        };
    }
}

// 响应式工具
class ResponsiveHelper {
    static init() {
        // 添加设备类型到 body
        document.body.classList.add(`device-${DeviceDetector.getDeviceType()}`);

        // 监听窗口大小变化
        window.addEventListener('resize', this.handleResize.bind(this));

        // 初始化移动端工具栏
        if (DeviceDetector.isMobile()) {
            document.body.classList.add('has-mobile-toolbar');
        }
    }

    static handleResize() {
        const deviceType = DeviceDetector.getDeviceType();
        document.body.className = document.body.className.replace(/device-\w+/, `device-${deviceType}`);
    }

    static optimizeTableForMobile() {
        if (!DeviceDetector.isMobile()) return;

        const tables = document.querySelectorAll('.table-responsive');
        tables.forEach(table => {
            // 添加触摸滚动提示
            if (!table.querySelector('.scroll-hint')) {
                const hint = document.createElement('div');
                hint.className = 'scroll-hint text-muted small text-center py-2';
                hint.innerHTML = '← 左右滑动查看更多 →';
                table.parentNode.insertBefore(hint, table);
            }
        });
    }

    static addPullToRefresh(callback) {
        if (!DeviceDetector.isMobile()) return;

        let startY = 0;
        let isPulling = false;

        document.addEventListener('touchstart', (e) => {
            if (window.scrollY === 0) {
                startY = e.touches[0].pageY;
                isPulling = true;
            }
        });

        document.addEventListener('touchmove', (e) => {
            if (!isPulling) return;

            const currentY = e.touches[0].pageY;
            const diff = currentY - startY;

            if (diff > 100) {
                isPulling = false;
                if (callback) callback();
            }
        });

        document.addEventListener('touchend', () => {
            isPulling = false;
        });
    }
}

// 触摸手势支持
class TouchGestures {
    static enableSwipe(element, onSwipeLeft, onSwipeRight) {
        if (!DeviceDetector.isMobile()) return;

        let startX = 0;
        let startY = 0;

        element.addEventListener('touchstart', (e) => {
            startX = e.touches[0].pageX;
            startY = e.touches[0].pageY;
        });

        element.addEventListener('touchend', (e) => {
            const endX = e.changedTouches[0].pageX;
            const endY = e.changedTouches[0].pageY;
            const diffX = endX - startX;
            const diffY = endY - startY;

            // 确保是水平滑动
            if (Math.abs(diffX) > Math.abs(diffY) && Math.abs(diffX) > 50) {
                if (diffX > 0 && onSwipeRight) {
                    onSwipeRight();
                } else if (diffX < 0 && onSwipeLeft) {
                    onSwipeLeft();
                }
            }
        });
    }
}

// 性能优化
class PerformanceOptimizer {
    static debounce(func, wait = 300) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    static throttle(func, limit = 300) {
        let inThrottle;
        return function executedFunction(...args) {
            if (!inThrottle) {
                func(...args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }

    static lazyLoadImages() {
        const images = document.querySelectorAll('img[data-src]');
        const imageObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.removeAttribute('data-src');
                    imageObserver.unobserve(img);
                }
            });
        });

        images.forEach(img => imageObserver.observe(img));
    }
}

// 页面初始化
document.addEventListener('DOMContentLoaded', () => {
    // 初始化响应式助手
    ResponsiveHelper.init();

    // 优化移动端表格
    ResponsiveHelper.optimizeTableForMobile();

    // 懒加载图片
    PerformanceOptimizer.lazyLoadImages();

    // 创建 Toast 容器
    if (!document.getElementById('toastContainer')) {
        Toast.createToastContainer();
    }
});

// 导出新增工具类
window.DeviceDetector = DeviceDetector;
window.ResponsiveHelper = ResponsiveHelper;
window.TouchGestures = TouchGestures;
window.PerformanceOptimizer = PerformanceOptimizer;
