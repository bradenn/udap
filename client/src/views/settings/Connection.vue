<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Header from "@/components/Header.vue";
import {inject} from "vue";

let remote = inject('remote')
</script>

<template>
  <div>
    <Header class="px-2 pb-1" icon="cloud" name="Connection"></Header>
    <div class="px-2">
      <div>
        <h4>Nexus</h4>
      </div>
      <div class="mt-2">
        <h4>Connection</h4>
        {{ remote }}
      </div>
    </div>
  </div>
</template>

<style scoped>


</style>