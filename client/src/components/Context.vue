<script>

export default {
  name: "Context",
  data() {
    return {
      show: false,
    }
  },
  created() {
    this.show = this.open
    if (!this.open) this.show = false
  },
  props: {
    name: String,
    open: Boolean,
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

    <div class="bar-but" @click="show = !show">
      <div class="h-bar">
        <div class="label-xs label-o2 label-w600 px-1">{{ icon }}</div>
        <div class="label-xs label-o4 label-w500">&nbsp;&nbsp;{{ name }}</div>
      </div>
      <div v-if="show" class="context h-100" v-on:click.stop>
        <slot></slot>
      </div>
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
