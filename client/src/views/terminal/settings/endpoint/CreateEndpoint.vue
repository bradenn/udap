<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint, Remote, Task, TaskOption} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import Page from "@/components/Page.vue";
import endpointService from "@/services/endpointService";
import core from "@/core";

const router = useRouter()

const remote = inject("remote") as Remote

interface Props {
  done?: () => void,
}

const props = defineProps<Props>()

onMounted(() => {
  setOptions()
})

const state = reactive({
  tasks: [] as Task[],
  loaded: false,
})

watchEffect(() => {
  setOptions()
  return remote
})

function setOptions() {
  if (!remote) return;

  let types: TaskOption[] = ["terminal", "console", "padd"].map(t => {
    return {title: t, description: t, value: t}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "What should the endpoint be named?",
      type: TaskType.String,
      value: "",
      preview: "unset"
    },
    {
      title: "Type",
      description: "What type of endpoint is this?",
      type: TaskType.Radio,
      options: types,
      value: types[0],
      preview: types[0].title
    }
  ]
  state.loaded = true
}

function goBack() {
  if (props.done) {
    props.done()
  }
}

let notifications = core.notifications()

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const type = tasks.find(t => t.title === "Type");
  if (!type) return;

  endpointService.createEndpoint({
    name: name.value,
    type: type.value,
  } as Endpoint).then(res => {
    notifications.show("Endpoint", "Endpoint created", 0, 8)
    goBack()
  }).catch(err => {
    notifications.show("Endpoint", "Failed to create endpoint.", 2, 8)
    console.log(err)
  })

}


function createSubroutine() {


}

</script>

<template>
  <div class="ctx " @click="goBack()">
    <div class="context-grid">
      <div class="d-flex justify-content-center h-100 w-100"
           style="grid-column: 4 / span 6 !important; grid-row: 2 / span 10 !important;">

        <div class="generic-grid h-100" @click.stop>
          <div>
            <Page :done="() => {goBack()}" :save="() => {finish(state.tasks)}" title="New Endpoint"></Page>

            <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Endpoint`" class="h-100"></TaskManager>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>


.context-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(12, 1fr);
  grid-template-columns: repeat(12, 1fr);
}

.ctx {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2;
  /*background-color: rgba(0, 0, 0, 0.1);*/
  backdrop-filter: blur(18px) brightness(90%);
}

.generic-grid {
  width: 100%;
}

</style>