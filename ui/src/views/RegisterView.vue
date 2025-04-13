<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="120" sm="80" md="60" lg="400">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark>
            <v-toolbar-title>注册</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form @submit.prevent="handleRegister">
              <v-text-field
                v-model="username"
                label="用户名"
                prepend-icon="mdi-account"
                :rules="[v => !!v || '用户名必填']"
                required
              />
              <v-text-field
                v-model="email"
                label="邮箱"
                prepend-icon="mdi-email"
                type="email"
                :rules="[
                  v => !!v || '邮箱必填',
                  v => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址'
                ]"
                required
              />
              <v-text-field
                v-model="password"
                label="密码"
                prepend-icon="mdi-lock"
                type="password"
                :rules="[v => !!v || '密码必填']"
                required
              />
              <v-text-field
                v-model="confirmPassword"
                label="确认密码"
                prepend-icon="mdi-lock"
                type="password"
                :rules="[
                  v => !!v || '请确认密码',
                  v => v === password || '两次输入的密码不匹配'
                ]"
                required
              />
              <v-card-actions>
                <v-spacer />
                <v-btn color="primary" type="submit" :loading="loading" :disabled="!isFormValid">
                  注册
                </v-btn>
              </v-card-actions>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" to="/login">
              已有账号？立即登录
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

const isFormValid = computed(() => {
  return username.value &&
         email.value &&
         password.value &&
         confirmPassword.value &&
         password.value === confirmPassword.value &&
         /.+@.+\..+/.test(email.value)
})

const handleRegister = async () => {
  if (!isFormValid.value) {
    notificationStore.showError('请检查表单填写是否正确')
    return
  }

  loading.value = true
  try {
    const success = await authStore.register({
      username: username.value,
      email: email.value,
      password: password.value
    })
    if (success) {
      notificationStore.showSuccess('注册成功，请登录')
      router.push('/login')
    } else {
      notificationStore.showError('注册失败，用户名或邮箱可能已被使用')
    }
  } catch (error) {
    console.error('Registration error:', error)
    notificationStore.showError('注册失败，请稍后重试')
  }
  loading.value = false
}
</script>
