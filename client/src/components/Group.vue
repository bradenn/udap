<script>

export default {
  name: "Group",
  data() {
    return {}
  },
  props: {
    name: String,
    icon: String,
    subtitle: String,
    small: Boolean,
    medium: Boolean,
    large: Boolean
  }
}

</script>

<template>
  <div>
    <div v-if="name" class="h-bar justify-content-between group surface mb-a">
      <div class="h-bar align-items-center">
        <div class="label-md label-o6 label-w500 mx-1 text-capitalize">{{ icon ? `${icon}&nbsp;` : '' }}{{ name }}</div>
        <div class="label-sm label-o2 label-w300 mx-2">{{ subtitle }}</div>
      </div>
      <div class="label-sm label-o2 label-w400 mx-2">􀫺</div>
    </div>
    <div :class="`${small?'group-small':''} ${medium?'group-medium':''} ${large?'group-large':''} gap`">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.group-small {

  display: grid;
  max-width: calc(100%);
  grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));
  grid-auto-flow: row;
  flex-wrap: wrap;


}

.group-medium {

  display: grid;
  max-width: calc(100%);
  grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
  grid-auto-flow: row;
  flex-wrap: wrap;
  gap: 0.25rem;

}

.group-large {

  display: grid;
  max-width: calc(100%);
  grid-template-columns: repeat(auto-fill, minmax(100%, 1fr));
  grid-auto-flow: row;
  flex-wrap: wrap;
  gap: 0.25rem;

}
</style>
