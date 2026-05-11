<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '../stores/user'
import { CreateQrcode, CheckQrcodeStatus } from '../../bindings/auth'

const userStore = useUserStore()

const qrImage = ref('')
const qrKey = ref('')
const status = ref('')
const loading = ref(false)
let polling = null

async function generate() {
  loading.value = true
  status.value = ''
  try {
    const result = await CreateQrcode()
    qrKey.value = result.key
    qrImage.value = 'data:image/png;base64,' + result.image
    status.value = '请使用网易云音乐 App 扫码'
    startPolling()
  } catch (e) {
    status.value = '生成失败: ' + e
  } finally {
    loading.value = false
  }
}

function startPolling() {
  stopPolling()
  polling = window.setInterval(async () => {
    try {
      const s = await CheckQrcodeStatus(qrKey.value)
      switch (s.code) {
        case 800: status.value = '二维码已过期'; stopPolling(); break
        case 801: status.value = '等待扫码'; break
        case 802: status.value = '已扫码，请在手机确认'; break
        case 803: status.value = '登录成功'; stopPolling(); await userStore.fetchProfile(); break
      }
    } catch { stopPolling() }
  }, 3000)
}

function stopPolling() {
  if (polling !== null) { clearInterval(polling); polling = null }
}

onMounted(() => generate())
onUnmounted(() => stopPolling())
</script>

<template>
  <div class="qr-login">
    <div class="qr-wrapper">
      <div v-if="loading" class="qr-loading">
        <div class="spinner" />
      </div>
      <img v-else-if="qrImage" :src="qrImage" class="qr-image" />
      <div v-else class="qr-placeholder">
        <span>点击下方按钮生成二维码</span>
      </div>
    </div>
    <p class="status" :class="{ error: status.includes('过期') || status.includes('失败') }">{{ status }}</p>
    <button class="btn-secondary" @click="generate" :disabled="loading">
      {{ qrImage ? '刷新二维码' : '生成二维码' }}
    </button>
  </div>
</template>

<style scoped>
.qr-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.qr-wrapper {
  width: 180px;
  height: 180px;
  border-radius: 8px;
  border: 1px solid var(--divider);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
}

.qr-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.qr-loading, .qr-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.qr-placeholder {
  color: var(--text-tertiary);
  font-size: 13px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--divider);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status {
  font-size: 13px;
  color: var(--text-secondary);
  text-align: center;
  min-height: 18px;
}
.status.error {
  color: var(--red);
}

.btn-secondary {
  padding: 7px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--accent);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-secondary:hover {
  background: rgba(0, 122, 255, 0.04);
  border-color: rgba(0, 122, 255, 0.3);
}
.btn-secondary:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
