<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Task} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import {onMounted, reactive} from "vue";


// const remote = inject("remote") as Remote


const state = reactive({
  tasks: [] as Task[],
  loaded: false,
})

// remote.triggers.map(t => {
//   return {title: t.name, description: t.description, value: t.id}
// })


onMounted(() => {

  state.tasks = [
    {
      title: "Name",
      description: "What should the subroutine be named?",
      type: TaskType.String,
      value: "",
      preview: "Subroutine 1"
    },
    {
      title: "Trigger",
      description: "What trigger should evoke this subroutine?",
      type: TaskType.Radio,
      options: [{title: "Trigger 1", description: "When your mom arrives", value: "991"}],
      value: "",
      preview: "unset"
    },
    {
      title: "type",
      description: "Which macros should be run when this subroutine is triggered?",
      type: TaskType.Radio,
      options: [
        {
          title: "Brightness",
          description: "Change the brightness of a light",
          value: "dim"
        }, {
          title: "Color",
          description: "Change the color of a light",
          value: "hue"
        }, {
          title: "Color Temperature",
          description: "Change the color of a light",
          value: "cct"
        },
        {
          title: "Power",
          description: "Change the power state of a device",
          value: "on"
        }],
      value: "",
      preview: "unset"
    }
  ]

  state.loaded = true
})


function finish(tasks: Task[]) {

}


</script>
<template>
  <div class="d-flex justify-content-center">
    <div class="mt-5" style="height: 50%">
      <TaskManager v-if="state.loaded" :on-complete="finish" :tasks="state.tasks" title="Title"></TaskManager>
    </div>
  </div>
</template>

<style scoped>

</style>