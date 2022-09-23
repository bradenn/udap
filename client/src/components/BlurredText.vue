<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface TextProps {
  text: string
  size?: number
}

const props = defineProps<TextProps>()
</script>

<template>
  <div class="headline">
    <span class="visually-hidden">{{ props.text }}</span>
    <svg id="textPath2" aria-hidden="true" class="headline lockup-headline-mask visually-hidden">
      <clipPath id="textPath">
        <text dominant-baseline="hanging" dy="0" text-anchor="middle" x="50%" y="0em">{{ props.text }}
        </text>
      </clipPath>
    </svg>
  </div>
</template>


<style lang="scss" scoped>
div.headline {

  position: absolute;
  color: white;
}

.headline {
  width: 100%;
  font: 8vw "SF Pro Display", sans-serif;
  font-weight: 500;
  height: 150px;
  overflow: visible;
}

.visually-hidden {
  position: absolute;
  top: -9999px;
  left: -9999px;
}


div.headline {
  backdrop-filter: blur(5px) saturate(160%);
  overflow: visible !important;
  background: rgba(245, 245, 246, 0.2);
  clip-path: url(#textPath);
}
</style>