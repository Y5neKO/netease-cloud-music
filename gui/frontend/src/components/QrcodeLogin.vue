<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { NButton, NImage, NSpace, NText, NSpin } from 'naive-ui'
import { useUserStore } from '../stores/user'
import { CreateQrcode, CheckQrcodeStatus } from '../../bindings/auth'

const userStore = useUserStore()

const qrImage = ref('')
const qrKey = ref('')
const status = ref('点击下方按钮生成二维码')
const loading = ref(false)
const polling = ref(null)

async function generate() {
  loading.value = true
  status.value = '生成中...'
  try {
    const result = await CreateQrcode()
    qrKey.value = result.key
    qrImage.value = 'data:image/png;base64,' + result.image
    status.value = '请使用网易云音乐 App 扫描二维码'
    startPolling()
  } catch (e) {
    status.value = '生成失败: ' + e
  } finally {
    loading.value = false
  }
}

function startPolling() {
  stopPolling()
  polling.value = window.setInterval(async () => {
    try {
      const s = await CheckQrcodeStatus(qrKey.value)
      switch (s.code) {
        case 800: status.value = '二维码已过期，请刷新'; stopPolling(); break
        case 801: status.value = '等待扫码...'; break
        case 802: status.value = '已扫码，请在手机上确认'; break
        case 803: status.value = '登录成功！'; stopPolling(); await userStore.fetchProfile(); break
      }
    } catch { stopPolling() }
  }, 3000)
}

function stopPolling() {
  if (polling.value !== null) { clearInterval(polling.value); polling.value = null }
}

onMounted(() => generate())
onUnmounted(() => stopPolling())
</script>

<template>
  <NSpace vertical align="center" :size="20">
    <NSpin :show="loading">
      <NImage v-if="qrImage" :src="qrImage" :width="200" :height="200" preview-disabled />
    </NSpin>
    <NText>{{ status }}</NText>
    <NButton @click="generate" :loading="loading" type="primary" ghost>刷新二维码</NButton>
  </NSpace>
</template>
