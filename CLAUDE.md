# Go + Vue3 项目代码生成规范和系统提示词

## 项目开发指导原则

### 1. 整体架构规范

#### 1.1 技术栈约定
- **后端**: Go 1.21+ + Gin + GORM + PostgreSQL + Redis
- **前端**: Vue 3 + Vite + Pinia + Vue Router + TypeScript
- **通信**: HTTP + WebSocket (实时交互)
- **AI 集成**: OpenAI API + Azure Speech + Whisper

#### 1.2 项目结构约定
```
ai-vtuber/
├── backend/          # Go 后端服务
│   ├── cmd/          # 应用入口
│   ├── internal/     # 内部代码
│   ├── pkg/          # 公共包
│   └── configs/      # 配置文件
├── frontend/         # Vue3 前端应用
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── views/       # 页面
│   │   ├── stores/      # 状态管理
│   │   ├── services/    # API 服务
│   │   └── utils/       # 工具函数
│   └── public/
└── docs/             # 项目文档
```

### 2. Go 后端代码规范

#### 2.1 Go 代码风格
- 遵循 Go 官方代码规范和 `gofmt` 格式化
- 使用 `golangci-lint` 进行代码质量检查
- 包名使用小写字母，简短且有意义
- 接口名以 -er 结尾或描述行为
- 结构体字段使用驼峰命名，JSON 标签使用下划线

#### 2.2 错误处理规范
```go
// 定义明确的错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// 错误包装和上下文
if err != nil {
    return fmt.Errorf("failed to process user: %w", err)
}
```

#### 2.3 依赖注入模式
```go
// 使用构造函数注入依赖
type ChatService struct {
    repo   repository.ChatRepository
    llm    infrastructure.LLMClient
    logger logger.Logger
}

func NewChatService(repo repository.ChatRepository, llm infrastructure.LLMClient, logger logger.Logger) *ChatService {
    return &ChatService{
        repo:   repo,
        llm:    llm,
        logger: logger,
    }
}
```

#### 2.4 接口设计规范
```go
// 定义清晰的接口，实现依赖倒置
type ChatRepository interface {
    Save(ctx context.Context, chat *domain.Chat) error
    FindByID(ctx context.Context, id string) (*domain.Chat, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Chat, error)
}

// 实现 LLM 客户端接口
type LLMClient interface {
    GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    GenerateStream(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)
    GetModels(ctx context.Context) ([]Model, error)
}
```

#### 2.5 WebSocket 处理规范
```go
// WebSocket Hub 管理所有连接
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

#### 2.6 数据库操作规范
```go
// 使用事务处理复杂操作
func (r *chatRepository) SaveWithMessages(ctx context.Context, chat *domain.Chat, messages []*domain.Message) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(chat).Error; err != nil {
            return err
        }

        for _, msg := range messages {
            msg.ChatID = chat.ID
            if err := tx.Create(msg).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

// 使用 GORM 进行数据库操作
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}
```

#### 2.7 中间件开发规范
```go
// JWT 认证中间件
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}

// 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Second), 10)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. Vue3 前端代码规范

#### 3.1 Vue3 组件开发规范
```vue
<!-- 使用 Composition API + TypeScript -->
<template>
  <div class="chat-window">
    <div class="message-list" ref="messageListRef">
      <MessageItem
        v-for="message in messages"
        :key="message.id"
        :message="message"
        :is-own="message.userId === currentUserId"
      />
    </div>

    <div class="input-area">
      <MessageInput
        v-model="inputMessage"
        :loading="isLoading"
        @send="handleSendMessage"
        @voice-start="handleVoiceStart"
        @voice-end="handleVoiceEnd"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'
import { WebSocketService } from '@/services/websocket'
import { AudioProcessor } from '@/utils/audio'
import type { Message, User } from '@/types'

// Props 定义
interface Props {
  chatId: string
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '100%'
})

// Emits 定义
interface Emits {
  (e: 'message-sent', message: Message): void
  (e: 'voice-toggle', isActive: boolean): void
}

const emit = defineEmits<Emits>()

// Stores
const chatStore = useChatStore()
const userStore = useUserStore()

// Refs
const messageListRef = ref<HTMLElement>()
const inputMessage = ref('')
const isLoading = ref(false)
const isRecording = ref(false)

// Services
const wsService = new WebSocketService()
const audioProcessor = new AudioProcessor()

// Computed
const messages = computed(() => chatStore.getCurrentMessages(props.chatId))
const currentUserId = computed(() => userStore.currentUser?.id)

// Methods
const handleSendMessage = async () => {
  if (!inputMessage.value.trim() || isLoading.value) return

  isLoading.value = true
  try {
    const message = await chatStore.sendMessage({
      content: inputMessage.value,
      chatId: props.chatId,
      type: 'text'
    })

    inputMessage.value = ''
    emit('message-sent', message)
    await scrollToBottom()
  } catch (error) {
    console.error('Failed to send message:', error)
  } finally {
    isLoading.value = false
  }
}

const handleVoiceStart = async () => {
  try {
    isRecording.value = true
    await audioProcessor.startRecording()
    emit('voice-toggle', true)
  } catch (error) {
    console.error('Failed to start recording:', error)
  }
}

const handleVoiceEnd = async () => {
  try {
    isRecording.value = false
    const audioBlob = await audioProcessor.stopRecording()
    const transcription = await chatService.transcribeAudio(audioBlob)

    inputMessage.value = transcription
    emit('voice-toggle', false)
  } catch (error) {
    console.error('Failed to process voice:', error)
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

// Lifecycle
onMounted(() => {
  wsService.connect()
  wsService.onMessage((message) => {
    chatStore.addMessage(message)
    scrollToBottom()
  })
})
</script>

<style scoped lang="scss">
.chat-window {
  display: flex;
  flex-direction: column;
  height: v-bind(height);
  background: var(--bg-primary);
  border-radius: 12px;
  overflow: hidden;

  .message-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    scroll-behavior: smooth;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: var(--bg-secondary);
    }

    &::-webkit-scrollbar-thumb {
      background: var(--border-color);
      border-radius: 3px;
    }
  }

  .input-area {
    border-top: 1px solid var(--border-color);
    padding: 16px;
    background: var(--bg-secondary);
  }
}
</style>
```

#### 3.2 Pinia Store 设计规范
```typescript
// stores/chat.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { chatService } from '@/services/chat'
import type { Chat, Message, ChatState } from '@/types'

export const useChatStore = defineStore('chat', () => {
  // State
  const chats = ref<Chat[]>([])
  const currentChatId = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const currentChat = computed(() =>
    chats.value.find(chat => chat.id === currentChatId.value)
  )

  const currentMessages = computed(() =>
    currentChat.value?.messages || []
  )

  // Actions
  const loadChats = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await chatService.getChats()
      chats.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load chats'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const sendMessage = async (messageData: Partial<Message>) => {
    if (!currentChatId.value) {
      throw new Error('No active chat')
    }

    // 乐观更新
    const tempMessage: Message = {
      id: `temp-${Date.now()}`,
      content: messageData.content || '',
      type: messageData.type || 'text',
      userId: 'current-user',
      chatId: currentChatId.value,
      timestamp: new Date(),
      status: 'sending'
    }

    addMessage(tempMessage)

    try {
      const response = await chatService.sendMessage({
        ...messageData,
        chatId: currentChatId.value
      })

      // 替换临时消息
      updateMessage(tempMessage.id, response.data)
      return response.data
    } catch (err) {
      // 更新消息状态为失败
      updateMessage(tempMessage.id, { status: 'failed' })
      throw err
    }
  }

  const addMessage = (message: Message) => {
    const chat = chats.value.find(c => c.id === message.chatId)
    if (chat) {
      chat.messages.push(message)
      chat.updatedAt = new Date()
    }
  }

  const updateMessage = (messageId: string, updates: Partial<Message>) => {
    const chat = chats.value.find(c =>
      c.messages.some(m => m.id === messageId)
    )

    if (chat) {
      const messageIndex = chat.messages.findIndex(m => m.id === messageId)
      if (messageIndex !== -1) {
        chat.messages[messageIndex] = {
          ...chat.messages[messageIndex],
          ...updates
        }
      }
    }
  }

  const setCurrentChat = (chatId: string) => {
    currentChatId.value = chatId
  }

  const createChat = async (title: string) => {
    try {
      const response = await chatService.createChat({ title })
      chats.value.unshift(response.data)
      setCurrentChat(response.data.id)
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create chat'
      throw err
    }
  }

  return {
    // State
    chats,
    currentChatId,
    isLoading,
    error,

    // Computed
    currentChat,
    currentMessages,

    // Actions
    loadChats,
    sendMessage,
    addMessage,
    updateMessage,
    setCurrentChat,
    createChat
  }
})
```

#### 3.3 API 服务设计规范
```typescript
// services/api.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useUserStore } from '@/stores/user'

class ApiService {
  private instance: AxiosInstance

  constructor() {
    this.instance = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        const userStore = useUserStore()
        const token = userStore.token

        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }

        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      (error) => {
        if (error.response?.status === 401) {
          const userStore = useUserStore()
          userStore.logout()
          window.location.href = '/login'
        }

        return Promise.reject(error)
      }
    )
  }

  // HTTP 方法封装
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.get(url, config)
    return response.data
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.post(url, data, config)
    return response.data
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.put(url, data, config)
    return response.data
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.delete(url, config)
    return response.data
  }

  // 文件上传
  async upload<T>(url: string, file: File, onProgress?: (progress: number) => void): Promise<T> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await this.instance.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      }
    })

    return response.data
  }
}

export const apiService = new ApiService()

// services/chat.ts
import { apiService } from './api'
import type { Chat, Message, SendMessageRequest } from '@/types'

export const chatService = {
  async getChats() {
    return apiService.get<Chat[]>('/chats')
  },

  async getChat(chatId: string) {
    return apiService.get<Chat>(`/chats/${chatId}`)
  },

  async createChat(data: { title: string; model?: string }) {
    return apiService.post<Chat>('/chats', data)
  },

  async sendMessage(data: SendMessageRequest) {
    return apiService.post<Message>('/chats/messages', data)
  },

  async transcribeAudio(audioBlob: Blob) {
    const formData = new FormData()
    formData.append('audio', audioBlob, 'recording.webm')

    return apiService.post<{ transcription: string }>('/speech/transcribe', formData)
  },

  async synthesizeSpeech(text: string, options?: { voice?: string; speed?: number }) {
    return apiService.post<{ audioUrl: string }>('/speech/synthesize', {
      text,
      ...options
    })
  }
}
```

#### 3.4 WebSocket 服务规范
```typescript
// services/websocket.ts
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'

export interface WebSocketMessage {
  type: string
  data: any
  timestamp: string
}

export class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectInterval = 1000
  private isConnecting = false
  private messageHandlers: Map<string, Function[]> = new Map()

  constructor(url?: string) {
    this.url = url || `${import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080'}/ws`
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      if (this.isConnecting) {
        reject(new Error('Connection already in progress'))
        return
      }

      this.isConnecting = true
      const userStore = useUserStore()
      const token = userStore.token

      this.ws = new WebSocket(`${this.url}?token=${token}`)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.isConnecting = false
        this.reconnectAttempts = 0
        resolve()
      }

      this.ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          this.handleMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.isConnecting = false
        this.handleReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        this.isConnecting = false
        reject(error)
      }
    })
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  send(type: string, data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      const message: WebSocketMessage = {
        type,
        data,
        timestamp: new Date().toISOString()
      }
      this.ws.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket is not connected')
    }
  }

  onMessage(type: string, handler: Function) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, [])
    }
    this.messageHandlers.get(type)!.push(handler)
  }

  offMessage(type: string, handler: Function) {
    const handlers = this.messageHandlers.get(type)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  private handleMessage(message: WebSocketMessage) {
    const handlers = this.messageHandlers.get(message.type)
    if (handlers) {
      handlers.forEach(handler => handler(message.data))
    }

    // 特殊处理某些消息类型
    switch (message.type) {
      case 'chat_message':
        this.handleChatMessage(message.data)
        break
      case 'voice_stream':
        this.handleVoiceStream(message.data)
        break
      case 'avatar_expression':
        this.handleAvatarExpression(message.data)
        break
    }
  }

  private handleChatMessage(data: any) {
    const chatStore = useChatStore()
    chatStore.addMessage(data)
  }

  private handleVoiceStream(data: any) {
    // 处理语音流数据
  }

  private handleAvatarExpression(data: any) {
    // 处理虚拟形象表情
  }

  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      setTimeout(() => {
        this.reconnectAttempts++
        console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
        this.connect().catch(error => {
          console.error('Reconnect failed:', error)
        })
      }, this.reconnectInterval * this.reconnectAttempts)
    } else {
      console.error('Max reconnect attempts reached')
    }
  }
}
```

#### 3.5 类型定义规范
```typescript
// types/index.ts
export interface User {
  id: string
  username: string
  email: string
  avatar?: string
  settings: UserSettings
  createdAt: Date
  updatedAt: Date
}

export interface UserSettings {
  theme: 'light' | 'dark' | 'auto'
  language: string
  voiceSettings: VoiceSettings
  modelPreferences: ModelPreferences
  notifications: NotificationSettings
}

export interface VoiceSettings {
  inputEnabled: boolean
  outputEnabled: boolean
  inputDevice?: string
  outputDevice?: string
  voiceId: string
  speed: number
  pitch: number
}

export interface ModelPreferences {
  defaultModel: string
  temperature: number
  maxTokens: number
  systemPrompt: string
}

export interface Chat {
  id: string
  title: string
  userId: string
  model: string
  messages: Message[]
  createdAt: Date
  updatedAt: Date
  settings: ChatSettings
}

export interface Message {
  id: string
  content: string
  type: 'text' | 'image' | 'audio' | 'system'
  userId: string
  chatId: string
  timestamp: Date
  status: 'sending' | 'sent' | 'delivered' | 'failed'
  metadata?: MessageMetadata
}

export interface MessageMetadata {
  tokens?: number
  model?: string
  latency?: number
  audioUrl?: string
  imageUrl?: string
}

export interface ChatSettings {
  model: string
  temperature: number
  maxTokens: number
  systemPrompt: string
  voiceEnabled: boolean
  avatarEnabled: boolean
}

export interface LLMModel {
  id: string
  name: string
  provider: string
  description: string
  maxTokens: number
  supportsStreaming: boolean
  pricing: {
    input: number
    output: number
  }
}

export interface Voice {
  id: string
  name: string
  language: string
  gender: 'male' | 'female' | 'neutral'
  provider: string
  sampleUrl?: string
}

export interface AvatarModel {
  id: string
  name: string
  type: 'live2d' | '3d'
  thumbnailUrl: string
  modelUrl: string
  expressions: Expression[]
}

export interface Expression {
  id: string
  name: string
  trigger: string[]
  duration: number
  intensity: number
}

// API 请求/响应类型
export interface SendMessageRequest {
  content: string
  type: 'text' | 'voice'
  chatId: string
  metadata?: {
    voiceId?: string
    stream?: boolean
  }
}

export interface ApiResponse<T> {
  data: T
  message: string
  success: boolean
  timestamp: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}
```

### 4. 开发工作流程和代码生成提示词

#### 4.1 Go 后端功能开发提示词模板

```
请为 AI-VTuber 项目的 Go 后端实现 [功能名称] 功能。

具体要求：
1. 遵循 Clean Architecture 架构模式
2. 在 internal/service/ 中实现业务逻辑
3. 在 internal/api/handlers/ 中实现 HTTP 处理器
4. 在 internal/repository/ 中实现数据访问层
5. 在 internal/domain/ 中定义领域模型
6. 添加必要的错误处理和日志记录
7. 包含单元测试和集成测试
8. 更新 API 文档和路由配置

技术要求：
- 使用 Gin 框架实现 RESTful API
- 使用 GORM 进行数据库操作
- 实现适当的错误处理和状态码
- 添加输入验证和中间件
- 使用结构化日志记录
- 实现缓存策略（如适用）

请生成完整的代码实现，包括：
1. 领域模型定义
2. 接口设计
3. 服务层实现
4. 处理器实现
5. 仓储层实现
6. 测试用例
```

#### 4.2 Vue3 前端功能开发提示词模板

```
请为 AI-VTuber 项目的 Vue3 前端实现 [功能名称] 功能。

具体要求：
1. 使用 Vue 3 Composition API 和 TypeScript
2. 创建可复用的组件
3. 使用 Pinia 进行状态管理
4. 实现响应式设计和移动端适配
5. 添加加载状态和错误处理
6. 实现动画和过渡效果
7. 遵循无障碍设计原则

技术要求：
- 使用 TypeScript 严格模式
- 组件使用 <script setup> 语法
- 实现适当的 Props 和 Emits 类型定义
- 使用 SCSS 编写样式，支持主题切换
- 实现组件的懒加载和代码分割
- 添加必要的错误边界处理

请生成完整的代码实现，包括：
1. Vue 组件文件（.vue）
2. TypeScript 类型定义
3. Pinia store 模块
4. API 服务函数
5. 工具函数和常量
6. 样式文件
7. 单元测试
```

#### 4.3 实时功能开发提示词模板

```
请为 AI-VTuber 项目实现实时 [功能名称] 功能，需要 WebSocket 支持。

后端要求：
1. 实现 WebSocket 连接管理器
2. 处理客户端连接和断开
3. 实现消息广播和路由
4. 添加连接认证和授权
5. 实现心跳检测和重连机制
6. 处理高并发和性能优化

前端要求：
1. 实现 WebSocket 客户端封装
2. 处理连接状态和错误恢复
3. 实现消息队列和重试机制
4. 添加连接状态指示器
5. 实现数据的实时同步
6. 处理网络异常和重连

技术栈：
- 后端：Go + Gorilla WebSocket
- 前端：Vue 3 + 原生 WebSocket API
- 通信：JSON 格式消息
- 认证：JWT Token

请生成完整的实时通信实现，包括服务端和客户端代码。
```

### 5. 代码质量检查清单

#### 5.1 Go 后端检查项
- [ ] 代码符合 Go 官方规范
- [ ] 错误处理完整且一致
- [ ] 接口设计清晰合理
- [ ] 数据库操作使用事务
- [ ] 实现适当的缓存策略
- [ ] 添加必要的日志记录
- [ ] 单元测试覆盖率 > 80%
- [ ] API 文档完整

#### 5.2 Vue3 前端检查项
- [ ] 使用 TypeScript 严格模式
- [ ] 组件 Props 和 Emits 类型完整
- [ ] 响应式设计适配
- [ ] 无障碍设计符合标准
- [ ] 性能优化（懒加载、虚拟滚动）
- [ ] 错误边界处理
- [ ] 单元测试覆盖关键功能
- [ ] 代码分割和包大小优化

### 6. 性能优化指导

#### 6.1 后端优化
- 使用连接池管理数据库和 Redis 连接
- 实现 Goroutine 池处理并发任务
- 使用缓存减少数据库查询
- 实现流式处理减少内存占用
- 添加适当的限流和熔断机制

#### 6.2 前端优化
- 实现路由级别的代码分割
- 使用虚拟滚动处理大列表
- 实现图片和模型的懒加载
- 使用 Web Workers 处理重计算
- 优化打包体积和缓存策略

这些规范和指导原则将确保 AI-VTuber 项目的高质量实现，Go 后端提供稳定高效的服务，Vue3 前端提供现代化的用户体验。