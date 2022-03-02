<script>

export default {
  name: "Device",
  data() {
    return {
      waiting: false,
      context: false,
      advanced: false,
    }
  },
  props: {
    device: Object,
  },
  methods: {
    toggleContext() {
      this.context = !this.context
    },
    toggleAdvanced() {
      this.advanced = !this.advanced
    }
  }
}

</script>

<template>
  <div class="entity entity-small" @click="toggleContext()">
    <div class="entity-header">
      <div class="icon">
        {{ device.controller ? 'A' : 'B' }}
      </div>
      <div class="label-xxs label-w400 label-o4">
        {{ device.ipv4 }}
      </div>
    </div>
    <div class="fill"></div>
    <div class="label-xxs label-o3 label-w400 px-2 not-important">
      {{ device.mac }}
    </div>
  </div>
  <div v-if="context" class="context" @click="toggleContext()">
    <div class="entity-context gap top flex-column">
      <div class="d-flex flex-column gap px-3 ">
        <div class="d-flex justify-content-start align-items-end align-content-end" v-on:click.stop>
          <div>
              <span
                  class="label-md label-w600 label-o3">d</span>
            <span class="label-md label-w600 label-o6 px-2">Unnamed Device</span>
          </div>
          <div class="fill"></div>
          <div class="h-bar">
            <div class="mx-1 my-1"
                 style="width: 0.0625rem; border-radius: 1rem; background-color: rgba(255,255,255,0.2);"></div>
            <div class="label-sm label-w600 label-o3 px-2 text-uppercase" @click="toggleAdvanced()">
              <span v-if="advanced">􁅦</span> <span v-else>􀍟</span>
            </div>
          </div>
        </div>

        <div class="context-container context-container-lg gap v-bar">
          <div class="entity-small">
            <div class="entity-header">
              <div class="label-xxs label-w400 label-o4 px-1">
                Identifier
              </div>
            </div>
            <div class="fill"></div>
            <div class="label-xxs label-o3 label-w400 px-2">
              <span class="label-w600">MAC</span> {{ device.mac }}
            </div>
          </div>
          <div class="entity-small">
            <div class="entity-header">
              <div class="label-xxs label-w400 label-o4 px-1">
                Addresses
              </div>
            </div>
            <div class="fill"></div>
            <div class="label-xxs label-o3 label-w400 px-2">
              <span class="label-w600">IPV4</span> {{ device.ipv4 || "Unknown" }}
            </div>
            <div class="label-xxs label-o3 label-w400 px-2">
              <span class="label-w600">IPV6</span> {{ device.ipv6 || "::1" }}
            </div>
          </div>
          <div class="entity-small">
            <div class="entity-header">
              <div class="label-xxs label-w400 label-o4 px-1">
                Connected
              </div>
            </div>
            <div class="fill"></div>
            <div class="label-xxs label-o3 label-w400 px-2">
              <span class="label-w600">SINCE</span> {{ this.$timeSince(device.created) || "Unknown" }} ago
            </div>
            <div class="label-xxs label-o3 label-w400 px-2">
              <span class="label-w600">SEEN</span> {{ this.$timeSince(device.updated) || "Unknown" }} ago
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>

</style>
