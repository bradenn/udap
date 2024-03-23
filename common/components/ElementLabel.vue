<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, reactive} from "vue"
import {useRouter} from "vue-router";
import Element from "./Element.vue";

const props = defineProps<{
  title?: string
  icon?: string
  to?: string
}>()

const state = reactive({})

const router = useRouter();


</script>

<template>
  <Element :foreground="true"
           class="d-flex align-items-center gap-2 ent" style="height: 3.25rem">

    <div class="d-flex flex-column gap-0 justify-content-center-center px-2">
      <div></div>
      <div class="label-c4 lh-1 label-o5 d-flex">
        <div v-if="props.icon" class="sf-icon label-c6" style="width: 15px; margin-right: 12px">
          {{ props.icon }}
        </div>
        <div>{{ props.title }}</div>
      </div>
      <!--      <div class="label-c5 lh-1 label-o3">sds</div>-->
    </div>
    <div class="flex-grow-1 d-flex justify-content-end px-2">
      <slot></slot>
    </div>
  </Element>
</template>

<style lang="scss" scoped>
.sf-icon {
  color: hsl(0, 0%, 50%);
}

</style>