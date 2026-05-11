<script setup>
import { NTabs, NTabPane, NCard, NSpace, NText, NAvatar, NButton } from 'naive-ui'
import QrcodeLogin from '../components/QrcodeLogin.vue'
import PhoneLogin from '../components/PhoneLogin.vue'
import CookieLogin from '../components/CookieLogin.vue'
import CCloudLogin from '../components/CCloudLogin.vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
</script>

<template>
  <div class="login-page">
    <NCard class="login-card" :bordered="false">
      <NSpace vertical align="center" :size="20">
        <span class="title">网易云音乐</span>

        <template v-if="userStore.loggedIn && userStore.profile">
          <NAvatar :size="64" :src="userStore.profile.avatarUrl" round />
          <NText strong style="font-size: 18px">{{ userStore.profile.nickname }}</NText>
          <NText depth="3">VIP: {{ userStore.profile.vipType > 0 ? '是' : '否' }}</NText>
          <NButton type="error" ghost @click="userStore.logout">退出登录</NButton>
        </template>

        <template v-else>
          <NTabs type="line" animated style="width: 100%">
            <NTabPane name="qrcode" tab="扫码登录">
              <QrcodeLogin />
            </NTabPane>
            <NTabPane name="phone" tab="手机登录">
              <PhoneLogin />
            </NTabPane>
            <NTabPane name="cookie" tab="Cookie">
              <CookieLogin />
            </NTabPane>
            <NTabPane name="cc" tab="CookieCloud">
              <CCloudLogin />
            </NTabPane>
          </NTabs>
        </template>
      </NSpace>
    </NCard>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #c62f2f 100%);
}
.title {
  font-size: 24px;
  font-weight: bold;
  color: #c62f2f;
}
.login-card {
  width: 420px;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}
</style>
