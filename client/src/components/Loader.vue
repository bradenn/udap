<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
interface Loader {
  size?: string
}

let props = defineProps<Loader>()
</script>

<template>

  <span class="loader" style="height: 0.65rem"></span>
</template>
s
<style lang="scss" scoped>
.loader {
  width: 4px;
  height: 0.45rem;
  display: flex;
  //margin: auto;
  gap: 0;
  left: -20px;
  position: relative;
  border-radius: 4px;
  box-sizing: border-box;
  animation: animloader 500ms cubic-bezier(0, 1, 0, 1) infinite alternate;
}

$clarity: rgba(255, 255, 255, 0.24);
$climax: rgba(255, 149, 0, 0.8);

@keyframes animloader {
  0% {
    box-shadow: 20px 0 $climax, 26px 0 $clarity, 32px 0 $clarity, 38px 0 $clarity;
  }
  33% {
    box-shadow: 20px 0 $clarity, 26px 0 $climax, 32px 0 $clarity, 38px 0 $clarity;
  }
  66% {
    box-shadow: 20px 0 $clarity, 26px 0 $clarity, 32px 0 $climax, 38px 0 $clarity;
  }
  100% {
    box-shadow: 20px 0 $clarity, 26px 0 $clarity, 32px 0 $clarity, 38px 0 $climax;
  }
}

</style>