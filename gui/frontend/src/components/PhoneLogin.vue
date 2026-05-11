<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { NInput, NButton, NSpace, NText, NTabs, NTabPane, NInputGroup, NCheckbox } from 'naive-ui'
import { useUserStore } from '../stores/user'

import { SendSmsCode, LoginByPhone } from '../../bindings/auth'

const userStore = useUserStore()

const phone = ref('')
const password = ref('')
const captcha = ref('')
const loading = ref(false)
const msg = ref('')
const smsCountdown = ref(0)
const rememberPwd = ref(false)
let smsTimer: number | null = null

const phoneValid = computed(() => /^\d{5,15}$/.test(phone.value))

onMounted(() => {
  const saved = localStorage.getItem('ncm.phone_login')
  if (saved) {
    try {
      const data = JSON.parse(saved)
      phone.value = data.phone || ''
      password.value = data.password || ''
      rememberPwd.value = !!data.password
    } catch {}
  }
})

function saveLogin() {
  if (rememberPwd.value && phone.value && password.value) {
    localStorage.setItem('ncm.phone_login', JSON.stringify({
      phone: phone.value,
      password: password.value
    }))
  } else {
    localStorage.setItem('ncm.phone_login', JSON.stringify({
      phone: phone.value
    }))
  }
}

async function sendSms() {
  if (!phoneValid.value) {
    msg.value = '请输入正确的手机号'
    return
  }
  try {
    loading.value = true
    msg.value = '发送中...'
    await SendSmsCode(phone.value, 86)
    msg.value = '验证码已发送'
    smsCountdown.value = 60
    smsTimer = window.setInterval(() => {
      smsCountdown.value--
      if (smsCountdown.value <= 0 && smsTimer) {
        clearInterval(smsTimer)
        smsTimer = null
      }
    }, 1000)
  } catch (e: any) {
    msg.value = '发送失败: ' + e
  } finally {
    loading.value = false
  }
}

async function loginWithSms() {
  if (!phoneValid.value || !captcha.value) {
    msg.value = '请填写手机号和验证码'
    return
  }
  try {
    loading.value = true
    msg.value = '登录中...'
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
  if (!phoneValid.value || !password.value) {
    msg.value = '请填写手机号和密码'
    return
  }
  try {
    loading.value = true
    msg.value = '登录中...'
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
  <NSpace vertical :size="16">
    <NInput v-model:value="phone" placeholder="手机号" :maxlength="15" />

    <NTabs type="segment">
      <NTabPane name="sms" tab="验证码登录">
        <NSpace vertical :size="12" style="margin-top: 12px">
          <NInputGroup>
            <NInput v-model:value="captcha" placeholder="短信验证码" :maxlength="6" />
            <NButton @click="sendSms" :disabled="smsCountdown > 0" :loading="loading">
              {{ smsCountdown > 0 ? `${smsCountdown}s` : '获取验证码' }}
            </NButton>
          </NInputGroup>
          <NButton type="primary" block @click="loginWithSms" :loading="loading">
            登录
          </NButton>
        </NSpace>
      </NTabPane>

      <NTabPane name="password" tab="密码登录">
        <NSpace vertical :size="12" style="margin-top: 12px">
          <NInput
            v-model:value="password"
            type="password"
            showPasswordOn="click"
            placeholder="密码"
          />
          <NCheckbox v-model:checked="rememberPwd">记住密码</NCheckbox>
          <NButton type="primary" block @click="loginWithPassword" :loading="loading">
            登录
          </NButton>
        </NSpace>
      </NTabPane>
    </NTabs>

    <NText v-if="msg" :type="msg.includes('失败') ? 'error' : 'success'">
      {{ msg }}
    </NText>
  </NSpace>
</template>
