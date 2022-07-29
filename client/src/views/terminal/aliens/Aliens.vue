<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
onerror = function (e) {
  console.log(e)
}
</script>

<template>
  <div class="h-100">
    <router-view></router-view>
  </div>
</template>

<style scoped>

</style>
