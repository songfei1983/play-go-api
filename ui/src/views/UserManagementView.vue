<template>
  <v-container>
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">用户管理</h1>
        <v-card>
          <v-card-text>
            <v-form @submit.prevent="updateProfile" v-if="currentUser">
              <v-text-field
                v-model="currentUser.username"
                label="用户名"
                readonly
                disabled
              />
              <v-text-field
                v-model="currentUser.email"
                label="邮箱"
                :rules="[v => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址']"
              />
              <v-text-field
                v-model="newPassword"
                label="新密码"
                type="password"
                hint="如需修改密码请输入新密码"
                persistent-hint
              />
              <v-text-field
                v-model="confirmPassword"
                label="确认新密码"
                type="password"
                :rules="[v => !newPassword || v === newPassword || '两次输入的密码不匹配']"
                :disabled="!newPassword"
              />
              <v-card-actions>
                <v-spacer />
                <v-btn
                  color="primary"
                  type="submit"
                  :loading="loading"
                  :disabled="!isFormValid"
                >
                  保存修改
                </v-btn>
              </v-card-actions>
            </v-form>
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              class="mt-4"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const loading = ref(false)
const currentUser = ref<any>(null)
const newPassword = ref('')
const confirmPassword = ref('')

const isFormValid = computed(() => {
  if (!currentUser.value) return false

  const isEmailValid = /.+@.+\..+/.test(currentUser.value.email)
  const isPasswordValid = !newPassword.value || newPassword.value === confirmPassword.value

  return isEmailValid && isPasswordValid
})

onMounted(async () => {
  try {
    const response = await axios.get('/api/v1/users/current')
    currentUser.value = response.data
  } catch (error) {
    console.error('Failed to fetch user profile:', error)
    notificationStore.showError('获取用户信息失败')
  }
})

const updateProfile = async () => {
  if (!currentUser.value || !isFormValid.value) {
    notificationStore.showError('请检查表单填写是否正确')
    return
  }

  loading.value = true
  try {
    const updateData = {
      email: currentUser.value.email,
      ...(newPassword.value ? { password: newPassword.value } : {})
    }

    await axios.put(`/api/v1/users/${currentUser.value.id}`, updateData)
    notificationStore.showSuccess('个人信息更新成功')
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (error) {
    console.error('Failed to update profile:', error)
    notificationStore.showError('更新失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>
