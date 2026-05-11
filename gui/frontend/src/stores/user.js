import { defineStore } from 'pinia'
import { ref } from 'vue'

import { IsLoggedIn, GetUserProfile, Logout } from '../../bindings/auth'

export const useUserStore = defineStore('user', () => {
  const profile = ref(null)
  const loggedIn = ref(false)
  const loading = ref(false)

  async function checkLogin() {
    try {
      const ok = await IsLoggedIn()
      loggedIn.value = ok
      if (ok) await fetchProfile()
    } catch {
      loggedIn.value = false
      profile.value = null
    }
  }

  async function fetchProfile() {
    try {
      loading.value = true
      const p = await GetUserProfile()
      profile.value = p
      loggedIn.value = true
    } catch {
      profile.value = null
      loggedIn.value = false
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try { await Logout() } finally {
      profile.value = null
      loggedIn.value = false
    }
  }

  function onLoginSuccess(p) {
    profile.value = p
    loggedIn.value = true
  }

  return { profile, loggedIn, loading, checkLogin, fetchProfile, logout, onLoginSuccess }
})
