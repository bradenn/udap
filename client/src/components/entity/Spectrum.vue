<script>

import Attribute from "./Attribute.vue";

export default {
  name: "Spectrum",
  components: {Attribute},
  data() {
    return {
      waiting: false,
      cache: {}
    }
  },
  props: {
    entity: Object,
  },
  created() {
  },
  watch: {
    'entity': function (d) {
      this.waiting = false
      this.cache = d
    }
  },
  computed: {
    liveColor: function () {

    }
  },
  methods: {}
}

</script>

<template>
  <div v-if="!entity">
    Loading!
  </div>
  <div v-else>
    {{ entity.name }}
  </div>
</template>

<style scoped>

</style>
