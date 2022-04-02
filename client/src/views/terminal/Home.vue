<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Light from "@/components/widgets/Light.vue";
import {inject, onMounted, reactive, watch} from "vue";
import Weather from "@/components/widgets/Weather.vue";
import router from "@/router";
import Earth from "@/components/widgets/Earth.vue";
import Widget from "@/components/widgets/Widget.vue";
import App from "@/components/App.vue";
import Shortcut from "@/components/Shortcut.vue";

// Define the local reactive data for this view
let state = reactive<{
  lights: any[]
  hideHome: boolean
  apps: any[]
  shortcuts: any[]
  page?: string
}>({
  lights: [],
  hideHome: false,
  apps: [],
  shortcuts: [
    {
      name: "Good night",
      icon: "fa-moon"
    }
  ],
  page: router.currentRoute.value.name as string
})

// Compare the names of the entities to sort them accordingly
function compareName(a: any, b: any): number {
  if (a.name < b.name)
    return -1;
  if (a.name > b.name)
    return 1;
  return 0;
}

// Inject the remote udap context
const remote: any = inject('remote')


// Force the light state to be read on load
onMounted(() => {
  updateLights(remote.entities)
  getRoutes()
  state.hideHome = false
})

// Update the Lights based on the remote injection changes
watch(remote.entities, (newEntities, oldEntities) => {
  updateLights(newEntities)
})

// Update the current set of lights based on the entities provided
function updateLights(entities: any) {
  // Find all applicable entities
  let candidates = entities.filter((f: any) => f.type === 'spectrum' || f.type === 'switch' || f.type === 'dimmer');
  // Sort and assign them to the reactive object
  state.lights = candidates.sort(compareName)
}

function getRoutes() {
  let routes = router.getRoutes()
  let bingo = routes.find(r => r.path === '/terminal')
  if (!bingo) return
  state.apps = bingo.children
}

let ui: any = inject("ui")

</script>

<template>

  <div :class="`${ui.outlines?'show-outline':''}`" class="widget-grid mt-1">
    <Widget :cols="1" :pad="3" :rows="1" size="md">
      <Weather></Weather>
    </Widget>

    <Widget :cols="2" :pad="3" :rows="2" size="md">
      <Light v-for="light in state.lights.slice(0, 4)" :key="light.id" :entity="light"
             class="flex-grow-1 h-100"></Light>
    </Widget>

    <Widget :cols="4" :pad="3" :rows="3" size="md">
      <App v-for="i in state.apps" :icon="i.icon || 'fa-square'" :name="i.name" @click="router.push(i.path)"></App>
    </Widget>

    <Widget :cols="2" :pad="3" :rows="4" size="md">
      <Shortcut v-for="i in state.shortcuts" :icon="i.icon || 'fa-square'" :name="i.name"></Shortcut>
    </Widget>

    <Widget :pad="3" size="md">
      <Earth></Earth>
    </Widget>
    <Widget :pad="3" size="md"></Widget>
    <Widget :pad="3" size="md"></Widget>
    <Widget :pad="3" size="md"></Widget>

  </div>


</template>

<style lang="scss">
@import "../../assets/sass/element.scss";

$widget-gap: 0.5rem;
$widget-aspect: 1.1618;
$widget-aspect-nested: 1.1618;
$widget-rows: 2;
$widget-cols: 4;

$widget-2x-width: 15rem;

.widget-sm {
  flex-grow: 1;
  aspect-ratio: $widget-aspect !important;
  //outline: 1px solid white;
}

.widget-app-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;

}

$macro-width: 6.5rem;
$macro-height: 2rem;

.widget-macro {
  width: $macro-width;
  height: $macro-height;
  display: flex;
  align-self: center;
  align-items: center;
  justify-content: start;
  padding-inline: 0.5rem;
  font-size: 0.8rem;
  gap: 0.375rem;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.4);
  /*outline: 1px solid white;*/
}

.widget-macro > i {
  text-shadow: 0 0 2px black;
  float: left;
}

.widget-macro > div {
  font-size: 0.7rem;
  font-weight: 500;
}

.widget-app {
  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-self: center;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1/1 !important;
  backdrop-filter: blur(12px);
  border-radius: 0.6rem;
  font-size: 0.9rem;
  color: rgba(200, 200, 200, 0.6);
  /*outline: 1px solid white;*/
}

.widget-app > i {
  text-shadow: 0 0 2px black;

}

.widget-app-container > span {
  font-size: 0.6rem;
  font-family: "SF Pro Display", sans-serif;
  font-weight: 400;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.55)
}


@keyframes fill {
  0% {
    position: relative;
    top: 0;
    left: 0;

  }
  100% {
    position: relative;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
  }
}

.widget-apps {

  aspect-ratio: $widget-aspect !important;

  display: grid;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;
  grid-row-gap: 0.5rem !important;
  grid-column-gap: 0.25rem !important;
  padding: 1rem;
  grid-template-rows: repeat(4, minmax(3.25rem, 1fr));
  grid-template-columns: repeat(4, minmax(3.25rem, 1fr));
}

.widget-macros {

  aspect-ratio: $widget-aspect !important;

  display: grid;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;
  grid-gap: 0;

  padding: 1rem;
  grid-template-rows: repeat(5, minmax($macro-height, 1fr));
  grid-template-columns: repeat(2, minmax($macro-width, 1fr));
}


.widget-container {

  padding: 1rem;
  display: grid;
  align-items: start;
  align-content: start;
  grid-gap: 1rem;
  grid-template-rows: repeat(2, minmax(5rem, 1fr));
  grid-template-columns: repeat(2, minmax(5rem, 1fr));
}

.widget-md {

  gap: 0.25rem !important;
  display: flex;
  justify-content: center;
  align-content: center;
  align-items: center;
  flex-wrap: wrap;
  overflow: clip;
  aspect-ratio: $widget-aspect !important;
  /*outline: 1px solid rgba(255, 64, 255, 0.5);*/
}

.widget-lg {

  overflow: clip;
  aspect-ratio: $widget-aspect !important;
  /*outline: 1px solid rgba(255, 64, 255, 0.5);*/
}

@keyframes focusIn {
  0% {
    filter: blur(24px);
  }
  100% {
    filter: blur(0px);
  }
}

.widget-grid {
  display: grid;
  align-items: center;
  align-content: center;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: repeat(auto-fill, minmax($widget-2x-width, 1fr));
  gap: $widget-gap;

}

.widget {
  animation: animateIn 42ms ease;
}

@keyframes animateIn {
  0% {
    transform: scale(0.94);
    filter: blur(4px);
    opacity: 0.4;
  }
  50% {
    transform: scale(0.99);
    filter: blur(0px);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    filter: blur(0px);
  }
}

.show-outline > div {
  border-radius: 0.5rem;
  outline: 3px dashed rgba(255, 255, 255, 0.125);
}

.terminal-center {
}
</style>
