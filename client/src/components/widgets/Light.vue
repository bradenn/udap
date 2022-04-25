<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import AttributeComponent from "@/components/entity/Attribute.vue"
import type {Attribute, Entity} from "@/types"
// Establish a local reactive state
let state = reactive<{
  loading: boolean,
  active: boolean,
  showMenu: boolean,
  activeColor: string,
  powerAttribute: Attribute,
  shortStatus: string
  levelAttribute: Attribute,
  attributes: Attribute[],
}>({
  loading: true,
  active: false,
  showMenu: false,
  activeColor: "rgba(255,255,255,1)",
  shortStatus: "",
  levelAttribute: {} as Attribute,
  powerAttribute: {} as Attribute,
  attributes: []
})


function generateState() {
  if (state.levelAttribute.value) {
    state.shortStatus = (state.active) ? `ON${(state.levelAttribute) ? ', ' + state.levelAttribute.value + '%' : ''}` : 'OFF'
  } else {
    state.shortStatus = (state.active) ? `ON` : 'OFF'
  }
}

// Define the prop for this entity
let props = defineProps<{
  entity: Entity
}>()

// Inject the remote manifest
let remote: any = inject('remote')
let context: any = inject('context')

// When the view loads, force the local state to update
onMounted(() => {
  state.loading = true
  updateLight(remote.attributes)
  generateState()
})

// Ripple any attribute changes down to the local reactive state
watchEffect(() => updateLight(remote.attributes))

// Compare the defined order to sort the lights
function compareOrder(a: any, b: any): number {
  return a.order - b.order
}

// Update the reactive model for the light
function updateLight(attributes: Attribute[]): void {
  // Define the attributes for the light
  state.attributes = attributes.filter((a: Attribute) => a.entity === props.entity.id).sort(compareOrder)
  // Get the current power state of the light
  let on = state.attributes.find((a: Attribute) => a.key === 'on')
  let dim = state.attributes.find((a: Attribute) => a.key === 'dim')
  // Assign the power state to a local attribute
  if (!on) return
  state.powerAttribute = on as Attribute
  if (dim) {
    state.levelAttribute = dim as Attribute
  }
  state.active = on.value === "true" || on.request === "true"
  generateState()
  state.loading = false
}

// Toggle the state of the context menu
function toggleMenu(): void {
  state.showMenu = !state.showMenu
  // context(state.showMenu)
}
</script>

<template>

  <div v-if="state.showMenu" class="context context-light" @click="toggleMenu"></div>
  <div v-if="state.loading" class="w-100 h-100">
    <div class="entity-small element">
      <div class="entity-header mb-2 ">
        <div class="label-o5">
          {{ props.entity.icon || '􀛮' }}
        </div>
        <div class="label-c1 label-w400 label-o4 px-2">
          {{ props.entity.name }}
        </div>
        <div class="fill"></div>

      </div>
    </div>
  </div>
  <div v-else class="w-100 h-100">
    <div v-if="!state.showMenu" :class="state.active?'active':''" class="entity-small element" @click="toggleMenu">
      <div class="entity-header mb-2 ">
        <div class="label-o5">
          {{ props.entity.icon || '􀛮' }}
        </div>
        <div class="label-c1 label-w400 label-o4 px-2">
          {{ props.entity.name }}
        </div>
        <div class="fill"></div>

      </div>
      <div v-if="state.powerAttribute" class="d-flex justify-content-between align-items-center w-100">

        <div class="label-c1 label-o2 label-w500">{{ state.shortStatus }}</div>

        <div v-if="state.powerAttribute.value ==='true'" class="label-c2 float-end label-o1 label-w600 label-mono">
          {{ Math.round((parseInt(state.levelAttribute.value, 10) || 20) / 100.0 * 7 * 100) / 100 }}W
        </div>
      </div>
    </div>
    <div v-else class="entity-small element quickedit">
      <div class="entity-header mb-1 d-flex justify-content-between align-items-center w-100">
        <div class="d-flex">
          <div class="label-o5">
            {{ props.entity.icon || '􀛮' }}
          </div>
          <div class="label-c1 label-w400 label-o4 px-2">
            {{ props.entity.name }}
          </div>
        </div>
        <AttributeComponent :attribute="state.powerAttribute" :entity-id="props.entity.id"
                            small></AttributeComponent>
        <div class="" @click="toggleMenu">
          <i class="fa-solid fa-circle-xmark label-o3 label-c1 label-w400 px-1"></i>
        </div>
      </div>
      <div class="w-100 d-flex flex-column gap">
        <div
            v-for="attribute in state.attributes.filter((attribute: Attribute) => attribute.key === 'dim' )">
          <AttributeComponent :key="attribute.id" :attribute="attribute" :entity-id="props.entity.id" primitive
                              small></AttributeComponent>
        </div>
      </div>
      <div v-if="false" class="entity-small element ">
        <div class="element surface d-flex flex-column gap py-4 px-3 pt-2">
          <div class="d-flex justify-content-start align-items-end align-content-end" v-on:click.stop>
            <div class="mt-1">
              <span :style="`text-shadow: 0 0 8px ${state.activeColor};`"
                    class="label-md label-w600 label-o3">{{ props.entity.icon }}</span>
              <span class="label-md label-w600 label-o6 px-2">{{ props.entity.name }}</span>
            </div>
            <div class="fill "></div>

            <div class="h-bar">
              <AttributeComponent :attribute="state.powerAttribute" :entity-id="props.entity.id"
                                  small></AttributeComponent>
            </div>
          </div>
          <div class="h-sep"></div>
          <div
              v-for="attribute in state.attributes.filter((attribute: Attribute) => attribute.key !== 'on')">
            <AttributeComponent :key="attribute.id" :attribute="attribute" :entity-id="props.entity.id" primitive
                                small></AttributeComponent>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.entity-small:not(.quickedit):active {
  animation: click 100ms ease forwards;
}

.entity-small {
  animation: click 100ms ease forwards;
}


@keyframes click {
  0% {
    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
  }
  30% {
    transform: scale(0.97);
  }
  100% {
    transform: scale(1);
  }
}

.entity-context {
  position: absolute;
  top: 0;
  width: 100%;
  padding: 1rem;
  height: calc(100% - 4.5rem);
}


</style>