<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { NInput, NButton, NSpace, NText } from 'naive-ui'
import { useUserStore } from '../stores/user'

import { LoginByCookieCloud } from '../../bindings/auth'

const userStore = useUserStore()

const server = ref('http://127.0.0.1:8088')
const uuid = ref('')
const password = ref('')
const loading = ref(false)
const msg = ref('')

onMounted(() => {
  const saved = localStorage.getItem('ncm.cookiecloud')
  if (saved) {
    try {
      const data = JSON.parse(saved)
      server.value = data.server || 'http://127.0.0.1:8088'
      uuid.value = data.uuid || ''
      password.value = data.password || ''
    } catch {}
  }
})

async function login() {
  if (!uuid.value || !password.value) {
    msg.value = '请填写 UUID 和密码'
    return
  }
  try {
    loading.value = true
    msg.value = '连接 CookieCloud...'
    localStorage.setItem('ncm.cookiecloud', JSON.stringify({
      server: server.value,
      uuid: uuid.value,
      password: password.value
    }))
    const profile = await LoginByCookieCloud(server.value, uuid.value, password.value)
    userStore.onLoginSuccess(profile)
  } catch (e: any) {
    msg.value = '登录失败: ' + e
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <NSpace vertical :size="16">
    <NText depth="3">
      从 CookieCloud 浏览器插件同步网易云音乐 Cookie
    </NText>
    <NInput v-model:value="server" placeholder="CookieCloud 服务器地址" />
    <NInput v-model:value="uuid" placeholder="UUID" />
    <NInput
      v-model:value="password"
      type="password"
      showPasswordOn="click"
      placeholder="密码"
    />
    <NButton type="primary" block @click="login" :loading="loading">
      登录
    </NButton>
    <NText v-if="msg" :type="msg.includes('失败') ? 'error' : 'success'">
      {{ msg }}
    </NText>
  </NSpace>
</template>
