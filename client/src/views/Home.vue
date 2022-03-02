<script>
import Entity from "../components/entity/Entity.vue";

export default {
  components: {Entity},
  data() {
    return {}
  },
  computed: {
    entities: function () {
      function compare(a, b) {
        if (a.created < b.created)
          return -1;
        if (a.created > b.created)
          return 1;
        return 0;
      }

      return this.$root.entities.sort(compare)
    }
  },
}

</script>

<template>
  <div>
    <div class="d-flex gap-1 mt-1">
      <div v-if="false" class="my-1" style="width: 22rem;">
        <div class="notification">
          <div class="icon">􀛮</div>
          <div class="notification-header">
            <div class="title">
              Remote Proximity Alert
            </div>
            <div class="message">
              Dr. QT Pi has arrived at your location.
            </div>
          </div>
          <div class="time">8m ago</div>
        </div>
      </div>

      <div class="cluster gap">
        <div v-for="entity in entities.filter(f => f.type==='spectrum'||f.type==='switch'||f.type==='dimmer')"
             :key="entity.id">
          <Entity :id="entity.id" :key="entity.id" :entity="entity" small></Entity>
        </div>
      </div>
    </div>
  </div>
</template>

<style>


</style>
