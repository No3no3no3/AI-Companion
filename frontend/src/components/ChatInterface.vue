<!--
  聊天界面主组件
  这是整个聊天应用的核心组件，包含三个主要区域：
  1. 左侧1/3：聊天消息展示区域
  2. 右侧2/3：背景图片展示区域
  3. 底部：聊天输入框区域

  后端开发者注意事项：
  - 这里使用了Vue3的Composition API语法 (<script setup>)
  - 消息数据存储在messages数组中，每个消息包含id、content、type、timestamp等字段
  - sendMessage方法负责处理消息发送，后端需要提供对应的API接口
  - 组件使用了响应式布局，适配不同屏幕尺寸
-->
<template>
  <div class="chat-container">
    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 左侧聊天消息展示区域 - 占1/3宽度 -->
      <div class="chat-messages-area">
        <div class="messages-header">
          <h3>聊天记录</h3>
        </div>

        <!-- 消息列表容器 -->
        <div class="messages-list" ref="messagesListRef">
          <!-- 遍历显示所有消息 -->
          <div
            v-for="message in messages"
            :key="message.id"
            class="message-item"
            :class="{ 'message-user': message.type === 'user', 'message-ai': message.type === 'ai' }"
          >
            <div class="message-content">
              {{ message.content }}
            </div>
            <div class="message-time">
              {{ formatTime(message.timestamp) }}
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧背景图片展示区域 - 占2/3宽度 -->
      <div class="background-area">
        <div class="background-image">
          <!-- 使用网络图片或者本地图片作为背景 -->
          <img
            src="/images/background.png"
            alt="AI Assistant Background"
            class="ai-background"
          />
          <div class="background-overlay">
            <div class="ai-info">
              <h2>AI 助手</h2>
              <p>我是您的智能助手，有什么可以帮助您的吗？</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部聊天输入框区域 -->
    <div class="chat-input-area">
      <div class="input-container">
        <!-- 消息输入框 -->
        <input
          v-model="inputMessage"
          @keydown.enter="sendMessage"
          type="text"
          class="message-input"
          placeholder="请输入您的消息..."
          :disabled="isLoading"
        />

        <!-- 发送按钮 -->
        <button
          @click="sendMessage"
          class="send-button"
          :disabled="isLoading || !inputMessage.trim()"
        >
          {{ isLoading ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'

// 响应式数据声明
const messages = ref([]) // 存储所有聊天消息
const inputMessage = ref('') // 当前输入框的内容
const isLoading = ref(false) // 加载状态，用于显示发送中状态
const messagesListRef = ref(null) // 消息列表的DOM引用，用于自动滚动到底部

// 初始化一些示例消息，方便后端开发者理解数据结构
onMounted(() => {
  messages.value = [
    {
      id: '1',
      content: '您好！我是AI助手，很高兴为您服务！',
      type: 'ai',
      timestamp: new Date(Date.now() - 60000) // 1分钟前
    }
  ]
})

/**
 * 发送消息的方法
 * 后端开发者需要在这里调用API接口发送消息到后端服务
 * 当前是模拟实现，实际使用时需要替换为真实的API调用
 */
const sendMessage = async () => {
  // 检查输入是否为空或正在加载
  if (!inputMessage.value.trim() || isLoading.value) return

  // 创建用户消息对象
  const userMessage = {
    id: Date.now().toString(), // 使用时间戳作为临时ID
    content: inputMessage.value.trim(),
    type: 'user',
    timestamp: new Date()
  }

  // 将用户消息添加到消息列表
  messages.value.push(userMessage)

  // 清空输入框并设置加载状态
  const messageContent = inputMessage.value.trim()
  inputMessage.value = ''
  isLoading.value = true

  // 滚动到底部显示新消息
  await scrollToBottom()

  try {
    // 调用真实的后端API
    const apiResponse = await callChatAPI(messageContent)
    console.log("apiResponse : ",apiResponse)
    // 检查API响应是否成功
    if (apiResponse.code === 0) {
      // 创建AI回复消息
      const reply = apiResponse.data
      const aiMessage = {
        id: reply.messageId,
        content: reply.reply,
        type: 'ai',
        timestamp: new Date()
      }

      // 添加AI回复到消息列表
      messages.value.push(aiMessage)

      // 再次滚动到底部
      await scrollToBottom()
    } else {
      // 处理API返回的错误
      throw new Error(apiResponse.error || 'Unknown error occurred')
    }

  } catch (error) {
    console.error('发送消息失败:', error)

    // 创建错误消息显示给用户
    const errorMessage = {
      id: (Date.now() + 1).toString(),
      content: '抱歉，消息发送失败，请稍后重试。错误信息：' + (error.message || '未知错误'),
      type: 'ai',
      timestamp: new Date()
    }

    // 添加错误消息到消息列表
    messages.value.push(errorMessage)
    await scrollToBottom()
  } finally {
    isLoading.value = false
  }
}

/**
 * 滚动消息列表到底部
 * 确保新消息始终可见
 */
const scrollToBottom = async () => {
  await nextTick() // 等待DOM更新完成
  if (messagesListRef.value) {
    messagesListRef.value.scrollTop = messagesListRef.value.scrollHeight
  }
}

/**
 * 格式化时间显示
 * 将时间戳转换为更友好的显示格式
 */
const formatTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time

  // 如果是今天的消息，显示时间
  if (time.toDateString() === now.toDateString()) {
    return time.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 如果是昨天的消息，显示"昨天 时间"
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (time.toDateString() === yesterday.toDateString()) {
    return '昨天 ' + time.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 其他情况显示完整日期时间
  return time.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * 调用后端API发送消息
 * 发送请求到 http://127.0.0.1:8080/api/chat
 */
const callChatAPI = async (message) => {
  try {
    const response = await fetch('http://127.0.0.1:8080/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        message: message,
        userId: 'user_001' // 可以根据实际需求修改
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    return data
  } catch (error) {
    console.error('API call failed:', error)
    throw error
  }
}
</script>

<style scoped>
/* 聊天容器整体样式 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-family: 'Arial', sans-serif;
}

/* 主要内容区域 - 包含左侧聊天和右侧背景 */
.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 左侧聊天消息区域样式 - 占1/3宽度 */
.chat-messages-area {
  width: 33.333%;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
}

/* 消息区域头部样式 */
.messages-header {
  padding: 20px;
  background: rgba(103, 126, 234, 0.1);
  border-bottom: 1px solid rgba(103, 126, 234, 0.2);
}

.messages-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
  font-weight: 600;
}

/* 消息列表容器样式 */
.messages-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 自定义滚动条样式 */
.messages-list::-webkit-scrollbar {
  width: 6px;
}

.messages-list::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.messages-list::-webkit-scrollbar-thumb {
  background: rgba(103, 126, 234, 0.5);
  border-radius: 3px;
}

.messages-list::-webkit-scrollbar-thumb:hover {
  background: rgba(103, 126, 234, 0.7);
}

/* 消息项样式 */
.message-item {
  display: flex;
  flex-direction: column;
  max-width: 80%;
  word-wrap: break-word;
}

.message-user {
  align-self: flex-end;
}

.message-ai {
  align-self: flex-start;
}

/* 消息内容样式 */
.message-content {
  padding: 12px 16px;
  border-radius: 18px;
  font-size: 14px;
  line-height: 1.4;
  white-space: pre-wrap;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.message-user .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom-right-radius: 4px;
}

.message-ai .message-content {
  background: #f3f4f6;
  color: #333;
  border-bottom-left-radius: 4px;
}

/* 消息时间样式 */
.message-time {
  font-size: 11px;
  color: #666;
  margin-top: 4px;
  padding: 0 8px;
}

.message-user .message-time {
  text-align: right;
}

.message-ai .message-time {
  text-align: left;
}

/* 右侧背景图片区域样式 - 占2/3宽度 */
.background-area {
  width: 66.667%;
  position: relative;
  overflow: hidden;
}

.background-image {
  width: 100%;
  height: 100%;
  position: relative;
}

.ai-background {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
}

/* 背景图片遮罩层样式 */
.background-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg,
    rgba(0, 0, 0, 0.3) 0%,
    rgba(0, 0, 0, 0.5) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-info {
  text-align: center;
  color: white;
  max-width: 80%;
}

.ai-info h2 {
  font-size: 36px;
  margin-bottom: 16px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.ai-info p {
  font-size: 18px;
  line-height: 1.6;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  opacity: 0.9;
}

/* 底部聊天输入区域样式 */
.chat-input-area {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.3);
  padding: 20px;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
}

.input-container {
  display: flex;
  gap: 12px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 消息输入框样式 */
.message-input {
  flex: 1;
  padding: 16px 20px;
  border: 2px solid #e5e7eb;
  border-radius: 25px;
  font-size: 16px;
  outline: none;
  transition: all 0.3s ease;
  background: white;
}

.message-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.message-input:disabled {
  background: #f9fafb;
  cursor: not-allowed;
}

/* 发送按钮样式 */
.send-button {
  padding: 16px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
  min-width: 100px;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.send-button:active:not(:disabled) {
  transform: translateY(0);
}

.send-button:disabled {
  background: #d1d5db;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 响应式设计 - 平板设备 */
@media (max-width: 1024px) {
  .chat-messages-area {
    width: 40%;
  }

  .background-area {
    width: 60%;
  }

  .ai-info h2 {
    font-size: 28px;
  }

  .ai-info p {
    font-size: 16px;
  }
}

/* 响应式设计 - 手机设备 */
@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
  }

  .chat-messages-area {
    width: 100%;
    height: 50%;
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .background-area {
    width: 100%;
    height: 50%;
  }

  .ai-info h2 {
    font-size: 24px;
  }

  .ai-info p {
    font-size: 14px;
  }

  .message-item {
    max-width: 90%;
  }

  .input-container {
    gap: 8px;
  }

  .message-input {
    padding: 12px 16px;
    font-size: 14px;
  }

  .send-button {
    padding: 12px 20px;
    font-size: 14px;
    min-width: 80px;
  }
}

/* 添加一些动画效果 */
.message-item {
  animation: messageSlideIn 0.3s ease-out;
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>