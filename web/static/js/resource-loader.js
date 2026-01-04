/**
 * 资源加载器 - 支持本地资源和CDN回退
 * 优先使用本地资源，失败时回退到CDN
 */

class ResourceLoader {
    constructor() {
        this.loadedResources = new Set();
        this.failedResources = new Set();
    }

    /**
     * 加载CSS资源
     * @param {string} localPath - 本地路径
     * @param {string} cdnUrl - CDN备用地址
     * @param {string} id - 资源ID
     */
    async loadCSS(localPath, cdnUrl, id) {
        if (this.loadedResources.has(id)) {
            return Promise.resolve();
        }

        return new Promise((resolve, reject) => {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.id = id;

            const tryLoad = (url, isFallback = false) => {
                link.href = url;

                link.onload = () => {
                    this.loadedResources.add(id);
                    if (isFallback) {
                        console.warn(`🔄 CSS回退到CDN: ${id} -> ${url}`);
                    } else {
                        console.log(`✅ CSS加载成功: ${id} -> ${url}`);
                    }
                    resolve();
                };

                link.onerror = () => {
                    if (!isFallback && cdnUrl) {
                        console.warn(`⚠️ 本地CSS加载失败，尝试CDN: ${id}`);
                        tryLoad(cdnUrl, true);
                    } else {
                        this.failedResources.add(id);
                        console.error(`❌ CSS加载失败: ${id}`);
                        reject(new Error(`Failed to load CSS: ${id}`));
                    }
                };

                document.head.appendChild(link);
            };

            tryLoad(localPath);
        });
    }

    /**
     * 加载JavaScript资源
     * @param {string} localPath - 本地路径
     * @param {string} cdnUrl - CDN备用地址
     * @param {string} id - 资源ID
     * @param {function} testFunction - 测试函数，用于验证库是否正确加载
     */
    async loadJS(localPath, cdnUrl, id, testFunction) {
        if (this.loadedResources.has(id)) {
            return Promise.resolve();
        }

        return new Promise((resolve, reject) => {
            const script = document.createElement('script');
            script.id = id;

            const tryLoad = (url, isFallback = false) => {
                script.src = url;

                script.onload = () => {
                    // 如果提供了测试函数，验证库是否正确加载
                    if (testFunction && !testFunction()) {
                        if (!isFallback && cdnUrl) {
                            console.warn(`⚠️ JS库验证失败，尝试CDN: ${id}`);
                            document.head.removeChild(script);
                            const newScript = document.createElement('script');
                            newScript.id = id;
                            tryLoad(cdnUrl, true);
                            return;
                        } else {
                            this.failedResources.add(id);
                            console.error(`❌ JS库验证失败: ${id}`);
                            reject(new Error(`JS library validation failed: ${id}`));
                            return;
                        }
                    }

                    this.loadedResources.add(id);
                    if (isFallback) {
                        console.warn(`🔄 JS回退到CDN: ${id} -> ${url}`);
                    } else {
                        console.log(`✅ JS加载成功: ${id} -> ${url}`);
                    }
                    resolve();
                };

                script.onerror = () => {
                    if (!isFallback && cdnUrl) {
                        console.warn(`⚠️ 本地JS加载失败，尝试CDN: ${id}`);
                        document.head.removeChild(script);
                        const newScript = document.createElement('script');
                        newScript.id = id;
                        Object.assign(newScript, script);
                        tryLoad(cdnUrl, true);
                    } else {
                        this.failedResources.add(id);
                        console.error(`❌ JS加载失败: ${id}`);
                        reject(new Error(`Failed to load JS: ${id}`));
                    }
                };

                document.head.appendChild(script);
            };

            tryLoad(localPath);
        });
    }

    /**
     * 批量加载资源
     * @param {Array} resources - 资源配置数组
     */
    async loadResources(resources) {
        const promises = resources.map(resource => {
            if (resource.type === 'css') {
                return this.loadCSS(resource.local, resource.cdn, resource.id);
            } else if (resource.type === 'js') {
                return this.loadJS(resource.local, resource.cdn, resource.id, resource.test);
            }
        });

        try {
            await Promise.all(promises);
            console.log('🎉 所有资源加载完成');
        } catch (error) {
            console.error('❌ 部分资源加载失败:', error);
            throw error;
        }
    }

    /**
     * 获取加载状态报告
     */
    getLoadReport() {
        return {
            loaded: Array.from(this.loadedResources),
            failed: Array.from(this.failedResources),
            total: this.loadedResources.size + this.failedResources.size
        };
    }
}

// 全局实例
window.resourceLoader = new ResourceLoader();