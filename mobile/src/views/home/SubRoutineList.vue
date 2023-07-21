<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import SubRoutineDom from "@/views/home/SubRoutine.vue";
import {SubRoutine} from "udap-ui/types";
import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";

const remote = core.remote()

const state = reactive({
  subroutines: [] as SubRoutine[]
})

onMounted(() => {
  state.subroutines = remote.subroutines
})

watchEffect(() => {
  state.subroutines = remote.subroutines
  return remote.subroutines
})
</script>

<template>
  <div v-for="sr in state.subroutines" v-if="state.subroutines" :key="sr.id">
    <SubRoutineDom :subroutine="sr"></SubRoutineDom>
  </div>
</template>

<style scoped>
.surface {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 1rem 0.25rem;
}
</style>