<template>
  <div class="settings-page">
    <div class="settings-container">
      <h1>Settings</h1>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="General" name="general">
          <div class="settings-section">
            <h3>Appearance</h3>
            <el-form>
              <el-form-item label="Theme">
                <el-select v-model="theme">
                  <el-option label="Light" value="light" />
                  <el-option label="Dark" value="dark" />
                  <el-option label="Auto" value="auto" />
                </el-select>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="Voice" name="voice">
          <div class="settings-section">
            <h3>Voice Settings</h3>
            <el-form>
              <el-form-item label="Enable Voice Input">
                <el-switch v-model="voiceInputEnabled" />
              </el-form-item>
              <el-form-item label="Enable Voice Output">
                <el-switch v-model="voiceOutputEnabled" />
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="AI Model" name="model">
          <div class="settings-section">
            <h3>Model Preferences</h3>
            <el-form>
              <el-form-item label="Default Model">
                <el-select v-model="defaultModel">
                  <el-option label="GPT-3.5 Turbo" value="gpt-3.5-turbo" />
                  <el-option label="GPT-4" value="gpt-4" />
                </el-select>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref('general')
const theme = ref('light')
const voiceInputEnabled = ref(true)
const voiceOutputEnabled = ref(true)
const defaultModel = ref('gpt-3.5-turbo')
</script>

<style scoped lang="scss">
.settings-page {
  min-height: 100vh;
  background: var(--bg-secondary);
  padding: var(--spacing-xl);

  .settings-container {
    max-width: 800px;
    margin: 0 auto;
    background: var(--bg-primary);
    border-radius: var(--radius-lg);
    padding: var(--spacing-xl);
    box-shadow: var(--shadow-md);

    h1 {
      margin-bottom: var(--spacing-lg);
      color: var(--text-primary);
    }

    .settings-section {
      margin-top: var(--spacing-lg);

      h3 {
        margin-bottom: var(--spacing-md);
        color: var(--text-primary);
      }
    }
  }
}
</style>