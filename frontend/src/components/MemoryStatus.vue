<template>
  <v-card
    class="mt-8 mx-auto overflow-visible"
    max-width="90%"
    min-height="100px"
  >
    <v-sheet
      class="v-sheet--offset mx-auto"
      color=#7B68EE
      elevation="30"
      max-width="95%"
      rounded="lg"
    >
      <v-sparkline
        :labels=memoryLabels
        :model-value=memoryData
        color="white"
        line-width="2"
        padding="16"
        min="0"
        max="120"
        smooth
        auto-draw
      ></v-sparkline>

    </v-sheet>

    <v-card-text class="pt-0">
      <div class="text-h6 font-weight-light mb-2">
        <!-- <v-icon>mdi-memory</v-icon> -->
        内存占用情况
      </div>
      <div v-if="isRunning" class="subheading font-weight-light text-grey">Mihomo 正在运行</div>
      <div v-else class="subheading font-weight-light text-grey">Mihomo 未运行</div>
      <v-divider class="my-2"></v-divider>
      <v-icon
        class="me-2"
        size="small"
      >
        mdi-memory
      </v-icon>
      <span class="text-caption text-grey font-weight-light">Current memory: {{currentMemory}} MB</span>
    </v-card-text>
  </v-card>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from "vue";
// import { fetchMemoryStatus } from "@/services/memory";

export default {
  name: "MemoryStatus",
  setup() {
    const currentMemory = ref(0);
    const memoryData = ref([]);
    const memoryLabels = ref([]);
    const socket = ref(null);
    const isRunning = ref(false);

    // memoryData.value = [10, 20, 30, 40, 50];


    const connectWebSocket = () => {
      // 建立 WebSocket 连接
      socket.value = new WebSocket("ws://localhost:9097/memory");

      socket.value.onopen = () => {
        isRunning.value = true;
        console.log("WebSocket 连接已建立");
      };

      socket.value.onmessage = (event) => {
        try {
          isRunning.value = true;
          const data = JSON.parse(event.data); // 解析接收到的数据
          const inuseMB = (data.inuse / (1024 * 1024)).toFixed(2); // 转换为 MB

          // 更新当前内存
          currentMemory.value = inuseMB;
          console.log("当前内存使用量:", inuseMB, "MB");
          console.log("长度:", memoryData.value.length);
          // 更新曲线图数据
          if (memoryData.value.length >= 6) {
            memoryData.value.shift(); // 移除最旧的数据
            memoryLabels.value.shift();
          }
          memoryData.value.push(inuseMB);
          memoryLabels.value.push(new Date().toLocaleTimeString());
        } catch (error) {
          currentMemory.value = 0;
          console.error("WebSocket 数据解析失败:", error);
        }
      };

      socket.value.onclose = () => {
        console.log("WebSocket 连接已关闭");
      };

      socket.value.onerror = (error) => {
        currentMemory.value = 0;
        isRunning.value = false;
        console.error("WebSocket 错误:", error);
      };
    };

    onMounted(() => {
      connectWebSocket();
    });

    onBeforeUnmount(() => {
      if (socket.value) {
        socket.value.close(); // 关闭 WebSocket 连接
      }
    });

    return {
      currentMemory,
      memoryData,
      memoryLabels,
      isRunning,
    };
  },
};
</script>

<style scoped>
.v-sheet--offset {
  top: -10px;
  position: relative;
}

.v-card {
  max-width: 800px;
  width: 90%; /* 支持宽度自适应 */
}

.v-sheet {
  height: auto;
  min-height: 100px; /* 确保内容区足够大 */
}

.v-sparkline {
  height: 80%; /* 确保图表高度更清晰 */
}

.d-flex {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
