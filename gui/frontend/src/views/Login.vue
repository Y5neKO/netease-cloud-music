<script setup>
import QrcodeLogin from '../components/QrcodeLogin.vue'
import PhoneLogin from '../components/PhoneLogin.vue'
import CookieLogin from '../components/CookieLogin.vue'
import CCloudLogin from '../components/CCloudLogin.vue'
import { useUserStore } from '../stores/user'
import { ref } from 'vue'

const userStore = useUserStore()
const activeTab = ref('qrcode')

const tabs = [
  { key: 'qrcode', label: '扫码' },
  { key: 'phone', label: '手机号' },
  { key: 'cookie', label: 'Cookie' },
  { key: 'cc', label: 'CookieCloud' },
]
</script>

<template>
  <div class="page">
    <div class="container">
      <!-- 头部 -->
      <div class="header">
        <svg class="logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/>
          <path d="M9 9.5a3 3 0 0 1 5.23-2"/>
          <path d="M12 17a5 5 0 0 0 4.58-3"/>
          <circle cx="15" cy="8.5" r="1" fill="currentColor" stroke="none"/>
        </svg>
        <h1 class="title">网易云音乐</h1>
      </div>

      <!-- 已登录 -->
      <div v-if="userStore.loggedIn && userStore.profile" class="profile">
        <img :src="userStore.profile.avatarUrl" class="avatar" />
        <div class="profile-info">
          <span class="nickname">{{ userStore.profile.nickname }}</span>
          <span class="badge" v-if="userStore.profile.vipType > 0">VIP</span>
        </div>
        <button class="btn-logout" @click="userStore.logout">退出登录</button>
      </div>

      <!-- 未登录 -->
      <template v-else>
        <!-- Tab 栏 -->
        <div class="tab-bar">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            :class="['tab-item', { active: activeTab === tab.key }]"
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </button>
          <div class="tab-indicator" :style="{ transform: `translateX(${tabs.findIndex(t => t.key === activeTab) * 100}%)` }" />
        </div>

        <!-- 内容区 -->
        <div class="tab-content">
          <QrcodeLogin v-if="activeTab === 'qrcode'" />
          <PhoneLogin v-else-if="activeTab === 'phone'" />
          <CookieLogin v-else-if="activeTab === 'cookie'" />
          <CCloudLogin v-else-if="activeTab === 'cc'" />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
  padding: 24px;
}

.container {
  width: 100%;
  max-width: 380px;
  background: var(--surface);
  border-radius: 14px;
  box-shadow: var(--shadow-md);
  border: 0.5px solid var(--border);
  overflow: hidden;
}

.header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 28px 24px 20px;
}

.logo {
  width: 28px;
  height: 28px;
  color: var(--red);
}

.title {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.3px;
  color: var(--text-primary);
}

/* Profile */
.profile {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 24px 28px;
  gap: 12px;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
}

.profile-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nickname {
  font-size: 17px;
  font-weight: 600;
}

.badge {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
  background: linear-gradient(135deg, #f6d365, #fda085);
  color: #7c4700;
}

.btn-logout {
  margin-top: 4px;
  padding: 6px 20px;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(255, 59, 48, 0.2);
  background: none;
  color: var(--red);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}
.btn-logout:hover {
  background: rgba(255, 59, 48, 0.06);
}

/* Tab Bar */
.tab-bar {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  position: relative;
  margin: 0 20px;
  background: rgba(118, 118, 128, 0.08);
  border-radius: 8px;
  padding: 2px;
}

.tab-item {
  position: relative;
  z-index: 1;
  padding: 6px 0;
  border: none;
  background: none;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s;
}
.tab-item.active {
  color: var(--text-primary);
}

.tab-indicator {
  position: absolute;
  top: 2px;
  left: 2px;
  width: calc(25% - 2px);
  height: calc(100% - 4px);
  background: var(--surface);
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  transition: transform 0.25s cubic-bezier(0.25, 0.1, 0.25, 1);
}

.tab-content {
  padding: 20px 24px 24px;
}
</style>
