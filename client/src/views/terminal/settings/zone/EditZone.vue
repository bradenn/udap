<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Entity, Task, TaskOption, Zone} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import zoneService from "@/services/zoneService";
import type {Remote} from "@/remote";
import Button from "@/components/Button.vue";

const props = defineProps<{
  zone: Zone
}>()

const router = useRouter()

const remote = inject("remote") as Remote

onMounted(() => {
  setOptions()
})

const state = reactive({
  tasks: [] as Task[],
  loaded: false,
})

watchEffect(() => {
  setOptions()

})

function setOptions() {
  if (!remote) return;

  let entities: TaskOption[] = props.zone.entities.map(t => {
    return {title: t.alias ? t.alias : t.name, description: t.module, value: t.id}
  }) as TaskOption[]

  state.tasks = [
    {
      title: "Name",
      description: "The name of the Macro",
      type: TaskType.String,
      value: props.zone.name,
      preview: props.zone.name
    },
    {
      title: "Entities",
      description: "A brief description of the macro's function.",
      type: TaskType.List,
      options: entities,
      value: props.zone.entities.map((e: Entity) => `${e.id}`),
      preview: `${props.zone.entities.length} items`
    },

  ]
  state.loaded = true
}


function goBack() {
  router.push("/terminal/settings/zones")
}

function finish(tasks: Task[]) {
  const name = tasks.find(t => t.title === "Name");
  if (!name) return;

  const entities = tasks.find(t => t.title === "Entities");
  if (!entities) return;

  zoneService.updateZone({
    name: name.value as string,
    entities: remote.entities.filter(e => (entities.value as string[]).includes(e.id)) as Entity[],
  } as Zone).then(res => {
    console.log(res)
    goBack()
  }).catch(err => {
    console.log(err)
  })

}


function createSubroutine() {


}

</script>

<template>
  <div class="ctx " @click="() => {}">
    <div class="context-grid">
      <div class="context-pane  d-flex flex-column" style="" @click.stop>
        <div class="nav-grid gap-1 pb-1 w-100">
          <div class="d-flex justify-content-start">

            <Button accent active fill icon="􀆉" text="Back"
                    @click="() => {}"></Button>
          </div>
          <div class="d-flex justify-content-center">
            <div class="label-w500 label-c1 label-w600 align-self-center">Create Macro</div>
          </div>
          <div class="d-flex justify-content-end gap-1 align-items-center">
            <!--            <div class="v-sep" style="height: 1rem"></div>-->

          </div>
        </div>
        <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Macro`"></TaskManager>


      </div>
    </div>
  </div>
</template>

<style scoped>
.context-pane {
  grid-column: 3 / span 6;
  grid-row: 2 / span 5;
}

.nav-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.generic-grid > div {
  width: 18rem;
}

.generic-grid {
  display: flex;
  justify-content: center;
  grid-column-gap: 0.25rem;
  width: max(60%, 0%);
}

.context-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(10, 1fr);
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
</style>