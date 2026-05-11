<script lang="ts" setup>
import { ref } from 'vue'
import { NInput, NButton, NSpace, NText } from 'naive-ui'
import { useUserStore } from '../stores/user'

import { LoginByCookie } from '../../bindings/auth'

const userStore = useUserStore()

const cookieStr = ref('')
const loading = ref(false)
const msg = ref('')

async function login() {
  if (!cookieStr.value.trim()) {
    msg.value = '请输入 Cookie'
    return
  }
  try {
    loading.value = true
    msg.value = '验证中...'
    const profile = await LoginByCookie(cookieStr.value.trim())
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
      从浏览器中复制网易云音乐站点的 Cookie（需包含 MUSIC_U 字段）
    </NText>
    <NInput
      v-model:value="cookieStr"
      type="textarea"
      placeholder="粘贴 Cookie 字符串..."
      :rows="4"
    />
    <NButton type="primary" block @click="login" :loading="loading">
      登录
    </NButton>
    <NText v-if="msg" :type="msg.includes('失败') ? 'error' : 'success'">
      {{ msg }}
    </NText>
  </NSpace>
</template>
