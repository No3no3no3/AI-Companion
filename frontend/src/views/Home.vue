<template>
  <div class="home">
    <div class="hero">
      <h1>AI-VTuber</h1>
      <p>Intelligent Virtual Companion</p>
      <el-button type="primary" size="large" @click="startChat">
        Start Chatting
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const startChat = () => {
  router.push('/chat')
}
</script>

<style scoped lang="scss">
.home {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-hover));
  color: white;

  .hero {
    text-align: center;

    h1 {
      font-size: 3rem;
      margin-bottom: 1rem;
      font-weight: 700;
    }

    p {
      font-size: 1.2rem;
      margin-bottom: 2rem;
      opacity: 0.9;
    }
  }
}
</style>