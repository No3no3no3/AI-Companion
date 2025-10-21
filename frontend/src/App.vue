<!--
  App.vue - 应用程序的根组件

  这是整个Vue应用的入口组件，现在已经集成了我们的聊天界面。

  后端开发者注意事项：
  - 这里引入了ChatInterface组件作为应用的主要界面
  - 移除了Vue默认的欢迎内容，替换为我们的聊天应用
  - 全局样式设置确保应用占满整个屏幕
-->
<script setup>
// 导入聊天界面组件
import ChatInterface from './components/ChatInterface.vue'
</script>

<template>
  <!--
    应用的根容器，包含整个聊天界面
    id为app的div会被Vue挂载，这在main.js中定义
  -->
  <div id="app">
    <ChatInterface />
  </div>
</template>

<style>
/*
  全局样式重置和基础设置
  确保应用占满整个屏幕，没有默认的边距和内边距
*/
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

#app {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

/* 禁用文本选择，提升用户体验（可选） */
.user-select-none {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
</style>
