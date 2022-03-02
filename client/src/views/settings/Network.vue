<script>
import Pane from "../../components/Pane.vue";
import Group from "../../components/Group.vue";
import NetworkDevice from "./NetworkDevice.vue";
import Header from "../../components/Header.vue";

export default {
  components: {Header, NetworkDevice, Pane, Group},
  data() {
    return {
      header: {icon: '􀆪', name: 'Network'}
    }
  },
  created() {

  },
  computed: {
    devices: function () {
      function compare(a, b) {
        if (a.created < b.created)
          return -1;
        if (a.created > b.created)
          return 1;
        return 0;
      }

      return this.$root.session.metadata.devices.sort(compare)
    }
  },
  methods: {
    groupBy(xs, key) {
      return xs.reduce(function (rv, x) {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
      }, {});
    },
    networkById(networkId) {
      return this.$root.session.metadata.networks.find(n => n.id === networkId)
    }
  }
}
</script>

<template>
  <div>
    <Header :target="header"></Header>
    <div class="context-container context-container-sm gap">

      <div v-for="device in devices">
        <NetworkDevice :device="device" subtitle="Sandbox"/>
      </div>

    </div>
  </div>
</template>

<style lang="scss">
code {
}
</style>
