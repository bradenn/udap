<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import Element from "./Element.vue";
import {onMounted, reactive} from "vue";

export interface Item {
  key: string,
  value: string,
}

const props = defineProps<{
  change?: (n: string) => void,
  value: string,
  options: Item[]
}>()

const state = reactive({
  selectedKey: "",
  selectedValue: "",
  toggled: false,
})

onMounted(() => {
  state.selectedKey = props.options.find(f => f.value == props.value)?.key
  state.selectedValue = props.value
})

function selectItem(i: Item) {
  state.selectedKey = i.key
  state.selectedValue = i.value
  state.toggled = false
  if (props.change) {
    props.change(state.selectedValue)
  }
}

</script>

<template>
  <div class="element-container" @touchstart.stop>
    <Element v-if="!state.toggled" foreground @click="() => state.toggled = !state.toggled">
      <div>
        {{ state.selectedKey }}
      </div>

    </Element>
    <div v-if="state.toggled">
      <div class="d-flex justify-content-between gap-1">
        <Element v-for="elem in props.options" :key="elem.value" class="w-100 text-center"
                 foreground
                 style="border-radius: 4px !important;"
                 @click="() => selectItem(elem)">
          {{ elem.key }}
        </Element>
      </div>
    </div>
  </div>
</template>


<style lang="scss" scoped>

.ud-button {
  padding: 8px 16px;

  background-image: linear-gradient(to bottom, rgba(41, 43, 48, 1) 0%, rgba(41, 43, 48, 0.9) 100%);
  text-align: left;
  border-radius: 8px;
}

.ud-card {
  position: absolute;
  background-color: rgba(19, 21, 25, 1);
  background-image: linear-gradient(to bottom, rgba(41, 43, 48, 1) 0%, rgba(41, 43, 48, 1) 100%);
  padding: 8px;
  width: 100%;
  border-radius: 8px;
  flex-direction: column;
  display: flex;
  min-height: 48px;
  min-width: 80px;
  z-index: 10000 !important;
}

.element-container {
  position: relative;
  z-index: 1000000 !important;
}

.element-dropdown {
  position: absolute;

}
</style>