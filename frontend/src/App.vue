<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Events } from '@wailsio/runtime'
import HelloWorld from './components/HelloWorld.vue'

const wailsReady = ref(false)

onMounted(() => {
  Events.On('common:WindowRuntimeReady', () => {
    wailsReady.value = true
  })
})
</script>

<template>
  <HelloWorld msg="ZTools on Wails" />
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
}
#app {
  width: 100%;
  height: 100%;
}
</style>
