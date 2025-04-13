import { defineStore } from 'pinia'
import axios from 'axios'

interface User {
  id: number
  username: string
  email: string
}

interface AuthState {
  token: string | null
  user: User | null
}

interface LoginCredentials {
  username: string
  password: string
}

interface RegisterData extends LoginCredentials {
  email: string
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: localStorage.getItem('token'),
    user: null
  }),

  getters: {
    isAuthenticated(): boolean {
      return !!this.token
    }
  },

  actions: {
    async login({ username, password }: LoginCredentials) {
      try {
        const response = await axios.post('/api/v1/login', { username, password })
        this.token = response.data.token
        localStorage.setItem('token', this.token!)
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
        return true
      } catch (error) {
        console.error('Login failed:', error)
        return false
      }
    },

    async register(data: RegisterData) {
      try {
        await axios.post('/api/v1/register', data)
        return true
      } catch (error) {
        console.error('Registration failed:', error)
        return false
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
    }
  }
})
