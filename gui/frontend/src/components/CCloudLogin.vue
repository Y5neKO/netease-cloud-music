<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { LoginByCookieCloud, SaveConfig, LoadConfig } from '../../bindings/auth'

const userStore = useUserStore()

const server = ref('http://127.0.0.1:8088')
const uuid = ref('')
const password = ref('')
const loading = ref(false)
const msg = ref('')
const msgIsError = computed(() => msg.value.includes('失败'))

onMounted(async () => {
  try {
    const raw = await LoadConfig('cookiecloud')
    if (raw) {
      const data = JSON.parse(raw)
      server.value = data.server || 'http://127.0.0.1:8088'
      uuid.value = data.uuid || ''
      password.value = data.password || ''
    }
  } catch {}
})

async function login() {
  if (!uuid.value || !password.value) { msg.value = '请填写 UUID 和密码'; return }
  try {
    loading.value = true; msg.value = ''
    await SaveConfig('cookiecloud', JSON.stringify({
      server: server.value, uuid: uuid.value, password: password.value
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
  <div class="cc-login">
    <p class="hint">从 CookieCloud 浏览器插件同步网易云音乐 Cookie</p>

    <div class="field">
      <label class="label">服务器</label>
      <input v-model="server" type="url" class="input" placeholder="http://127.0.0.1:8088" />
    </div>

    <div class="field">
      <label class="label">UUID</label>
      <input v-model="uuid" type="text" class="input" placeholder="用户 UUID" />
    </div>

    <div class="field">
      <label class="label">密码</label>
      <input v-model="password" type="password" class="input" placeholder="加密密码" />
    </div>

    <button class="btn-primary" @click="login" :disabled="loading">同步并登录</button>
    <p v-if="msg" class="msg" :class="{ error: msgIsError, success: !msgIsError }">{{ msg }}</p>
  </div>
</template>

<style scoped>
.cc-login {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hint {
  font-size: 12px;
  color: var(--text-tertiary);
  line-height: 1.4;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  letter-spacing: 0.2px;
}

.input {
  width: 100%;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg);
  font-size: 14px;
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.12);
}
.input::placeholder {
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
