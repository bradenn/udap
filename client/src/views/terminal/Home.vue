<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {defineAsyncComponent, inject, onMounted, reactive, watchEffect} from "vue";
import router from "@/router";
import App from "@/components/App.vue";
import Macros from "@/components/widgets/Macros.vue";
import type {Attribute, Entity, Remote} from "@/types";
import Error from "@/components/Error.vue";
import Camera from "@/components/widgets/Camera.vue";

const Spotify = defineAsyncComponent({
    loader: () => import('@/components/widgets/Spotify.vue'),
    errorComponent: Error,
    timeout: 250
})
const Calendar = defineAsyncComponent({
    loader: () => import('@/components/widgets/Calendar.vue'),
    errorComponent: Error,
    timeout: 250
})

const Weather = defineAsyncComponent({
    loader: () => import('@/components/widgets/Weather.vue'),
    errorComponent: Error,
    timeout: 250
})


// Define the local reactive data for this view
let state = reactive<{
    lights: any[]
    hideHome: boolean
    apps: any[]
    shortcuts: any[]
    atlas: string,
    page?: string
}>({
    lights: [],
    hideHome: false,
    apps: [],
    atlas: "",
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
let remote: Remote = inject('remote') as Remote


// Force the light state to be read on load
onMounted(() => {
    updateLights(remote.entities)
    getRoutes()
    state.hideHome = false
})

// Update the Lights based on the remote injection changes
watchEffect(() => updateLights(remote.entities))

watchEffect(() => updateAtlas(remote.attributes))

function updateAtlas(attributes: Attribute[]) {
    let target = attributes.find(a => a.key === "buffer")
    if (!target) return;

    state.atlas = target.value
}

// Update the current set of lights based on the entities provided
function updateLights(entities: Entity[]) {
    // Find all applicable entities
    let candidates = entities.filter((f: Entity) => f.type === 'spectrum' || f.type === 'switch' || f.type === 'dimmer');
    candidates = candidates.filter((e: Entity) => remote.attributes.filter((a: Attribute) => a.entity === e.id).length >= 1)
    // Sort and assign them to the reactive object
    state.lights = candidates.sort(compareName)
    return entities
}

function getRoutes() {
    let routes = router.getRoutes()
    let bingo = routes.find(r => r.path === '/terminal')
    if (!bingo) return
    state.apps = bingo.children
}

</script>

<template>

    <div :class="``" class="h-100 w-100">
        <div class="home-grid w-100 h-100">

            <div style="grid-column: 1/span 4; grid-row: 1/span 4;">
                <Macros></Macros>
            </div>

            <div class="d-flex justify-content-center flex-column align-items-center gap-1"
                 style=" grid-column: 6/span 5; grid-row: 1/span 6;">
                <Calendar></Calendar>
                <Camera></Camera>
            </div>

            <div
                    style="grid-column: 1/span 11; grid-row: 10 / span 4; ">
                <div class="label-c1 label-w600 label-o4 px-1 mb-1"
                     style="font-family: 'SF Pro Display', sans-serif">
                    <span class="label-o2">􀎪</span>
                    Recent Macros
                </div>
                <div class="mc-grid">

                    <div v-for="macro in remote.macros.sort((a, b) => new Date(a.updated).valueOf() >= new Date(b.updated).valueOf() ? -1 : 1).slice(0, 10)"
                         class="element d-flex justify-content-between align-items-center py-0">
                        <div style="padding-left: 0.125rem">
                            <div class="label-c2 label-o3 label-w500 lh-1">
                                {{ macro.name }}
                            </div>
                            <div class="label-c3 label-o2 label-w500 lh-1">
                                {{ macro.description.slice(0, 18) }}
                            </div>
                        </div>
                        <div class="subplot d-flex align-items-center justify-content-center"
                             style="height: 1.5rem; width: 1.6rem">
                            <div class="label-c3 d-block label-o3 label-w700">􀊃</div>
                        </div>
                    </div>

                </div>
            </div>
            <div style="grid-column: 12 / span 4; grid-row: 1 / span 8;">
                <div class="pb-0 d-flex flex-column gap-2">
                    <Spotify></Spotify>
                    <Weather></Weather>
                    <div>
                        <div class="app-grid">
                            <App v-for="i in state.apps" :key="i.name" :icon="i.icon || 'fa-square'"
                                 :img="i?.meta?.icon"
                                 :name="i.name" :status="i?.meta?.status"
                                 style="grid-column: span 1"
                                 @click="router.push(i.path)"></App>
                        </div>
                    </div>
                </div>
            </div>

        </div>

    </div>
</template>

<style lang="scss">
@import "../../assets/sass/element.scss";

$widget-gap: 0.5rem;
$widget-aspect: 1.1618;
$widget-aspect-nested: 1.1618;
$widget-rows: 2;
$widget-cols: 4;

$widget-2x-width: 14rem;

$macro-width: 12rem;
$macro-height: 2rem;
.terminal-home > * {
  //outline: 1px solid white;
}

.mc-grid {
  display: grid;
  grid-gap: 0.25rem;

  grid-template-rows: repeat(2, minmax(2rem, 1fr));
  grid-template-columns: repeat(5, minmax(6rem, 1fr));
}

.app-grid {

  display: grid;
  padding: 0.75rem 0.5rem 0.5rem;
  flex-grow: 1;
  grid-gap: 0.75rem;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.home-grid {
  padding-top: 0.25rem;
  display: grid;
  flex-grow: 1;
  grid-gap: 0.25rem;
  grid-template-rows: repeat(12, 1fr);
  grid-template-columns: repeat(15, 1fr);
}

.home-grid > div {
  //outline: 1px solid white;
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

  display: grid;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;

  grid-row-gap: 0.5rem;
  margin-top: 0.375rem;
  grid-column-gap: 0;
  grid-template-rows: repeat(3, minmax(3.25rem, 1fr));
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

.widget-grid-vertical {

  display: grid;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;

  grid-template-rows: repeat(9, minmax(3rem, 1fr));
  grid-template-columns: repeat(1, minmax(13rem, 1fr));
  gap: $widget-gap;

}

.widget-grid {
  display: grid;
  justify-content: flex-start;
  align-items: flex-start;
  align-content: flex-start;

  grid-template-rows: repeat(8, minmax(96px, 1fr));
  grid-template-columns: repeat(16, 1fr);
  gap: $widget-gap;

}

.widget {
  animation: animateIn 20ms ease;
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

.show-outline > * > .widget {
  border-radius: 0.5rem;
  outline: 3px dashed rgba(255, 255, 255, 0.125);
}

.terminal-center {
}
</style>
