<script>

export default {
  data() {
    return {}
  },
  props: {
    icon: String,
    title: String,
    link: String,
    opaque: Boolean,
    small: Boolean,
  }
}

</script>

<template>
  <router-link :to="link">
    <div :class="`${opaque?'app-opaque':''} ${small?'app-small':''}` " class="surface">
      <div class="app-icon" v-html="icon">

      </div>
    </div>
    <div v-if="opaque" class="app-title">{{ title }}</div>
  </router-link>

</template>
<style lang="scss" scoped>


.app-small .app-outer {
  height: 1.5rem !important;
  width: 1.5rem !important;
  padding: 0;
  margin: 0;
}

.app-outer {
  border-radius: 8px;

  display: flex;
  justify-content: center;
  align-content: center;
  flex-direction: column;
}

.app-opaque {
  height: 4rem !important;
  width: 4rem !important;
}

.app-link {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.app-opaque .app-icon {
  width: 3em;
  height: 3em;
  background-color: rgb(255, 255, 255, 0);
}

.app-small .app-icon {
  font-size: 0.8em;
  padding: 0;
  height: 24px !important;
  width: 32px !important;
}

.app-icon {
  padding: 0.125rem;
  background-color: rgba(255, 255, 255, 0.0425);
  font-family: "SF Pro Display", sans-serif;
  font-size: 1.2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  align-content: center;
  color: rgb(255, 255, 255, 0.5);
}

.router-link-active .app-icon {
  background-color: rgba(255, 255, 255, 0.08);
  background-blend-mode: color-dodge;
  color: rgba(255, 255, 255, 0.75) !important;
  border-radius: 4px;
}

.app-title {
  width: 5em;
  text-align: center;

  margin-top: 4px;
  font-size: 8px;
  font-weight: 400;
  font-family: "SF Pro Display", sans-serif;
}

</style>
