<script setup>
import { onMounted } from 'vue'
import Login from './views/Login.vue'
import { useUserStore } from './stores/user'

const userStore = useUserStore()

onMounted(() => {
  userStore.checkLogin()
})
</script>

<template>
  <Login />
</template>

<style>
:root {
  --bg: #f5f5f7;
  --surface: #ffffff;
  --text-primary: #1d1d1f;
  --text-secondary: #86868b;
  --text-tertiary: #aeaeb2;
  --accent: #007aff;
  --accent-hover: #0066d6;
  --red: #ff3b30;
  --border: rgba(0, 0, 0, 0.08);
  --divider: rgba(0, 0, 0, 0.06);
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.04);
  --shadow-md: 0 2px 12px rgba(0, 0, 0, 0.06);
  --radius: 10px;
  --radius-sm: 6px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Helvetica Neue', sans-serif;
  background: var(--bg);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

input, textarea, button {
  font-family: inherit;
}
</style>
