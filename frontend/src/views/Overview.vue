<template>
  <v-container fluid class="overview-container">
    <v-row>
      <!-- 左侧内容 -->
      <v-col cols="12" md="6">
        <MemoryStatus />
      </v-col>

      <!-- 右侧内容 -->
      <v-col cols="12" md="6">
        <TrafficStatus />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import MemoryStatus from "@/components/MemoryStatus.vue";
import TrafficStatus from "@/components/TrafficStatus.vue";
export default {
  name: "Overview",
  components: {
    MemoryStatus,
    TrafficStatus,
  },
};
</script>

<style scoped>
/* 左侧文本块样式 */
.overview-text {
  padding: 16px;
  text-align: left;
}

.overview-container {
  background-color: #ffffff; /* 设置容器背景色 */
  min-height: 100vh; /* 确保占满屏幕高度 */
  padding: 16px;
}

h2 {
  margin: 0 0 8px 0;
  color: #3f51b5;
}

p {
  color: #666;
}
</style>
