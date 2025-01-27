<template>
    <div>
      <h1>日志</h1>
      <v-card>
        <v-card-text>
          <v-list>
            <v-list-item
              v-for="(log, index) in logs"
              :key="index"
            >
              <v-list-item-content>
                <v-list-item-title>{{ log.title }}</v-list-item-title>
                <v-list-item-subtitle>{{ log.timestamp }}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Logs',
    data() {
      return {
        logs: [
          { title: '启动应用', timestamp: '2025-01-25 14:30:00' },
          { title: '测试代理连接', timestamp: '2025-01-25 14:35:00' }
        ]
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
  