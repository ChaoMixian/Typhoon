<template>
    <v-card
      class="mt-8 mx-auto overflow-visible"
      max-width="90%"
      min-height="100px"
    >
      <v-sheet
        class="v-sheet--offset mx-auto"
        color="#7B68EE"
        elevation="30"
        max-width="95%"
        rounded="lg"
      >
        <!-- 上行流量 -->
        <v-sparkline
          :labels="trafficLabels"
          :model-value="trafficDataUp"
          color="red"
          line-width="2"
          padding="16"
          min="0"
          max="120"
          smooth
          auto-draw
        ></v-sparkline>
        
        <!-- 下行流量 -->
        <v-sparkline
          :labels="trafficLabels"
          :model-value="trafficDataDown"
          color="green"
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
          网络使用情况
        </div>
        <div v-if="isRunning" class="subheading font-weight-light text-grey">Mihomo 正在运行</div>
        <div v-else class="subheading font-weight-light text-grey">Mihomo 未运行</div>
        <v-divider class="my-2"></v-divider>
        <v-icon
          class="me-2"
          size="small"
        >
          mdi-upload-box-outline
        </v-icon>
        <span class="text-caption text-grey font-weight-light">Current Up traffic: {{currentUp}} MB</span>
        <br/>
        <v-icon
          class="me-2"
          size="small"
        >
          mdi-download-box-outline
        </v-icon>
        <span class="text-caption text-grey font-weight-light">Current Down traffic: {{currentDown}} MB</span>
      </v-card-text>
    </v-card>
  </template>
  
  <script>
  import { ref, onMounted, onBeforeUnmount } from "vue"; // 导入 Vue 的响应式 API
  
  export default {
    name: "TrafficStatus",
    setup() {
      const currentUp = ref(0);
      const currentDown = ref(0);
      const trafficDataUp = ref([]);   // 上行流量数据
      const trafficDataDown = ref([]); // 下行流量数据
      const trafficLabels = ref([]);   // 时间标签
      const socket = ref(null);
      const isRunning = ref(false);
  
      const connectWebSocket = () => {
        socket.value = new WebSocket("ws://localhost:9097/traffic");
  
        socket.value.onopen = () => {
          isRunning.value = true;
          console.log("WebSocket 连接已建立");
        };
  
        socket.value.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            const up = (data.up / (1024 * 1024)).toFixed(2);
            const down = (data.down / (1024 * 1024)).toFixed(2);
  
            currentUp.value = up;
            currentDown.value = down;
  
            // 更新上行流量数据
            if (trafficDataUp.value.length >= 6) {
              trafficDataUp.value.shift();
            }
            trafficDataUp.value.push(up);
  
            // 更新下行流量数据
            if (trafficDataDown.value.length >= 6) {
              trafficDataDown.value.shift();
            }
            trafficDataDown.value.push(down);
  
            trafficLabels.value.push(new Date().toLocaleTimeString());
          } catch (error) {
            currentUp.value = 0;
            currentDown.value = 0;
            console.error("WebSocket 数据解析失败:", error);
          }
        };
  
        socket.value.onclose = () => {
          console.log("WebSocket 连接已关闭");
        };
  
        socket.value.onerror = (error) => {
          currentUp.value = 0;
          currentDown.value = 0;
          isRunning.value = false;
          console.error("WebSocket 错误:", error);
        };
      };
  
      onMounted(() => {
        connectWebSocket();
      });
  
      onBeforeUnmount(() => {
        if (socket.value) {
          socket.value.close();
        }
      });
  
      return {
        currentUp,
        currentDown,
        trafficDataUp,
        trafficDataDown,
        trafficLabels,
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
    width: 90%;
  }
  
  .v-sheet {
    height: auto;
    min-height: 100px;
  }
  
  .v-sparkline {
    height: 80%;
  }
  
  .d-flex {
    display: flex;
    align-items: center;
    justify-content: center;
  }
  </style>
  