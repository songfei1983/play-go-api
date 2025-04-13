import { defineStore } from 'pinia'

interface NotificationState {
  message: string
  color: string
  show: boolean
  timeout: number
}

export const useNotificationStore = defineStore('notification', {
  state: (): NotificationState => ({
    message: '',
    color: 'success',
    show: false,
    timeout: 3000
  }),

  actions: {
    showSuccess(message: string) {
      this.message = message
      this.color = 'success'
      this.show = true
    },

    showError(message: string) {
      this.message = message
      this.color = 'error'
      this.show = true
    },

    hide() {
      this.show = false
    }
  }
})
