<script>
import Pane from "../../components/Pane.vue";
import Group from "../../components/Group.vue";
import Header from "../../components/Header.vue";

export default {
  components: {Header, Pane, Group},
  data() {
    return {
      header: {
        icon: "􀄭",
        name: "Connection",
      }
    }
  },
  created() {

  },
  computed: {
    devices: function () {
      function compare(a, b) {
        if (a.name < b.name)
          return -1;
        if (a.name > b.name)
          return 1;
        return 0;
      }

      return this.$root.session.metadata.devices.sort(compare)
    }
  },
  methods: {
    updateStatistics() {

    }
  },
}
</script>

<template>
  <Header :target="header"></Header>
  <div class="h-75 overflow-scroll">

    <div class="element v-bar gap">
      <div class="h-bar justify-content-start align-items-center align-content-center">
        <div class="label-xxs label-o2 label-w600">􀉣</div>
        <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;Controller</div>
        <div class="fill"></div>
        <div class="label-xxs label-o4 label-w500 px-2">&nbsp;&nbsp;Connected</div>

      </div>

      <div class="h-bar gap justify-content-end label-xxs w-100">
        <pre class="text-wrap" v-html="JSON.stringify(this.$root.session.metadata)"></pre>
      </div>
    </div>

  </div>
</template>

<style lang="scss" scoped>

.status > div {
  font-size: 14px;
  margin-right: 0.4em;
  color: rgba(var(--bs-fog-rgb), 0.3);
}

.value {
  font-family: "JetBrains Mono", monospace;
  font-size: 13px;

}

.status > h6 {

  line-height: 1.7em;

  margin-bottom: 0.125em;
}

.status {
  font-size: 0.7rem;
  display: flex;
  align-content: center;
}

.small {
  font-size: 0.7rem;
}
</style>
