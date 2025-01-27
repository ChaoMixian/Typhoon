<template>
    <div>
      <h1>配置</h1>
      <v-form>
        <v-text-field
          v-model="settings.server"
          label="服务器地址"
        ></v-text-field>
        <v-checkbox
          v-model="settings.autoConnect"
          label="自动连接"
        ></v-checkbox>
        <v-btn color="primary" @click="saveSettings">保存配置</v-btn>
      </v-form>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Settings',
    data() {
      return {
        settings: {
          server: 'http://localhost',
          autoConnect: true
        }
      }
    },
    methods: {
      saveSettings() {
        console.log('配置已保存：', this.settings)
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
  