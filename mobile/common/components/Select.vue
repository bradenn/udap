<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "./Element.vue";
import List from "./List.vue";
import ElementToggle from "./ElementToggle.vue";
import {onMounted, reactive, watchEffect} from "vue";
import ElementHeader from "./ElementHeader.vue";
import Title from "./Title.vue";

export interface SelectOption {
  name: string
  icon?: string
  category: string
  description?: string
  id: string
}

export interface SelectCategory {
  name: string
  options: SelectOption[]
}

const props = defineProps<{
  title: string

  prompt?: string
  options: SelectOption[],
  selections?: (out: SelectOption[]) => void,
  single?: boolean
  many?: boolean
}>()


const state = reactive({
  selected: [] as string[],
  categories: [] as SelectCategory[],

})

function buildOptionManifest(options: SelectOption[]) {
  let categories = [] as SelectCategory[]

  for (let option of options) {

    let find = categories.find(c => c.name === option.category)
    if (!find) {
      let cat = {
        name: option.category,
        options: [] as SelectOption[]
      } as SelectCategory
      cat.options.push(option)
      categories.push(cat)
    } else {

      categories = categories.map((c: SelectCategory) => {
        if (c.name === option.category) {
          let local: SelectCategory = c
          local.options.push(option)
          return local
        }
        return c
      })

    }

  }

  state.categories = categories;
}


onMounted(() => {
  buildOptionManifest(props.options)

})

watchEffect(() => {
  buildOptionManifest(props.options)

  return props.options
})

function isSelected(option: SelectOption): boolean {
  return state.selected.filter(s => s == option.id).length > 0
}


function toggleSelect(option: SelectOption) {
  if (props.single) {
    if (isSelected(option)) {
      state.selected = []
      if (props.selections) {
        props.selections(state.selected.map(s => props.options.find(o => o.id === s) || {} as SelectOption))
      }
      return
    } else {
      state.selected = []
    }

  }
  if (isSelected(option)) {
    state.selected = state.selected.filter(s => s != option.id)
  } else {
    state.selected.push(option.id)
  }
  if (props.selections) {
    props.selections(state.selected.map(s => props.options.find(o => o.id === s) || {} as SelectOption))
  }
}
</script>

<template>
  <Element class="scrollable-fixed">
    <ElementHeader v-if=props.title :subtitle="props.prompt" :title="props.title"></ElementHeader>
    <div style="height: 1px; width: 100%; background-color: rgba(255,255,255,0.05)"></div>
    <List scroll-y>
      <div v-for="category in state.categories">
        <Title v-if="category.name !== 'none'" :title="category.name" class="pt-2"></Title>
        <List>
          <ElementToggle v-for="option in category.options" :accent="isSelected(option)"
                         :cb="() => toggleSelect(option)"
                         :icon="option.icon" :selected="isSelected(option)"
                         :title="option.name"></ElementToggle>
        </List>
      </div>

    </List>
  </Element>
</template>

<style lang="scss" scoped>

</style>