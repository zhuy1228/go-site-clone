<template>
  <div>
    <a-breadcrumb style="margin: 10px 0">
      <a-breadcrumb-item>仿站</a-breadcrumb-item>
      <a-breadcrumb-item>整站下载</a-breadcrumb-item>
    </a-breadcrumb>
    <div class="google-like-page">
      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- Logo区域 -->
        <div class="logo-container" :class="{ 'logo-small': isSearching }">
          <h1 class="logo">Google</h1>
        </div>

        <!-- 搜索框容器，使用动态类名控制动画 -->
        <div class="search-container" :class="{ 'search-top': isSearching }" ref="searchContainer">
          <div class="search-box">
            <a-input v-model:value="searchQuery" class="search-input" placeholder="输入一个网址或者包含网址的链接，点击回车"
              @press-enter="handleSearch" ref="searchInput">
              <template #prefix>
                <SearchOutlined />
              </template>
            </a-input>

            <!-- <div class="search-buttons">
              <a-button type="primary" @click="handleSearch">Google 搜索</a-button>
              <a-button @click="handleFeelingLucky">手气不错</a-button>
            </div> -->
          </div>
        </div>

        <!-- 搜索结果区域 -->
        <div v-if="showResults" class="search-results">
          <h2>搜索结果</h2>
          <p>这是搜索结果区域...</p>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
import { SearchOutlined, MenuOutlined } from '@ant-design/icons-vue';

export default {
  components: {
    SearchOutlined,
    MenuOutlined
  },
  data() {
    return {
      searchQuery: '',
      isSearching: false,
      showResults: false
    };
  },
  methods: {
    async handleSearch() {
      if (!this.searchQuery.trim()) return;

      // 触发搜索状态
      this.isSearching = true;

      // 等待下一个 tick 确保 DOM 已更新
      await this.$nextTick();

      // 添加动画类
      if (this.$refs.searchContainer) {
        this.$refs.searchContainer.classList.add('animating');

        // 动画结束后显示结果
        setTimeout(() => {
          if (this.$refs.searchContainer) {
            this.$refs.searchContainer.classList.remove('animating');
          }
          this.showResults = true;

          // 可选：滚动到顶部
          window.scrollTo({ top: 0, behavior: 'smooth' });
        }, 500); // 匹配 CSS 动画时长
      }
    },
    handleFeelingLucky() {
      // 手气不错功能
      window.open('https://www.google.com/doodles', '_blank');
    }
  }
};
</script>

<style scoped>
.google-like-page {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 154px);
  background-color: #fff;
  font-family: Arial, sans-serif;
}

.header {
  display: flex;
  justify-content: flex-end;
  padding: 15px 20px;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 15px;
}

.nav-links a {
  text-decoration: none;
  color: rgba(0, 0, 0, 0.87);
  font-size: 13px;
}

.nav-links a:hover {
  text-decoration: underline;
}

.apps-button {
  color: #5f6368;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  margin-top: 60px;
  /* 为顶部导航留出空间 */
}

.logo-container {
  margin-bottom: 30px;
  transition: all 0.5s ease;
}

.logo {
  font-size: 5rem;
  font-weight: 400;
  color: #4285f4;
  margin: 0;
  transition: all 0.5s ease;
}

.logo-small .logo {
  font-size: 2rem;
  margin-bottom: 10px;
}

.search-container {
  width: 100%;
  max-width: 600px;
  transition: all 0.5s ease;
  display: flex;
  justify-content: center;
}

.search-container.animating {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-top {
  position: fixed;
  top: 180px;
  z-index: 999;
  max-width: 500px;
  transform: translateX(-50%);
  left: 50%;
}

.search-box {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.search-input {
  width: 100%;
  height: 44px;
  border-radius: 24px;
  padding: 0 20px;
  margin-bottom: 20px;
  box-shadow: 0 1px 6px rgba(32, 33, 36, 0.28);
  border: none;
}

.search-input :deep(.ant-input) {
  border-radius: 24px;
  height: 44px;
  border: none;
  outline: none;
}

.search-input :deep(.ant-input):focus {
  box-shadow: none;
}

.search-buttons {
  display: flex;
  gap: 10px;
}

.search-buttons .ant-btn {
  background-color: #f8f9fa;
  border: 1px solid #f8f9fa;
  border-radius: 4px;
  color: #3c4043;
  font-size: 14px;
  padding: 0 16px;
  height: 36px;
}

.search-buttons .ant-btn-primary {
  background-color: #f8f9fa;
  border-color: #f8f9fa;
  color: #3c4043;
}

.search-buttons .ant-btn-primary:hover {
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
  background-color: #f8f9fa;
  border: 1px solid #dadce0;
  color: #202124;
}

.search-results {
  margin-top: 50px;
  width: 100%;
  max-width: 600px;
}

.footer {
  background: #f2f2f2;
  border-top: 1px solid #e4e4e4;
  padding: 15px 20px;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.footer-links a {
  text-decoration: none;
  color: #5f6368;
  font-size: 14px;
}

.footer-links a:hover {
  text-decoration: underline;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .logo {
    font-size: 3rem;
  }

  .logo-small .logo {
    font-size: 1.5rem;
  }

  .search-container {
    max-width: 90%;
  }

  .search-top {
    max-width: 80%;
  }

  .search-buttons {
    flex-direction: column;
    width: 100%;
  }

  .search-buttons .ant-btn {
    width: 100%;
  }
}
</style>
