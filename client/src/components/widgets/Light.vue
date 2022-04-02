<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watch} from "vue";
import AttributeComponent from "@/components/entity/Attribute.vue"
import type {Attribute, Entity} from "@/types"

// Establish a local reactive state
let state = reactive<{
  active: boolean,
  showMenu: boolean,
  activeColor: string,
  powerAttribute: Attribute,
  levelAttribute: Attribute,
  attributes: Attribute[],
}>({
  active: false,
  showMenu: false,
  activeColor: "rgba(255,255,255,1)",
  levelAttribute: {} as Attribute,
  powerAttribute: {} as Attribute,
  attributes: []
})

// Define the prop for this entity
let props = defineProps<{
  entity: Entity
}>()

// Inject the remote manifest
let remote: any = inject('remote')
let hideHome: any = inject('hideHome')

// When the view loads, force the local state to update
onMounted(() => {
  updateLight(remote.attributes)
})

// Ripple any attribute changes down to the local reactive state
watch(remote.attributes, (newAttributes: Attribute[]) => {
  updateLight(newAttributes)
})

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
}

// Toggle the state of the context menu
function toggleMenu(): void {
  state.showMenu = !state.showMenu
  hideHome(state.showMenu)
}
</script>

<template>
  <div>
    <div :class="state.active?'active':''" class="entity-small element" @click="toggleMenu">
      <div class="entity-header mb-2 ">
        <div class="label-o5">
          {{ props.entity.icon || '􀛮' }}
        </div>
        <div class="label-c1 label-w400 label-o4 px-2">
          {{ props.entity.name }}
        </div>
        <div class="fill"></div>

      </div>
      <div v-if="state.powerAttribute" class="d-flex justify-content-between w-100">
        <div class="label-xxs label-o2 label-w500"
             v-text="(state.active) ? `ON${(state.levelAttribute) ? ', ' + state.levelAttribute.value + '%' : ''}` : 'OFF'"></div>
        <div class="label-c3 float-end align-self-center label-o2 label-w500">
          ~{{ Math.round(parseInt(state.levelAttribute.value, 10) / 100.0 * 7 * 100) / 100 }}w
        </div>
      </div>
    </div>

    <div v-if="state.showMenu" class="" @click="toggleMenu">
      <div class="entity-context">
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
          <div class="context-container gap v-bar">
            <div
                v-for="attribute in state.attributes.filter((attribute: Attribute) => attribute.key !== 'on')">
              <AttributeComponent :key="attribute.id" :attribute="attribute" :entity-id="props.entity.id" primitive
                                  small></AttributeComponent>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.entity-context {
  position: absolute;
  top: 0;
  width: 100%;
  padding: 1rem;
  height: calc(100% - 4.5rem);
}


</style>