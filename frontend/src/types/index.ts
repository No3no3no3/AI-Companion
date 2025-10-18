// Basic type definitions for the AI-VTuber frontend

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
}

export interface VoiceSettings {
  inputEnabled: boolean
  outputEnabled: boolean
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