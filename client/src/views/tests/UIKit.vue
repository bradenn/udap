<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Task} from "@/types";
import {TaskType} from "@/types";
import {onMounted, reactive} from "vue";
import Element from "@/components/element/Element.vue";
import Item from "@/components/element/Item.vue";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import ToolbarButton from "@/components/ToolbarButton.vue";


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
  <h1 class="px-2">UDAP::UI</h1>
  <div class="px-2">
    <div>
      <div class="label-md label-w700">Solids</div>
      <div class="d-flex flex-column gap-1">
        <div class="label-xs label-w600">Element</div>

        <div class="demo-grid">
          <Element :cols="2" :rows="2" style="grid-column: span 2" sublink="/dd" subtext="Link" title="4x4 Grid">
            <Item v-for="idx in Array(4).keys()" :description="`This is item number ${idx + 1}`"
                  :title="`Item ${idx + 1}`"></Item>
          </Element>
          <Element :cols="2" :rows="2" style="grid-column: span 2">
            <Item v-for="idx in Array(4).keys()" :selected="idx === 0" description="This is a penis" title="Penis"
                  toggle></Item>
          </Element>
          <Element :cols="4" :rows="4">
            <div v-for="idx in Array(16).keys()" :class="`${idx === 0?'':'subplot-inline'}`"
                 class="subplot flex-grow-1 justify-content-center">
              <div :class="`${idx === 0?' text-accent ':''}`" class="label-c2 lh-1 label-w800 label-r">
                {{ idx }}
              </div>
            </div>
          </Element>

        </div>

        <div>
          <Toolbar icon="􀥁" title="Title">
            <ToolbarButton text="Option 1"></ToolbarButton>
            <ToolbarButton text="Option 2"></ToolbarButton>
            <div class="flex-grow-1"></div>
            <ToolbarButton text="Option 3"></ToolbarButton>
          </Toolbar>
        </div>
        <div>
          <div class="label-c1 label-w500">Subplot</div>
        </div>
      </div>

    </div>
    <div>
      <h2 class="">Plots</h2>
      <div>
        <h3 class="">Opaque</h3>

      </div>
    </div>
  </div>


</template>

<style scoped>

.demo-grid {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

</style>