<template>
  <div class="app-container">
    <div class="top-decoration"></div>
    
    <div class="container mt-5">
      <div class="row justify-content-center">
        <div class="col-md-8 col-lg-6">
          <div class="card shadow animate-fade-in">
            <div class="card-header text-white py-4">
              <h2 class="mb-0">IP地址查询</h2>
            </div>
            <div class="card-body p-6">
              <!-- 查询表单 -->
              <div class="mb-6">
                <div class="input-group">
                  <input
                    type="text"
                    id="ipAddress"
                    v-model="ipAddress"
                    class="form-control form-control-lg"
                    placeholder="请输入IPv4或IPv6地址"
                    @keyup.enter="handleQuery"
                    :class="{ 'is-invalid': error }"
                  >
                  <button
                    class="btn btn-success btn-lg ms-2"
                    @click="handleQuery"
                    :disabled="isQuerying"
                  >
                    <font-awesome-icon icon="search" class="me-2"></font-awesome-icon>
                    {{ isQuerying ? '查询中...' : '查询' }}
                  </button>
                </div>
                
                <!-- 快捷按钮 -->
                <div class="mt-3 flex flex-wrap gap-2">
                  <button 
                    class="btn btn-outline-secondary btn-sm px-4"
                    @click="queryMyIp"
                    :disabled="isQuerying"
                  >
                    <font-awesome-icon icon="location-arrow" class="me-1"></font-awesome-icon>
                    查询我的IP
                  </button>
                  <template v-if="recentQueries.length">
                    <button 
                      class="btn btn-outline-secondary btn-sm px-4"
                      @click="clearHistory"
                    >
                      <font-awesome-icon icon="trash" class="me-1"></font-awesome-icon>
                      清除历史
                    </button>
                  </template>
                </div>
                
                <!-- 错误提示 -->
                <div v-if="error" class="invalid-feedback d-block mt-2 animate-fade-in">
                  <font-awesome-icon icon="exclamation-circle" class="me-1"></font-awesome-icon>
                  {{ error }}
                </div>
              </div>

              <!-- 加载状态 -->
              <div v-if="loading" class="text-center my-8">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
                <p class="mt-3 text-muted">正在查询，请稍候...</p>
              </div>
              
              <!-- 查询结果 -->
              <div v-if="result && !loading" class="animate-slide-up">
                <h3 class="text-primary mb-4">查询结果</h3>
                <div class="table-responsive">
                  <table class="table">
                    <tbody>
                      <tr>
                        <th scope="row" class="w-1/3">IP地址</th>
                        <td>{{ result.ip }}</td>
                      </tr>
                      <tr>
                        <th scope="row">国家</th>
                        <td>{{ result.country || '未知' }}</td>
                      </tr>
                      <tr>
                        <th scope="row">地区</th>
                        <td>{{ result.region || '未知' }}</td>
                      </tr>
                      <tr>
                        <th scope="row">城市</th>
                        <td>{{ result.city || '未知' }}</td>
                      </tr>
                      <template v-if="result.isp">
                        <tr>
                          <th scope="row">运营商</th>
                          <td>{{ result.isp }}</td>
                        </tr>
                      </template>
                      <template v-if="result.latitude && result.longitude">
                        <tr>
                          <th scope="row">坐标</th>
                          <td>{{ result.latitude }}, {{ result.longitude }}</td>
                        </tr>
                      </template>
                    </tbody>
                  </table>
                </div>
                
                <div class="text-center mt-4">
                  <button 
                    class="btn btn-outline-secondary"
                    @click="copyResult"
                  >
                    <font-awesome-icon icon="copy" class="me-1"></font-awesome-icon>
                    复制结果
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 历史记录 -->
          <div v-if="recentQueries.length && !loading" class="mt-5 card shadow animate-fade-in">
            <div class="card-header py-3">
              <h5 class="mb-0">查询历史</h5>
            </div>
            <div class="card-body">
              <div class="d-flex flex-wrap gap-2">
                <button 
                  v-for="(item, index) in recentQueries" 
                  :key="index"
                  class="btn btn-outline-secondary btn-sm"
                  @click="queryFromHistory(item)"
                >
                  {{ item }}
                </button>
              </div>
            </div>
          </div>
          
          <!-- 知识卡片 -->
          <div class="mt-5 card shadow animate-fade-in">
            <div class="card-body p-5 text-center">
              <h5 class="text-primary mb-2">IP查询小知识</h5>
              <p class="text-sm text-gray-600 mb-0">
                IPv4地址由32位二进制数组成，通常分为4个8位组；
                IPv6地址由128位二进制数组成，通常分为8个16位组，
                用于解决IPv4地址枯竭问题。
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 页脚 -->
    <footer class="bg-gray-50 border-t py-4 mt-8">
      <div class="container text-center text-sm text-gray-500">
        <p class="mb-0">© {{ new Date().getFullYear() }} IP查询工具 | 保护您的网络隐私</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import ipService from '@/services/ipService';

// 状态管理
const ipAddress = ref('');
const result = ref(null);
const error = ref('');
const loading = ref(false);
const isQuerying = ref(false);
const recentQueries = ref(JSON.parse(localStorage.getItem('ipQueryHistory') || '[]'));

// IP地址验证函数
const isValidIP = (ip) => {
  const trimmedIp = ip.trim();
  if (!trimmedIp) return false;

  // IPv4验证
  const ipv4Regex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/;
  if (ipv4Regex.test(trimmedIp)) {
    return trimmedIp.split('.').every(part => {
      const num = parseInt(part, 10);
      return num >= 0 && num <= 255 && segment === num.toString();
    });
  }

  // IPv6验证
  const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/;
  const ipv6CompressedRegex = /^(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?$/;
  
  return ipv6Regex.test(trimmedIp) || ipv6CompressedRegex.test(trimmedIp);
  };

  // 包含压缩的IPv6验证
  const ipv6CompressedRegex = /^([0-9a-fA-F]{1,4}:){1,7}:([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4}$/;
  return ipv6CompressedRegex.test(trimmedIp);
};

// 保存查询历史
const saveToHistory = (ip) => {
  // 去重
  const newHistory = [ip, ...recentQueries.value.filter(item => item !== ip)];
  // 限制历史记录数量
  if (newHistory.length > 10) {
    newHistory.pop();
  }
  recentQueries.value = newHistory;
  localStorage.setItem('ipQueryHistory', JSON.stringify(newHistory));
};

// 处理IP查询
const handleQuery = async () => {
  error.value = '';
  const ip = ipAddress.value.trim();
  
  if (!ip) {
    error.value = '请输入IP地址';
    return;
  }
  
  if (!isValidIP(ip)) {
    error.value = '请输入有效的IPv4或IPv6地址（例：192.168.1.1 或 2001:0db8:85a3:0000:0000:8a2e:0370:7334）';
    return;
  }
  
  try {
    loading.value = true;
    isQuerying.value = true;
    
    const data = await ipService.queryBasicIP(ip);
    if (!data) {
      error.value = '未查询到该IP的相关信息';
      result.value = null;
      return;
    }
    result.value = data;
    saveToHistory(ip);
  } catch (err) {
    error.value = `查询失败：${err.message || '未知错误'}`;
    result.value = null;
  } finally {
    loading.value = false;
    isQuerying.value = false;
  }
};

// 查询本机IP
const queryMyIp = async () => {
  try {
    loading.value = true;
    isQuerying.value = true;
    error.value = '';
    
    const data = await ipService.getMyIP();
    result.value = data;
    ipAddress.value = data.ip;
    
    // 保存到历史记录
    saveToHistory(data.ip);
  } catch (err) {
    error.value = err.message || '获取本机IP失败，请手动输入';
    result.value = null;
  } finally {
    loading.value = false;
    isQuerying.value = false;
  }
};

// 从历史记录查询
const queryFromHistory = (ip) => {
  ipAddress.value = ip;
  handleQuery();
};

// 清除历史记录
const clearHistory = () => {
  recentQueries.value = [];
  localStorage.removeItem('ipQueryHistory');
};

// 复制结果
const copyResult = () => {
  if (!result.value) return;
  
  const text = `IP地址: ${result.value.ip}\n国家: ${result.value.country || '未知'}\n地区: ${result.value.region || '未知'}\n城市: ${result.value.city || '未知'}`;
  
  navigator.clipboard.writeText(text).then(() => {
    alert('结果已复制到剪贴板');
  }).catch(err => {
    console.error('复制失败:', err);
    alert('复制失败，请手动复制');
  });
};

// 自动聚焦输入框
onMounted(() => {
  const input = document.getElementById('ipAddress');
  if (input) input.focus();
});
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background-color: #f8fafc;
  position: relative;
}

.top-decoration {
  height: 6px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6, #3b82f6);
  background-size: 200% 100%;
  animation: gradientShift 8s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.card {
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

/* 动画效果 */
.animate-fade-in {
  animation: fadeIn 0.3s ease-in-out;
}

.animate-slide-up {
  animation: slideUp 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

/* 响应式优化 */
@media (max-width: 768px) {
  .card-body {
    padding: 40px 20px !important;
  }
  
  h2 {
    font-size: 1.5rem !important;
  }
  
  .input-group {
    flex-direction: column;
    gap: 10px;
  }
}
</style>