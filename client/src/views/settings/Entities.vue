<script>

import Header from "../../components/Header.vue";
import Entity from "../../components/entity/Entity.vue";

export default {
  name: "Entities",
  components: {Header, Entity},
  data() {
    return {
      header: {
        icon: "􀛭",
        name: "Entities",
      }
    }
  },
  watch: {},
  beforeMount() {
    /*this.$root.session.subscriptions = []*/

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

  created() {

  },
  methods: {

    setState(name, module, state) {
      this.$root.request("entity", "state", {
        name: name,
        module: module,
        state: state
      })
      console.log("sent")
    },
    groupBy(xs, key) {
      return xs.reduce(function (rv, x) {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
      }, {});
    },

  },
}

</script>


<template>
  <Header :target="header"></Header>


  <div class="context-container context-container-sm gap">
    <div v-for="entity in entities" :key="entity.id">
      <Entity :id="entity.id" :key="entity.id" :entity="entity" small></Entity>
    </div>
  </div>

</template>

<style lang="scss" scoped>

</style>
