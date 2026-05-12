<script setup lang="ts">
import SearchBox from './components/SearchBox.vue'
import SearchResults from './components/SearchResults.vue'
</script>

<template>
  <div class="window">
    <SearchBox />
    <SearchResults />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
html, body {
  width: 100%;
  height: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}
#app {
  width: 100%;
  height: 100%;
}
.window {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1b2636;
  color: #e0e0e0;
  user-select: none;
}
</style>
