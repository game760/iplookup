
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const TIMEOUT = Number(import.meta.env.VITE_API_TIMEOUT) || 5000;

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
    
    return data;
  } catch (error) {
    if (error.name === 'AbortError') {
      throw new Error(`请求超时（${TIMEOUT}ms），请重试`);
    }
    throw error;
  }
};

export default {
  // 基础IP查询
  queryBasicIP(ip) {
    return request(`/ip/query?ip=${encodeURIComponent(ip)}`);
  },
  
  // 详细IP查询
  queryDetailIP(ip) {
    return request(`/ip/detail?ip=${encodeURIComponent(ip)}`);
  },
  
   queryBasicIP(ip) {
    return request(`/ip/basic?ip=${encodeURIComponent(ip)}`)
      .then(response => {
        // 统一处理后端返回格式（假设后端返回 {code:0, data:..., message:''}）
        if (response.code !== 0) {
          throw new Error(response.message || '查询失败')
        }
        return response.data
      })
  }
}

  // 获取本机IP
  getMyIP() {
    return request('/ip/my');
  }
};