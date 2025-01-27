<template>
    <div>
      <h1>规则</h1>
      <v-data-table
        :headers="headers"
        :items="rules"
        class="elevation-1"
      >
        <template #top>
          <v-toolbar flat>
            <v-toolbar-title>规则列表</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="addRule">添加规则</v-btn>
          </v-toolbar>
        </template>
      </v-data-table>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Rules',
    data() {
      return {
        headers: [
          { text: '名称', value: 'name' },
          { text: '描述', value: 'description' },
          { text: '操作', value: 'actions', sortable: false }
        ],
        rules: [
          { name: '规则1', description: '描述1' },
          { name: '规则2', description: '描述2' }
        ]
      }
    },
    methods: {
      addRule() {
        this.rules.push({ name: `规则${this.rules.length + 1}`, description: '新规则描述' })
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
  