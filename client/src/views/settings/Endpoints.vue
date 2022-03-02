<script>

import Header from "../../components/Header.vue";

export default {
  name: "Endpoints",
  components: {Header},
  data() {
    return {
      header: {
        icon: "􀞿",
        name: "Endpoints",
      }
    }
  },
  watch: {},
  beforeMount() {
    /*this.$root.session.subscriptions = []*/

  },
  computed: {
    endpoints: function () {
      function compareConnected(a, b) {
        if (a.connected < b.connected)
          return 1;
        if (a.connected > b.connected)
          return -1;
        return 0;
      }

      function compareName(a, b) {
        if (a.connected < b.connected)
          return 1;
        if (a.connected > b.connected)
          return -1;
        return 0;
      }

      return this.$root.endpoints.sort(compareConnected).sort(compareName)
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

  {{ this.$root.endpoints }}

  <div class="context-container context-container-sm gap">
    <div v-for="endpoint in endpoints" :key="endpoint.id" class="element">
      <div class="h-bar justify-content-start align-items-center align-content-center">
        <div class="label-xxs label-o2 label-w600">􀥔</div>
        <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ endpoint.name }}</div>
        <div class="fill"></div>
        <div v-if="endpoint" class="h-bar gap label-xxs label-o3 label-w400 px-2">
          <div class="label-xxs label-o3">{{ endpoint.connected ? 'Connected' : 'Disconnected' }}</div>
        </div>
      </div>

    </div>
  </div>

</template>

<style lang="scss" scoped>

</style>
