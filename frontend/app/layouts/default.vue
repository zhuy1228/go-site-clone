

<template>
  <div class="app-layout">
    <a-layout class="main-layout">
      <SiderMenu />
      <a-layout class="content-layout">
        <Header />
        <a-layout-content class="page-content">
          <div class="content-wrapper">
            <slot />
          </div>
        </a-layout-content>
        <a-layout-footer class="app-footer">
          <div class="footer-content">
            <div class="footer-text">
              <span class="copyright">© 2024 整站克隆工具</span>
              <span class="separator">•</span>
              <span class="version">Version 1.0.0</span>
            </div>
            <div class="footer-links">
              <a href="#" class="footer-link">使用文档</a>
              <span class="separator">•</span>
              <a href="#" class="footer-link">意见反馈</a>
              <span class="separator">•</span>
              <a href="#" class="footer-link">关于我们</a>
            </div>
          </div>
        </a-layout-footer>
      </a-layout>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

onMounted(() => {
  // 初始化逻辑
});
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: #f0f2f5;
}

.main-layout {
  min-height: 100vh;
  background: transparent;
}

.content-layout {
  background: transparent;
}

.page-content {
  margin: 0;
  padding: 0;
  min-height: calc(100vh - 128px);
  background: transparent;
  overflow-x: hidden;
}

.content-wrapper {
  min-height: calc(100vh - 128px);
}

/* 页面底部 */
.app-footer {
  background: white;
  padding: 20px 24px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  margin-top: auto;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
  flex-wrap: wrap;
  gap: 16px;
}

.footer-text,
.footer-links {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #666;
}

.copyright {
  font-weight: 500;
  color: #333;
}

.version {
  color: #999;
}

.separator {
  color: #d9d9d9;
  user-select: none;
}

.footer-link {
  color: #666;
  text-decoration: none;
  transition: color 0.3s ease;
}

.footer-link:hover {
  color: #1890ff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-footer {
    padding: 16px;
  }

  .footer-content {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }

  .footer-text,
  .footer-links {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>