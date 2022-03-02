<script>

export default {
  components: {},
  data() {
    return {}
  },
  watch: {},
  props: {
    items: Object
  },
  created() {

  },
  methods: {},
}


</script>


<template>
  <div class="history">
    <div v-for="item in items" :style="`height:calc(${item.level/100} * 0.6rem);`" class="history-event label-ys"></div>
  </div>
</template>

<style lang="scss" scoped>
.history {
  display: flex;
  flex-direction: row;
  justify-content: start;
  align-items: center;
  gap: 0.0625rem;
}

.history-event {
  flex-shrink: 0;
  flex-grow: 0;
  background-color: rgba(255, 255, 255, 0.35);
  border-radius: 1px;
  width: 0.125rem;
  animation: historyIn 500ms forwards ease;
}

@keyframes historyIn {
  0% {
    transform: scaleY(0);
    opacity: 0;
  }
  100% {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>
