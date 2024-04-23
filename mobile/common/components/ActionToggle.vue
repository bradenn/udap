<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "./Element.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Action} from "../types";
import core from "../core";
import actionService from "../services/actionService";

const props = defineProps<{
  id: string
}>()

const state = reactive({
  action: {} as Action
})

const remote = core.remote()
const router = core.router()

onMounted(() => {
  updateData(remote.actions)
  return remote.actions
})

watchEffect(() => {
  updateData(remote.actions)
  return remote.actions
})

function updateData(actions: Action[]) {
  let action = actions.find(a => a.id == props.id)
  if (action) {
    state.action = action
  }
  return state.action
}

function execute() {
  actionService.trigger(state.action.id).then(res => {

  }).catch(err => {

  })
}

function settings() {
  router.push(`/apps/actions/${state.action.id}`)
}

</script>

<template>
  <Element :cb="() => execute()" :long-cb="() => settings()" class="subplot " mutable style="padding: 0.8rem 1.125rem">
    <div class="d-flex gap-2 justify-content-between align-items-center py-1 gap-3 w-100">
      <div class="d-flex justify-content-center align-items-center gap-3">
        <div :class="`icon-on`" class="sf-icon">􁓽</div>
        <div>
          <div class="label-c5 label-w600 name" style="overflow-x: clip; word-break: normal; word-wrap: break-word">
            {{ state.action.name }}
          </div>
        </div>
      </div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>

</style>