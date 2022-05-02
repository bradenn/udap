<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject} from "vue";

let remote = inject('remote')
let preferences = inject('preferences')


</script>

<template>
  <div>
    {{ remote }}
  </div>
</template>

<style scoped>


</style>