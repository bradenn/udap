<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Remote} from "@/remote";
import core from "@/core";

const remote: Remote = core.remote()

</script>

<template>
  <div class="">
    <div v-for="user in remote.users">
      {{ user.username }}
      {{ user.first }}
      {{ user.middle }}
      {{ user.last }}
      {{ user.classification }}
      {{ user.residency }}
    </div>
  </div>
</template>

<style scoped>

</style>