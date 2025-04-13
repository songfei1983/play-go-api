<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const logout = () => {
  authStore.logout()
  notificationStore.showSuccess('已成功退出登录')
  router.push('/login')
}
</script>

<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <v-toolbar-title>管理系统</v-toolbar-title>
      <v-spacer></v-spacer>
      <template v-if="authStore.isAuthenticated">
        <v-btn icon @click="logout">
          <v-icon>mdi-logout</v-icon>
        </v-btn>
      </template>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>

    <v-snackbar
      v-model="notificationStore.show"
      :color="notificationStore.color"
      :timeout="notificationStore.timeout"
    >
      {{ notificationStore.message }}
      <template v-slot:actions>
        <v-btn
          variant="text"
          @click="notificationStore.hide()"
        >
          关闭
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
