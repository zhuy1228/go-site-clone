<template>
  <div class="header-outside">
    <a-layout-header class="modern-header">
      <div class="header-content">
        <!-- 左侧面包屑导航 -->
        <div class="header-left">
          <a-breadcrumb class="breadcrumb">
            <a-breadcrumb-item>
              <home-outlined />
              <span>{{ breadcrumbTitle }}</span>
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>

        <!-- 右侧操作区 -->
        <div class="header-right">
          <!-- 搜索框 -->
          <div class="search-box">
            <a-input 
              v-model:value="searchValue" 
              placeholder="搜索功能..."
              class="search-input"
              allow-clear
            >
              <template #prefix>
                <search-outlined class="search-icon" />
              </template>
            </a-input>
          </div>

          <!-- 快捷操作 -->
          <div class="header-actions">
            <a-tooltip title="通知">
              <a-badge :count="3" :offset="[-4, 4]">
                <div class="action-btn">
                  <bell-outlined class="action-icon" />
                </div>
              </a-badge>
            </a-tooltip>

            <a-tooltip title="设置">
              <div class="action-btn">
                <setting-outlined class="action-icon" />
              </div>
            </a-tooltip>

            <a-divider type="vertical" class="action-divider" />

            <!-- 用户信息 -->
            <a-dropdown placement="bottomRight">
              <div class="user-info">
                <a-avatar 
                  class="user-avatar" 
                  :style="{ backgroundColor: '#667eea' }"
                >
                  <template #icon>
                    <user-outlined />
                  </template>
                </a-avatar>
                <span class="user-name">管理员</span>
                <down-outlined class="dropdown-icon" />
              </div>
              <template #overlay>
                <a-menu class="user-menu">
                  <a-menu-item key="profile">
                    <user-outlined />
                    <span>个人中心</span>
                  </a-menu-item>
                  <a-menu-item key="settings">
                    <setting-outlined />
                    <span>系统设置</span>
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item key="logout">
                    <logout-outlined />
                    <span>退出登录</span>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>
    </a-layout-header>
  </div>
</template>

<script lang="ts" setup>
import { 
  HomeOutlined,
  SearchOutlined,
  BellOutlined,
  SettingOutlined,
  UserOutlined,
  DownOutlined,
  LogoutOutlined
} from '@ant-design/icons-vue';
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const searchValue = ref('');

// 根据路由动态生成面包屑标题
const breadcrumbTitle = computed(() => {
  const path = route.path;
  if (path === '/' || path === '') return '控制台';
  if (path.includes('download')) return '整站下载';
  if (path.includes('record')) return '下载记录';
  if (path.includes('webpage')) return 'Nginx网站';
  if (path.includes('pack')) return '打包应用';
  return '控制台';
});
</script>

<style scoped>
.header-outside {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 999;
}

.modern-header {
  background: white;
  padding: 0 24px;
  height: 64px;
  line-height: 64px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

/* 左侧区域 */
.header-left {
  flex: 1;
}

.breadcrumb {
  line-height: 64px;
}

:deep(.breadcrumb .ant-breadcrumb-link) {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

:deep(.breadcrumb .anticon) {
  font-size: 18px;
  color: #1890ff;
}

/* 右侧区域 */
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 搜索框 */
.search-box {
  margin-right: 8px;
}

.search-input {
  width: 240px;
  border-radius: 20px;
  background: #f5f5f5;
  border: 1px solid transparent;
  transition: all 0.3s ease;
}

.search-input:hover,
.search-input:focus {
  background: white;
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

:deep(.search-input .ant-input) {
  background: transparent;
  font-size: 14px;
}

.search-icon {
  color: #999;
}

/* 操作按钮 */
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f5f5f5;
}

.action-btn:hover {
  background: #e6f7ff;
  transform: translateY(-2px);
}

.action-icon {
  font-size: 18px;
  color: #666;
}

.action-btn:hover .action-icon {
  color: #1890ff;
}

.action-divider {
  height: 24px;
  margin: 0 4px;
  background: rgba(0, 0, 0, 0.1);
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f5f5f5;
}

.user-info:hover {
  background: #e6f7ff;
}

.user-avatar {
  flex-shrink: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  line-height: 1;
}

.dropdown-icon {
  font-size: 12px;
  color: #999;
  transition: transform 0.3s ease;
}

.user-info:hover .dropdown-icon {
  transform: rotate(180deg);
}

/* 用户菜单 */
:deep(.user-menu) {
  min-width: 160px;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

:deep(.user-menu .ant-dropdown-menu-item) {
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
}

:deep(.user-menu .ant-dropdown-menu-item:hover) {
  background: #e6f7ff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-header {
    padding: 0 16px;
  }

  .search-box {
    display: none;
  }

  .user-name {
    display: none;
  }

  .user-info {
    padding: 8px;
  }
}
</style>
