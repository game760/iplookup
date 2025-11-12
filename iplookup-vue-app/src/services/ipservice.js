// 定义常量
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7895/api';
const TIMEOUT = 10000; // 10秒超时

// 基础请求函数
const request = async (url) => {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), TIMEOUT);
    
    const response = await fetch(`${API_BASE_URL}${url}`, {
      signal: controller.signal,
      headers: {
        'Accept': 'application/json'
      }
    });
    
    clearTimeout(timeoutId);
    
    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.message || '请求失败');
    }
    
    if (data.code === 1) {
      throw new Error(data.message || '服务器返回错误');
    }
    
    return data.data || data;
  } catch (error) {
    if (error.name === 'AbortError') {
      throw new Error(`请求超时（${TIMEOUT}ms），请重试`);
    }
    throw error;
  }
};

// 查询IP基本信息
const queryBasicIP = async (ip) => {
  return request(`/ip/query?ip=${encodeURIComponent(ip)}`);
};

// 获取本机IP信息
const getMyIP = async () => {
  return request('/ip/myip');
};

export default {
  queryBasicIP,
  getMyIP
};