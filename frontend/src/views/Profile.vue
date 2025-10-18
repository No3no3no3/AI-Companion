<template>
  <div class="profile-page">
    <div class="profile-container">
      <h1>Profile</h1>
      <div class="profile-content">
        <div class="avatar-section">
          <el-avatar :size="100" :src="userAvatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <el-button type="primary" @click="changeAvatar">
            Change Avatar
          </el-button>
        </div>

        <div class="info-section">
          <el-form :model="userInfo" label-width="120px">
            <el-form-item label="Username">
              <el-input v-model="userInfo.username" />
            </el-form-item>
            <el-form-item label="Email">
              <el-input v-model="userInfo.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveProfile">Save Changes</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { User } from '@element-plus/icons-vue'

const userInfo = ref({
  username: 'User',
  email: 'user@example.com'
})

const userAvatar = ref('')

const changeAvatar = () => {
  // TODO: Implement avatar change
  console.log('Change avatar')
}

const saveProfile = () => {
  // TODO: Implement profile save
  console.log('Saving profile:', userInfo.value)
}
</script>

<style scoped lang="scss">
.profile-page {
  min-height: 100vh;
  background: var(--bg-secondary);
  padding: var(--spacing-xl);

  .profile-container {
    max-width: 600px;
    margin: 0 auto;
    background: var(--bg-primary);
    border-radius: var(--radius-lg);
    padding: var(--spacing-xl);
    box-shadow: var(--shadow-md);

    h1 {
      margin-bottom: var(--spacing-lg);
      color: var(--text-primary);
      text-align: center;
    }

    .profile-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: var(--spacing-xl);

      .avatar-section {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--spacing-md);
      }

      .info-section {
        width: 100%;
        max-width: 400px;
      }
    }
  }
}
</style>