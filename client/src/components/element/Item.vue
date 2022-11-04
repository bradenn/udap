<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {reactive} from "vue";

interface Props {
  title?: string
  alt?: string
  description?: string
  selected?: boolean
  fn?: () => void
  toggle?: boolean
}

const props = defineProps<Props>()

const state = reactive({
  selected: false,
})


function touch(e: MouseEvent) {

}

</script>

<template>
  <div class="item subplot">
    <div class="d-flex flex-row justify-content-between w-100 align-items-center">
      <div class="d-flex flex-column lh-1">
        <div class="label-c2 label-o4">{{ props.title }}</div>
        <div class="label-c3 label-o3">{{ props.description }}</div>
      </div>
      <div v-if="props.toggle" class=" px-1">
        <div v-if="props.selected"
             class="label-c2 label-o2 label-w500 item-selection text-accent" style="opacity: 0.6;">􀷙</div>
        <div v-else
             class="label-c2 label-o2 label-w400 item-selection">􀓞</div>
      </div>
      <div v-else>
        <div
            class="label-c2 label-o2 label-w500 text-accent" style="opacity: 0.6;"> {{ props.alt }}
        </div>
      </div>
    </div>

  </div>
</template>


<style lang="scss" scoped>
.item:active {
  animation: accent-animation 10ms forwards cubic-bezier(.35, 1.34, .11, 1.06);

  .item-selection {
    animation: twist-animation 100ms forwards ease-in-out;

  }
}


@keyframes accent-animation {
  0% {
    transform: scale3d(1, 1, 1);
  }
  100% {
    transform: scale3d(0.99, 0.99, 0.99);
  }
}

@keyframes twist-animation {
  0% {
    rotate: 0deg;
  }
  100% {
    rotate: 60deg;
  }
}

.item {

}

</style>