<script>

export default {
  data() {
    return {}
  },
  props: {},
  created() {

  },
  beforeDestroy() {

  },
  methods: {},
}
</script>

<template>
  <div class="whiteboard">

  </div>
</template>

<style scoped>


.whiteboard {
  position: absolute !important;
  z-index: 100;
  top: 0 !important;
  left: 0 !important;
  width: 100vw;
  height: 100vw;
  margin-top: 0.125em;
  padding-bottom: 1.5rem;
  font-family: "SF Pro Text", sans-serif;
  font-size: 15px;
  font-weight: 400;
  background-color: rgba(200, 200, 200, 1);
  border-radius: 10px;
  box-shadow: 0 0 32px 4px rgba(0, 0, 0, 0.06);
}

@keyframes fullscreen {
  100% {

  }
}

</style>
