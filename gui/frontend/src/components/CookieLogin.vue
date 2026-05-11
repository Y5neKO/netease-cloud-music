<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useUserStore } from '../stores/user'
import { LoginByCookie } from '../../bindings/auth'

const userStore = useUserStore()

const cookieStr = ref('')
const loading = ref(false)
const msg = ref('')
const msgIsError = computed(() => msg.value.includes('失败'))

async function login() {
  if (!cookieStr.value.trim()) { msg.value = '请输入 Cookie'; return }
  try {
    loading.value = true; msg.value = ''
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
  <div class="cookie-login">
    <p class="hint">粘贴浏览器中网易云音乐站点的 Cookie（需包含 MUSIC_U）</p>
    <textarea
      v-model="cookieStr"
      class="textarea"
      placeholder="Cookie 字符串..."
      rows="3"
    />
    <button class="btn-primary" @click="login" :disabled="loading">登录</button>
    <p v-if="msg" class="msg" :class="{ error: msgIsError, success: !msgIsError }">{{ msg }}</p>
  </div>
</template>

<style scoped>
.cookie-login {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hint {
  font-size: 12px;
  color: var(--text-tertiary);
  line-height: 1.4;
}

.textarea {
  width: 100%;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg);
  font-size: 13px;
  font-family: 'SF Mono', 'Menlo', monospace;
  color: var(--text-primary);
  outline: none;
  resize: vertical;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.textarea:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.12);
}
.textarea::placeholder {
  color: var(--text-tertiary);
}

.btn-primary {
  width: 100%;
  padding: 9px 0;
  border-radius: var(--radius-sm);
  border: none;
  background: var(--accent);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}
.btn-primary:hover { background: var(--accent-hover); }
.btn-primary:disabled { opacity: 0.45; cursor: not-allowed; }

.msg { font-size: 12px; text-align: center; min-height: 16px; }
.msg.error { color: var(--red); }
.msg.success { color: #34c759; }
</style>
