<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="120" sm="100" md="80" lg="40">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark>
            <v-toolbar-title>登录</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form @submit.prevent="handleLogin">
              <v-text-field
                v-model="username"
                label="用户名"
                prepend-icon="mdi-account"
                required
              />
              <v-text-field
                v-model="password"
                label="密码"
                prepend-icon="mdi-lock"
                type="password"
                required
              />
              <v-card-actions>
                <v-spacer />
                <v-btn color="primary" type="submit" :loading="loading">
                  登录
                </v-btn>
              </v-card-actions>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" to="/register">
              还没有账号？立即注册
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

const handleLogin = async () => {
  if (!username.value || !password.value) {
    notificationStore.showError('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const success = await authStore.login({
      username: username.value,
      password: password.value
    })
    if (success) {
      notificationStore.showSuccess('登录成功')
      router.push('/user')
    } else {
      notificationStore.showError('用户名或密码错误')
    }
  } catch (error) {
    console.error('Login error:', error)
    notificationStore.showError('登录失败，请稍后重试')
  }
  loading.value = false
}
</script>
