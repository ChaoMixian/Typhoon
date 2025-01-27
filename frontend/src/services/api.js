import axios from "axios";

// 创建 Axios 实例
const apiClient = axios.create({
  baseURL: "http://localhost:8090/api/v1", // 基础 API 路径
  timeout: 5000, // 请求超时时间
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 在发送请求之前可以添加一些逻辑（例如附加 Token）
    // config.headers.Authorization = `Bearer ${yourToken}`;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // 对响应数据进行处理
    return response.data;
  },
  (error) => {
    // 全局错误处理（可以自定义弹窗或日志）
    console.error("API 请求错误:", error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default apiClient;
