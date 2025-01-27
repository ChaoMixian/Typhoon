<template>
    <div>
      <h1>代理</h1>
      <v-btn color="primary" @click="testConnection">测试连接</v-btn>
      <v-alert v-if="status" :type="status.success ? 'success' : 'error'" dismissible>
        {{ status.message }}
      </v-alert>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Proxy',
    data() {
      return {
        status: null // 测试连接状态
      }
    },
    methods: {
      testConnection() {
        // 模拟测试连接的逻辑
        this.status = { success: true, message: '连接成功！' }
        setTimeout(() => (this.status = null), 3000) // 3秒后清除状态
      }
    }
  }
  </script>
  
  <style scoped>
  h1 {
    font-size: 24px;
    margin-bottom: 16px;
  }
  </style>
  