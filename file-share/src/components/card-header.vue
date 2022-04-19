<template>
  <div class='header'>
    <p class='title'>{{ pClient.name }}</p>
    <p class='subtitle'>{{ pClient.ipAddress }}</p>
  </div>
</template>

<script>
export default {
  name: 'CardHeader',
  props: {
    pClient: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped lang="scss">
</style>
