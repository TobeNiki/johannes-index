<template>
  <div class="container-bookmark-view">
    <iframe plugins="true" allow="fullscreen"
      frameborder="no"
      :src="url" />
  </div>
</template>

<script>
export default {
  props: ['url']
}
</script>
<style scoped>
.container-bookmark-view {
  right: 0;
}
</style>
