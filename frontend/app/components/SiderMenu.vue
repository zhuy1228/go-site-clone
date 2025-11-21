<template>
  <div>
    <div class="menu-out">
      <!-- 折叠按钮 - 移到外层以避免被遮挡 -->
      <div class="collapse-trigger" @click="collapsed = !collapsed">
        <menu-unfold-outlined v-if="collapsed" class="trigger-icon" />
        <menu-fold-outlined v-else class="trigger-icon" />
      </div>

      <a-layout-sider 
        v-model:collapsed="collapsed" 
        :trigger="null"
        :width="240"
        class="modern-sider"
      >
        <!-- Logo区域 -->
        <div class="logo-container">
          <div class="logo-icon">
            <cloud-server-outlined class="icon" />
          </div>
          <transition name="fade">
            <div v-if="!collapsed" class="logo-text">
              <div class="app-title">整站克隆</div>
              <div class="app-subtitle">Site Clone</div>
            </div>
          </transition>
        </div>

        <!-- 菜单 -->
        <a-menu 
          v-model:selectedKeys="selectedKeys" 
          mode="inline" 
          class="modern-menu"
          @click="menuClick"
        >
          <a-menu-item key="" class="menu-item">
            <template #icon>
              <home-outlined class="menu-icon" />
            </template>
            <span>控制台</span>
          </a-menu-item>

          <a-menu-divider class="menu-divider" />

          <a-menu-item-group title="站点管理" class="menu-group">
            <a-menu-item key="site/download" class="menu-item">
              <template #icon>
                <cloud-download-outlined class="menu-icon" />
              </template>
              <span>整站下载</span>
            </a-menu-item>
            <a-menu-item key="site/record" class="menu-item">
              <template #icon>
                <history-outlined class="menu-icon" />
              </template>
              <span>下载记录</span>
            </a-menu-item>
          </a-menu-item-group>

          <a-menu-divider class="menu-divider" />

          <a-menu-item-group title="部署运行" class="menu-group">
            <a-menu-item key="run/webpage" class="menu-item">
              <template #icon>
                <global-outlined class="menu-icon" />
              </template>
              <span>Nginx网站</span>
            </a-menu-item>
            <a-menu-item key="run/pack" class="menu-item">
              <template #icon>
                <appstore-outlined class="menu-icon" />
              </template>
              <span>打包应用</span>
            </a-menu-item>
          </a-menu-item-group>

          <a-menu-divider class="menu-divider" />

          <a-menu-item key="settings" class="menu-item">
            <template #icon>
              <setting-outlined class="menu-icon" />
            </template>
            <span>系统设置</span>
          </a-menu-item>
        </a-menu>

        <!-- 底部信息 -->
        <div class="sider-footer">
          <transition name="fade">
            <div v-if="!collapsed" class="footer-content">
              <div class="version-info">
                <api-outlined class="footer-icon" />
                <span>v1.0.0</span>
              </div>
            </div>
          </transition>
        </div>
      </a-layout-sider>
    </div>
    <div>
      <a-layout-sider 
        v-model:collapsed="collapsed" 
        :trigger="null"
        :width="240"
        style="height: 100%; visibility: hidden;"
      ></a-layout-sider>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  HomeOutlined,
  CloudDownloadOutlined,
  HistoryOutlined,
  GlobalOutlined,
  AppstoreOutlined,
  CloudServerOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ApiOutlined,
  SettingOutlined
} from '@ant-design/icons-vue';
import { ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const collapsed = ref<boolean>(false);
const selectedKeys = ref<string[]>(['']);
const route = useRoute();

// 监听路由变化更新选中状态
watch(() => route.path, (newPath) => {
  updateSelectedKeys(newPath);
}, { immediate: true });

onMounted(() => {
  updateSelectedKeys(route.path);
});

function updateSelectedKeys(path: string) {
  // 移除开头的斜杠
  const cleanPath = path.startsWith('/') ? path.slice(1) : path;
  selectedKeys.value = cleanPath ? [cleanPath] : [''];
}

function menuClick(e: any) {
  const key = e.key;
  if (key) {
    navigateTo('/' + key);
  } else {
    navigateTo('/');
  }
}
</script>

<style scoped>
.menu-out {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 1000;
  min-height: 100vh;
  overflow: visible;
}

.modern-sider {
  background: linear-gradient(180deg, #1a1f35 0%, #0f1419 100%) !important;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.15);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: visible;
  height: 100%;
}

:deep(.ant-layout-sider-children) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Logo区域 */
.logo-container {
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  margin-bottom: 8px;
  background: rgba(255, 255, 255, 0.02);
}

.logo-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.logo-icon .icon {
  font-size: 24px;
  color: white;
}

.logo-text {
  flex: 1;
  min-width: 0;
}

.app-title {
  font-size: 18px;
  font-weight: 700;
  color: white;
  line-height: 1.2;
  margin-bottom: 2px;
  letter-spacing: 0.5px;
}

.app-subtitle {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 500;
}

/* 折叠按钮 */
.collapse-trigger {
  position: absolute;
  top: 32px;
  right: -14px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 1001;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
}

.collapse-trigger:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.trigger-icon {
  color: white;
  font-size: 14px;
}

/* 菜单样式 */
.modern-menu {
  background: transparent !important;
  border: none !important;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 12px;
}

/* 隐藏滚动条但保持功能 */
.modern-menu::-webkit-scrollbar {
  width: 4px;
}

.modern-menu::-webkit-scrollbar-track {
  background: transparent;
}

.modern-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

.modern-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

:deep(.ant-menu-item-group-title) {
  color: rgba(255, 255, 255, 0.45) !important;
  font-size: 12px !important;
  font-weight: 600 !important;
  padding: 16px 16px 8px !important;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

:deep(.menu-item) {
  margin: 4px 0 !important;
  border-radius: 10px !important;
  height: 44px !important;
  line-height: 44px !important;
  color: rgba(255, 255, 255, 0.75) !important;
  transition: all 0.3s ease !important;
  font-size: 14px !important;
  font-weight: 500 !important;
}

:deep(.menu-item:hover) {
  background: rgba(255, 255, 255, 0.08) !important;
  color: white !important;
}

:deep(.menu-item.ant-menu-item-selected) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%) !important;
  color: white !important;
  position: relative;
}

:deep(.menu-item.ant-menu-item-selected::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 60%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 0 2px 2px 0;
}

:deep(.menu-icon) {
  font-size: 18px !important;
}

:deep(.menu-divider) {
  background: rgba(255, 255, 255, 0.08) !important;
  margin: 16px 16px !important;
}

/* 底部信息 */
.sider-footer {
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(0, 0, 0, 0.2);
}

.footer-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

.version-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
  font-weight: 500;
}

.footer-icon {
  font-size: 14px;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 折叠状态调整 */
:deep(.ant-layout-sider-collapsed) .logo-container {
  padding: 24px 16px;
  justify-content: center;
}

:deep(.ant-layout-sider-collapsed) .modern-menu {
  padding: 8px 8px;
}

:deep(.ant-layout-sider-collapsed .ant-menu-item-group-title) {
  padding: 8px 0 !important;
  text-align: center;
  opacity: 0;
}

:deep(.ant-layout-sider-collapsed .menu-item) {
  padding: 0 calc(50% - 18px / 2) !important;
}
</style>