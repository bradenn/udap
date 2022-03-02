<script>

import Room from "./Room.vue";
import Hexagon from "../components/Hexagon.vue";
import Pane from "../components/Pane.vue";

export default {
  components: {Pane, Hexagon, Room},
  data() {
    return {}
  },
  mounted() {
    pop();

  },
  methods: {
    go(link) {
      this.$router.replace({path: link})
    }
  },
}


function pop() {


}


</script>


<template>
  <div class="element">
    <Room></Room>
  </div>
</template>

<style scoped>
.subroutine {
  width: 10rem;
  display: flex;
  justify-content: start;
  flex-direction: row;
}
</style>
