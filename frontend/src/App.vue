<template>
  <v-app class="rounded rounded-md">

    <!-- 左侧导航栏 -->
    <v-navigation-drawer app expand-on-hover rail permanent>
      <v-divider></v-divider>
      <v-list density="compact" nav>

        <v-list-item :to="'/'" link prepend-icon="mdi-view-dashboard" title="总览" value="overview"></v-list-item>

        <v-list-item :to="'/proxy'" link prepend-icon="mdi-network" title="代理" value="proxy"></v-list-item>

        <v-list-item :to="'/rules'" link prepend-icon="mdi-routes" title="规则" value="rules"></v-list-item>

        <v-list-item :to="'/settings'" link prepend-icon="mdi-cogs" title="配置" value="settings"></v-list-item>

        <v-list-item :to="'/logs'" link prepend-icon="mdi-file-document" title="日志" value="logs"></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 顶部导航栏 -->
    <v-app-bar app color="white" dark>
      <v-toolbar-title>{{ pageTitle }}</v-toolbar-title>
    </v-app-bar>

    <!-- 主内容区域 -->
    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>

  </v-app>
</template>

<script>

export default {
  name: 'App',
  components: {
  },
  data() {
    return {
      pageTitle: '', // 当前页面标题
    };
  },
  watch: {
    // 监听路由变化，动态更新标题
    $route(to) {
      this.updateTitle(to);
    },
  },
  created() {
    // 初始化标题
    this.updateTitle(this.$route);
  },
  methods: {
    updateTitle(route) {
      this.pageTitle = route.meta.title || 'Typhoon';
    },
  },
}
</script>

<style>
html,
body,
#app {
  height: 100%;
  margin: 0;
  padding: 0;
  flex-direction: row;
  /* 横向布局 */
  background-color: #fff; /* 设置全局背景色 */

}


v-app {
  display: flex;
  flex-direction: column;
  height: 100%;
}

v-toolbar-title {
  background-color: #fff;
  color: #000;
  font: 32px;
}

v-main {
  flex-grow: 1;
  overflow: auto;
  width: 100%;
  height: 100%;
}
</style>
