<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>


import {reactive} from "vue";

interface Props {
  title?: string
  subtext?: string
  sublink?: string
  rows?: number
  cols?: number
  scroll?: boolean
}

let props = defineProps<Props>()

let state = reactive({})

</script>


<template>
  <div>
    <div class="element p-1">
      <div v-if="props.title" class="element-header">
        <div class="element-label">{{ props.title }}</div>
        <div :class="`${props.sublink ? 'text-accent':''}`" class="element-subtext">{{
            props.subtext
          }}
        </div>
      </div>
      <div :class="`${props.cols || props.rows ? 'element-grid':'element-flex'}`">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.text-accent {
  border-bottom: 1px solid rgba(255, 149, 0, 0.3);
}

.element {

  .element-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-inline: 0.25rem;
    padding-top: 0.0125rem;
    padding-bottom: 0.2rem;

  }

  .element-grid {
    display: grid;
    grid-gap: 0.25rem;

    grid-template-columns: repeat(v-bind('props.cols'), minmax(1rem, 1fr));
    grid-template-rows: repeat(v-bind('props.rows'), minmax(1.75rem, 1fr));
  }

  .element-flex {
    display: flex;
    height: 100%;

    div {

      flex-grow: 1;
    }
  }

}

.element-label {
  font-size: 20px;
  font-weight: 600;
  line-height: 20px;
  color: rgba(255, 255, 255, 0.65);
  font-family: "SF Pro Rounded", sans-serif;
}

.element-subtext {
  font-size: 18px;
  font-weight: 500;
  line-height: 18px;
  color: rgba(255, 255, 255, 0.35);
  font-family: "SF Pro", sans-serif;
}
</style>