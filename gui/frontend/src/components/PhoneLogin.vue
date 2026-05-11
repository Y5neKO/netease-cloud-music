<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { SendSmsCode, LoginByPhone, SaveConfig, LoadConfig } from '../../bindings/auth'

const userStore = useUserStore()

const phone = ref('')
const password = ref('')
const captcha = ref('')
const loading = ref(false)
const msg = ref('')
const smsCountdown = ref(0)
const rememberPwd = ref(false)
const mode = ref<'sms' | 'password'>('password')
let smsTimer: number | null = null

const phoneValid = computed(() => /^\d{5,15}$/.test(phone.value))
const msgIsError = computed(() => msg.value.includes('失败') || msg.value.includes('错误'))

onMounted(async () => {
  try {
    const raw = await LoadConfig('phone_login')
    if (raw) {
      const data = JSON.parse(raw)
      phone.value = data.phone || ''
      password.value = data.password || ''
      rememberPwd.value = !!data.password
    }
  } catch {}
})

async function saveLogin() {
  const data: any = { phone: phone.value }
  if (rememberPwd.value && password.value) {
    data.password = password.value
  }
  try { await SaveConfig('phone_login', JSON.stringify(data)) } catch {}
}

async function sendSms() {
  if (!phoneValid.value) { msg.value = '请输入正确的手机号'; return }
  try {
    loading.value = true
    msg.value = ''
    await SendSmsCode(phone.value, 86)
    msg.value = '验证码已发送'
    smsCountdown.value = 60
    smsTimer = window.setInterval(() => {
      smsCountdown.value--
      if (smsCountdown.value <= 0 && smsTimer) { clearInterval(smsTimer); smsTimer = null }
    }, 1000)
  } catch (e: any) {
    msg.value = '发送失败: ' + e
  } finally {
    loading.value = false
  }
}

async function loginWithSms() {
  if (!phoneValid.value || !captcha.value) { msg.value = '请填写手机号和验证码'; return }
  try {
    loading.value = true; msg.value = ''
    saveLogin()
    const profile = await LoginByPhone(phone.value, '', captcha.value, 86)
    userStore.onLoginSuccess(profile)
  } catch (e: any) {
    msg.value = '登录失败: ' + e
  } finally {
    loading.value = false
  }
}

async function loginWithPassword() {
  if (!phoneValid.value || !password.value) { msg.value = '请填写手机号和密码'; return }
  try {
    loading.value = true; msg.value = ''
    saveLogin()
    const profile = await LoginByPhone(phone.value, password.value, '', 86)
    userStore.onLoginSuccess(profile)
  } catch (e: any) {
    msg.value = '登录失败: ' + e
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="phone-login">
    <div class="field">
      <label class="label">手机号</label>
      <input v-model="phone" type="tel" class="input" placeholder="请输入手机号" maxlength="15" />
    </div>

    <!-- 模式切换 -->
    <div class="mode-switch">
      <button :class="['mode-btn', { active: mode === 'sms' }]" @click="mode = 'sms'">验证码</button>
      <button :class="['mode-btn', { active: mode === 'password' }]" @click="mode = 'password'">密码</button>
    </div>

    <!-- 验证码模式 -->
    <template v-if="mode === 'sms'">
      <div class="field">
        <label class="label">验证码</label>
        <div class="input-group">
          <input v-model="captcha" type="text" class="input" placeholder="短信验证码" maxlength="6" />
          <button class="btn-sms" @click="sendSms" :disabled="smsCountdown > 0 || loading">
            {{ smsCountdown > 0 ? `${smsCountdown}s` : '获取' }}
          </button>
        </div>
      </div>
      <button class="btn-primary" @click="loginWithSms" :disabled="loading">登录</button>
    </template>

    <!-- 密码模式 -->
    <template v-else>
      <div class="field">
        <label class="label">密码</label>
        <input v-model="password" type="password" class="input" placeholder="请输入密码" />
      </div>
      <div class="remember-row">
        <label class="checkbox-label">
          <input type="checkbox" v-model="rememberPwd" class="checkbox" />
          <span>记住密码</span>
        </label>
      </div>
      <button class="btn-primary" @click="loginWithPassword" :disabled="loading">登录</button>
    </template>

    <p v-if="msg" class="msg" :class="{ error: msgIsError, success: !msgIsError }">{{ msg }}</p>
  </div>
</template>

<style scoped>
.phone-login {
  display: flex;
  flex-direction: column;
  gap: 14px;
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

.input-group {
  display: flex;
  gap: 8px;
}
.input-group .input {
  flex: 1;
}

.btn-sms {
  white-space: nowrap;
  padding: 0 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--accent);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-sms:hover:not(:disabled) {
  background: rgba(0, 122, 255, 0.04);
}
.btn-sms:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.mode-switch {
  display: grid;
  grid-template-columns: 1fr 1fr;
  background: rgba(118, 118, 128, 0.08);
  border-radius: 6px;
  padding: 2px;
}

.mode-btn {
  padding: 5px 0;
  border: none;
  background: none;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 5px;
  transition: all 0.15s;
}
.mode-btn.active {
  background: var(--surface);
  color: var(--text-primary);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

.remember-row {
  display: flex;
  align-items: center;
  margin-top: -4px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  user-select: none;
}

.checkbox {
  accent-color: var(--accent);
  width: 14px;
  height: 14px;
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
.btn-primary:hover {
  background: var(--accent-hover);
}
.btn-primary:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.msg {
  font-size: 12px;
  text-align: center;
  min-height: 16px;
}
.msg.error { color: var(--red); }
.msg.success { color: #34c759; }
</style>
