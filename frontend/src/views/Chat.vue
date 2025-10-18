<template>
  <div class="chat-page">
    <div class="chat-container">
      <div class="chat-header">
        <h2>AI Assistant</h2>
      </div>

      <div class="chat-content">
        <div class="messages-area">
          <!-- Messages will be displayed here -->
          <div class="empty-state">
            <p>No messages yet. Start a conversation!</p>
          </div>
        </div>

        <div class="input-area">
          <el-input
            v-model="messageInput"
            placeholder="Type your message..."
            @keyup.enter="sendMessage"
          >
            <template #append>
              <el-button @click="sendMessage">Send</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const messageInput = ref('')

const sendMessage = () => {
  if (messageInput.value.trim()) {
    // TODO: Implement message sending
    console.log('Sending message:', messageInput.value)
    messageInput.value = ''
  }
}
</script>

<style scoped lang="scss">
.chat-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);

  .chat-container {
    width: 100%;
    max-width: 800px;
    height: 80vh;
    background: var(--bg-primary);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .chat-header {
      padding: var(--spacing-lg);
      border-bottom: 1px solid var(--border-color);

      h2 {
        margin: 0;
        color: var(--text-primary);
      }
    }

    .chat-content {
      flex: 1;
      display: flex;
      flex-direction: column;

      .messages-area {
        flex: 1;
        padding: var(--spacing-lg);
        overflow-y: auto;

        .empty-state {
          display: flex;
          align-items: center;
          justify-content: center;
          height: 100%;
          color: var(--text-tertiary);
        }
      }

      .input-area {
        padding: var(--spacing-lg);
        border-top: 1px solid var(--border-color);
      }
    }
  }
}
</style>