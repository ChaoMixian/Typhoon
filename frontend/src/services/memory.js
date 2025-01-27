import apiClient from "./api";

// 获取内存使用情况
export function fetchMemoryStatus() {
  return apiClient.get("/mihomo/memory");
}
