<script>
export default {
  components: {},
  data() {
    return {
      header: {
        icon: "􀌆",
        name: "Settings"
      }
    }
  },
  created() {
  },
  methods: {},
}
</script>

<template>
  <div class="window">
    <div class="main-sidebar top">
      <div class="context-container">
        <div class="d-flex justify-content-between align-items-end align-content-end mb-1" v-on:click.stop>
          <div class="label-xxl label-w600 label-o5">
            <span class="label-xl label-w600 label-o4">􀍟</span>
            <span class="label-xl label-w600 label-o6 px-2">Settings</span>
          </div>
        </div>
      </div>
      <div class="element">
        <div class="d-flex justify-content-start flex-column gap-1">
          <router-link v-for="page in this.$route.matched[1].children.filter(c => !c.hidden)"
                       class="dock-link d-flex justify-content-start px-2 py-2" draggable="false" v-bind:to="page.path">
            <div class="label-w500 label-o3 label-xxs">{{ page.sf }}&nbsp;</div>
            <div class="label-w400 label-o5 label-xxs px-1">{{ page.name }}</div>
          </router-link>
        </div>
      </div>
    </div>
    <div class="main-body top">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.window {
  width: 100%;
  max-width: 100%;

  display: flex;
  height: 100%;
  justify-content: center;
  align-content: center;
  gap: 1rem;
}

.main-sidebar {

  min-width: 12rem;
  max-width: 16vw;
}

.main-body {
  width: 100%;
  border-radius: 0.5rem;
}
</style>
